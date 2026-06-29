// Package anonstruct provides a go/analysis analyzer enforcing the gomatic Go
// standard that struct types are named rather than anonymous. Empty anonymous
// structs (idiomatic for sets and signaling channels) are allowed.
package anonstruct

import (
	"go/ast"
	"go/token"

	goyze "github.com/gomatic/go-yze"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const message = "anonymous struct with fields; define a named type"

// Analyzer reports anonymous struct types that carry fields.
var Analyzer = &analysis.Analyzer{
	Name:     "anonstruct",
	Doc:      "reports anonymous struct types (with fields), which the gomatic Go standard forbids in favor of named types",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

// Registration declares this analyzer to the yze framework.
var Registration = goyze.Registration{
	Name:       "anonstruct",
	Categories: []goyze.Category{"types", "structure"},
	URL:        "https://docs.gomatic.dev/yze/anonstruct",
	Analyzer:   Analyzer,
}

// run reports each anonymous struct type that has fields.
func run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.WithStack([]ast.Node{(*ast.StructType)(nil)}, func(n ast.Node, push bool, stack []ast.Node) bool {
		if push {
			check(pass, n.(*ast.StructType), stack)
		}
		return true
	})
	return nil, nil
}

// check reports st unless it is an empty struct or the type of a named type spec.
func check(pass *analysis.Pass, st *ast.StructType, stack []ast.Node) {
	if isEmpty(st) || namesAType(stack) {
		return
	}
	pass.Reportf(st.Pos(), message)
}

func isEmpty(st *ast.StructType) bool {
	return len(st.Fields.List) == 0
}

// namesAType reports whether the struct is the right-hand side of a type
// definition (its immediate parent in the traversal stack is a TypeSpec with no
// Assign position). A type alias (`type A = struct{...}`) carries an Assign
// position and does not name a new struct type, so it is not exempt.
func namesAType(stack []ast.Node) bool {
	ts, ok := stack[len(stack)-2].(*ast.TypeSpec)
	return ok && ts.Assign == token.NoPos
}
