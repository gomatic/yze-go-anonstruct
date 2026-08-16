// Package anonstruct provides a go/analysis analyzer enforcing the gomatic Go
// standard that struct types are named rather than anonymous. It reports an
// anonymous struct type that carries fields, and prescribes one remedy: give
// the type a name.
//
// Three things are not reported, and this list is the whole of them. An
// exemption absent from here gets no corpus case, no forgery probe and no
// reviewer, because every enumerator is told to start from the documentation.
//
// A struct with no fields is not reported. It has no fields to name, and
// struct{} is the idiom for a set's value and a signalling channel's element,
// so the rule has no work to do on it.
//
// A struct that a type declaration names is not reported, and BOTH spellings of
// a declaration name it: `type Point struct{...}` defines a new type, and
// `type Point = struct{...}` gives the same struct type a name without defining
// one. A declaration that names NOTHING is not a name and is reported —
// `type _ = struct{...}` produces no identifier any use site could read, which
// is the whole of what the exemption is for. Parentheses do not hide the
// declaration — `type Point (struct{...})` is legal Go the toolchain compiles as
// a definition, which gofmt leaves alone and gofumpt does not strip — so the
// declaration is looked for through them. gofmt collapses a NESTED pair to a
// single one, so a chain cannot survive in formatted source; the walk still
// steps over the whole chain, because an analyzer judges source as it is written
// rather than as it would be formatted, and an internal test pins that
// independently of any fixture a format sweep could rewrite.
//
// The alias counts because the standard is about NAMES, and because a defined
// type is not always available. Where an anonymous struct is fixed by a
// contract this package does not own — a method implementing an interface whose
// signature spells the struct, a func value passed to such a signature —
// method-signature identity and func-value assignability both require identical
// types, so the defined type fails to compile ("have Emit() Code, want Emit()
// struct{Code int}") and the alias is the only name Go allows. Reporting the
// author who has taken the only remedy available leaves them a disablement as
// their next move, which is worse for the standard than the alias is. A defined
// type remains the better spelling wherever it compiles, because it is a type
// of its own; that preference belongs to review rather than to this rule, which
// asks only that the struct be reachable through a name.
//
// The alias cannot be forged into an escape, because acquiring it IS the
// remedy: `type T = struct{...}` writes the struct literal exactly once, in a
// declaration whose whole purpose is to name it, and every use site then reads
// the name. That is the property the rule exists to produce.
//
// There is deliberately NO exemption for generic constraint position. Go
// forbids a DEFINED type after ~ ("invalid use of ~ (underlying type of Named
// is struct{x int})"), which reads like a position with no compliant rewrite,
// but an alias is accepted there — `type S = struct{x int}` makes `~S` legal —
// so the remedy exists in constraint position exactly as it does everywhere
// else, and a bare type-set term, a union arm and a type-parameter constraint
// all accept a defined type outright.
//
// Test files are not judged. A table-driven test's anonymous struct is the Go
// idiom rather than a defect: it is the test's own local shape, declared and
// consumed in one statement, and naming it publishes a type nothing else uses.
// That is the yze runner's standing decision for this rule — gomatic/yze lists
// anonstruct in its sourceOnly set — and the analyzer applies it itself so that
// the shipped binary, `go vet -vettool` and analysistest, none of which reach
// the runner, all enforce the same rule. The exemption cannot be forged without
// acquiring what it exempts: a file's name ending in _test.go is what makes the
// go tool compile it into the test binary alone, so a file cannot take the
// escape and stay production source.
//
// Generated files are not this analyzer's concern: the yze framework drops
// findings in files carrying the standard generated marker before they are
// reported.
package anonstruct

import (
	"go/ast"
	"go/token"
	"strings"

	goyze "github.com/gomatic/go-yze"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const message = "anonymous struct with fields; define a named type"

// Analyzer reports anonymous struct types that carry fields.
var Analyzer = &analysis.Analyzer{
	Name: "anonstruct",
	Doc: "reports anonymous struct types carrying fields, which the gomatic Go standard forbids in favor of named types. " +
		"Not reported: a struct with no fields, which has nothing to name; a struct a type declaration gives a name to, " +
		"whether it defines a type or aliases one, since either names it, though never `type _`, which names nothing; " +
		"and any struct in a _test.go file, where the table-driven " +
		"literal is the idiom. There is no exemption for generic constraint position: an alias is accepted after ~, so the " +
		"remedy exists there too",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

// Registration declares this analyzer to the yze framework.
//
// TestScope is DECLARED rather than left to the runner's table, because the zero
// value is TestScopeAll and TestScopeAll asserts "reports findings in every
// file" — which this analyzer does not. Leaving it unset would make the
// registration state the opposite of the rule, true only for as long as
// gomatic/yze's sourceOnly table keeps stamping the scope in from outside.
//
// Declaring it adds no second matcher an author could play off against the
// first: check emits nothing for a test file, so the driver's drop has nothing
// left to drop, and the driver decides from the compiled file name exactly as
// check does, so the two cannot disagree even on a file carrying a line
// directive.
var Registration = goyze.Registration{
	Name:       "anonstruct",
	Categories: []goyze.Category{"types", "structure"},
	URL:        "https://docs.gomatic.dev/yze/anonstruct",
	TestScope:  goyze.TestScopeSourceOnly,
	Analyzer:   Analyzer,
}

// testSuffix names the files the go tool compiles into the test binary alone.
const testSuffix = "_test.go"

// blank is the identifier that declares no name for anything to read.
const blank = "_"

// sourceName is a file's name as the FileSet knows it — the name the toolchain
// parsed the file under, which no directive inside the file rewrites.
type sourceName string

// isTest reports whether the name belongs to a file compiled into the test
// binary alone, which this analyzer does not judge.
func (n sourceName) isTest() bool {
	return strings.HasSuffix(string(n), testSuffix)
}

// run reports each anonymous struct type that has fields.
func run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.WithStack([]ast.Node{(*ast.StructType)(nil)}, func(n ast.Node, isPush bool, stack []ast.Node) bool {
		if isPush {
			check(pass, n.(*ast.StructType), stack)
		}
		return true
	})
	return nil, nil
}

// check reports st unless it has no fields, a type declaration names it, or it
// sits in a test file.
func check(pass *analysis.Pass, st *ast.StructType, stack []ast.Node) {
	if isEmpty(st) || namesAType(stack) {
		return
	}
	if compiledName(pass.Fset, st.Pos()).isTest() {
		return
	}
	pass.Reportf(st.Pos(), message)
}

// compiledName returns the name the toolchain parsed pos's file under.
//
// It is token.File.Name() and never fset.Position(...).Filename, because the
// only decision taken from it is a SUPPRESSION and a //line directive is a
// comment inside the judged file: a resolved position would let any source file
// claim the drop with one line naming a _test.go it neither owns nor creates.
//
// An invalid position belongs to no file and yields the empty name, which no
// suppression matches — an unresolvable position is not an argument for silence,
// so the finding stands. This is go-yze's rule for the same question, and the
// two must not answer it differently.
func compiledName(fset *token.FileSet, pos token.Pos) sourceName {
	file := fset.File(pos)
	if file == nil {
		return ""
	}
	return sourceName(file.Name())
}

func isEmpty(st *ast.StructType) bool {
	return len(st.Fields.List) == 0
}

// namesAType reports whether a type declaration gives the struct a name: its
// enclosing node in the traversal stack is a TypeSpec whose identifier is not
// blank, which covers a definition and an alias alike. `type _ = struct{...}`
// declares no identifier, so nothing can read the name and the exemption's whole
// reason is absent.
//
// Parentheses are stepped over because a parenthesized type literal is still
// the declared type, and each pair puts an ast.ParenExpr between the TypeSpec
// and the struct. The walk terminates without a bound because stack[0] is the
// *ast.File the inspector rooted at, which is not an ast.ParenExpr.
func namesAType(stack []ast.Node) bool {
	i := len(stack) - 2
	for isParenthesized(stack[i]) {
		i--
	}
	spec, ok := stack[i].(*ast.TypeSpec)
	return ok && spec.Name.Name != blank
}

// isParenthesized reports whether n is a parenthesized expression, which wraps
// its operand without changing what the operand is.
func isParenthesized(n ast.Node) bool {
	_, ok := n.(*ast.ParenExpr)
	return ok
}
