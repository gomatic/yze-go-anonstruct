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
	insp.WithStack([]ast.Node{(*ast.StructType)(nil)}, func(n ast.Node, isPush bool, stack []ast.Node) bool {
		if isPush {
			check(pass, n.(*ast.StructType), stack)
		}
		return true
	})
	return nil, nil
}

// check reports st unless it is an empty struct, the type of a named type
// spec, or a term in generic constraint position.
func check(pass *analysis.Pass, st *ast.StructType, stack []ast.Node) {
	if isEmpty(st) || namesAType(stack) || inConstraintPosition(stack) {
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

// inConstraintPosition reports whether the struct appears in generic
// constraint position: as the operand of a ~ underlying-type term, as a
// type-set term of an interface (including | unions), or directly as a
// type-parameter constraint. Go forbids a defined type after ~, so no
// compliant named-type rewrite exists for these positions.
func inConstraintPosition(stack []ast.Node) bool {
	i := len(stack) - 2
	for isTermExpr(stack[i]) {
		i--
	}
	if _, ok := stack[i].(*ast.Field); !ok {
		return false
	}
	list := stack[i-1].(*ast.FieldList) // a Field's parent is always a FieldList
	return isInterfaceBody(stack[i-2]) || isTypeParams(stack[i-2], list)
}

// isTermExpr reports whether n is part of a type-set term expression: a ~
// underlying-type operator or a | union joining terms.
func isTermExpr(n ast.Node) bool {
	switch e := n.(type) {
	case *ast.UnaryExpr:
		return e.Op == token.TILDE
	case *ast.BinaryExpr:
		return e.Op == token.OR
	}
	return false
}

// isInterfaceBody reports whether n is an interface type, making a Field list
// under it the interface's type-set/method list.
func isInterfaceBody(n ast.Node) bool {
	_, ok := n.(*ast.InterfaceType)
	return ok
}

// isTypeParams reports whether list is the type-parameter list of n (a generic
// function's FuncType or a generic type declaration's TypeSpec), as opposed to
// its parameters, results, or struct fields.
func isTypeParams(n ast.Node, list *ast.FieldList) bool {
	switch p := n.(type) {
	case *ast.FuncType:
		return p.TypeParams == list
	case *ast.TypeSpec:
		return p.TypeParams == list
	}
	return false
}
