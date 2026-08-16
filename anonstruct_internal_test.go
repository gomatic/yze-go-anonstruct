package anonstruct

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"
)

// TestIsTestReadsTheSuffixAndNothingElse varies the dimension every fixture in
// this repository holds constant. Every case in the corpus is a bare relative
// path inside one directory, so a widening keyed on a directory segment, on a
// prefix, or on an absolute path would be invisible to it — and being a disjunct
// added to an existing return, invisible to statement coverage too. Real runs
// hand this matcher absolute paths, so the cases are absolute.
func TestIsTestReadsTheSuffixAndNothingElse(t *testing.T) {
	for name, want := range map[sourceName]bool{
		"/Users/x/src/pkg/thing_test.go":          true,
		"/Users/x/src/pkg/internal/thing_test.go": true,
		"/Users/x/src/pkg/thing.go":               false,
		"/Users/x/src/pkg/internal/thing.go":      false,
		"/Users/x/src/pkg/thingtest.go":           false,
		"/Users/x/src/pkg/_test.go.orig":          false,
		"/Users/x/src/test/thing.go":              false,
		"/private/tmp/build/thing.go":             false,
	} {
		assert.Equal(t, want, name.isTest(), "isTest(%q)", name)
	}
}

// TestCheckJudgesTheCompiledNameNotTheRewrittenPosition pins the divergence no
// corpus case can see, because no fixture anywhere carries the construct. A
// //line directive is a comment INSIDE the judged file and it rewrites what
// fset.Position reports, so a scope decided from a resolved path is a decision
// the judged file writes for itself: one line naming a _test.go it neither owns
// nor has to create would silence an ordinary source file. token.File.Name() is
// the name the toolchain parsed the file under, and no directive rewrites it.
func TestCheckJudgesTheCompiledNameNotTheRewrittenPosition(t *testing.T) {
	const source = "//line zz_test.go:1\npackage p\n\nvar v = struct{ x int }{}\n"

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "pkg/thing.go", source, parser.ParseComments)
	require.NoError(t, err)

	structType := findStructType(t, file)
	// go/parser resolves the directive's relative name against the directory of
	// the file it parsed, so the rewritten position is pkg/zz_test.go.
	require.Equal(t, "pkg/zz_test.go", fset.Position(structType.Pos()).Filename,
		"the directive must take effect, or this case proves nothing")
	require.Equal(t, "pkg/thing.go", fset.File(structType.Pos()).Name())

	reported := 0
	pass := &analysis.Pass{Fset: fset, Report: func(analysis.Diagnostic) { reported++ }}
	check(pass, structType, []ast.Node{file, structType})

	assert.Equal(t, 1, reported, "a source file does not become a test file by naming one in a comment")
}

// findStructType returns the only struct type literal in file.
func findStructType(t *testing.T, file *ast.File) *ast.StructType {
	t.Helper()

	var found *ast.StructType
	ast.Inspect(file, func(n ast.Node) bool {
		if structType, ok := n.(*ast.StructType); ok {
			found = structType
		}
		return found == nil
	})
	require.NotNil(t, found)

	return found
}

// TestAPositionTheFileSetDoesNotKnowIsNotAnArgumentForSilence pins the DIRECTION
// of the fail-safe, which is the opposite of the obvious one. The only decision
// taken from the file name is a suppression, so a position with no name must not
// acquire it. go-yze answers the same question the same way for the framework's
// own drops (compiled.go), and an analyzer failing silent where the framework
// fails loud would be a second, quieter policy nobody declared.
func TestAPositionTheFileSetDoesNotKnowIsNotAnArgumentForSilence(t *testing.T) {
	reported := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(analysis.Diagnostic) { reported++ },
	}

	structType := &ast.StructType{Fields: &ast.FieldList{List: []*ast.Field{{}}}}
	check(pass, structType, []ast.Node{&ast.File{}, structType})

	assert.Equal(t, 1, reported, "an unresolvable position is not an argument for silence")
}

// TestNamesATypeStepsOverEveryParenthesis pins the walk where no fixture can.
// gofmt collapses a nested pair to a single one, so `type T ((struct{...}))`
// survives only until somebody formats the file — and the case that dies with it
// is the one proving the walk is a loop rather than a single step. The stack is
// built here instead, where no formatter reaches it.
func TestNamesATypeStepsOverEveryParenthesis(t *testing.T) {
	structType := &ast.StructType{Fields: &ast.FieldList{List: []*ast.Field{{}}}}
	named := &ast.TypeSpec{Name: ast.NewIdent("Point")}
	blankSpec := &ast.TypeSpec{Name: ast.NewIdent("_")}

	assert.True(t, namesAType(stackUnder(named, structType, 0)), "no parentheses")
	assert.True(t, namesAType(stackUnder(named, structType, 1)), "one pair")
	assert.True(t, namesAType(stackUnder(named, structType, 2)), "two pairs")
	assert.True(t, namesAType(stackUnder(named, structType, 5)), "five pairs")
	assert.False(t, namesAType(stackUnder(blankSpec, structType, 2)), "a blank declaration names nothing")
	assert.False(t, namesAType(stackUnder(&ast.ValueSpec{}, structType, 2)), "not a declaration at all")
}

// stackUnder builds an inspector-shaped stack: the file, the enclosing node,
// depth parenthesis wrappers, then the struct.
func stackUnder(enclosing ast.Node, structType *ast.StructType, depth int) []ast.Node {
	stack := []ast.Node{&ast.File{}, enclosing}
	for range depth {
		stack = append(stack, &ast.ParenExpr{})
	}

	return append(stack, structType)
}
