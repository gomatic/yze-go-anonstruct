package a

// Named is a proper named struct type and must not be flagged.
type Named struct {
	a int
}

// AliasAnon aliases an anonymous struct with fields; the alias does not name a
// new type, so the struct is still anonymous and is flagged.
type AliasAnon = struct { // want `anonymous struct with fields; define a named type`
	x int
}

// AliasEmpty aliases an empty anonymous struct (idiomatic); it is not flagged.
type AliasEmpty = struct{}

// Set uses an empty anonymous struct (idiomatic); it must not be flagged.
type Set map[string]struct{}

// emptyVar uses a standalone empty anonymous struct (idiomatic); not flagged.
var emptyVar struct{}

// withVar uses an anonymous struct with fields as a variable type and is flagged.
func withVar() {
	var v struct { // want `anonymous struct with fields; define a named type`
		x int
	}
	_ = v
}

// withParam takes an anonymous struct with fields and is flagged.
func withParam(p struct { // want `anonymous struct with fields; define a named type`
	y int
}) {
	_ = p
}

// compositeLit constructs an anonymous struct with fields via a composite literal
// and is flagged at the struct type of the literal.
func compositeLit() {
	_ = struct { // want `anonymous struct with fields; define a named type`
		c int
	}{}
}

// returnsAnon declares an anonymous struct with fields as its result type and is
// flagged there. It panics rather than returning a value so the result type is
// the only anonymous struct node in the function.
func returnsAnon() struct { // want `anonymous struct with fields; define a named type`
	r int
} {
	panic("unreachable")
}

// mapElem has an anonymous struct with fields as its value type and is flagged.
var mapElem map[string]struct { // want `anonymous struct with fields; define a named type`
	k int
}

// sliceElem has an anonymous struct with fields as its element type and is flagged.
var sliceElem []struct { // want `anonymous struct with fields; define a named type`
	e int
}

// Outer nests an anonymous struct in a field and is flagged at the inner struct.
type Outer struct {
	inner struct { // want `anonymous struct with fields; define a named type`
		z int
	}
}

// TildeConstraint uses a struct as the operand of a ~ underlying-type term. Go
// forbids a named (defined) type after ~, so no compliant rewrite exists and
// the struct is not flagged.
type TildeConstraint interface {
	~struct{ x int }
}

// BareTermConstraint embeds a struct directly as an interface type-set term
// (constraint position); it is not flagged.
type BareTermConstraint interface {
	struct{ b int }
}

// UnionConstraint unions struct terms in an interface type set (constraint
// position); neither term is flagged.
type UnionConstraint interface {
	~struct{ u int } | struct{ w int }
}

// genericFunc constrains its type parameter directly with a struct type
// (constraint position in the type-parameter list); it is not flagged.
func genericFunc[P struct{ p int }]() {
	var _ P
}

// genericFuncTilde constrains its type parameter with a ~ struct term directly
// in the type-parameter list; it is not flagged.
func genericFuncTilde[P ~struct{ q int }]() {
	var _ P
}

// GenericType constrains its type parameter in a generic type declaration
// (constraint position); the constraint struct is not flagged.
type GenericType[P struct{ g int }] struct {
	v P
}

// genericBody proves type parameters do not blanket-exempt a generic function:
// an anonymous struct in the *body* is still flagged.
func genericBody[P any]() {
	var v struct { // want `anonymous struct with fields; define a named type`
		x int
	}
	_ = v
}

// MethodParamInInterface proves interface bodies are not blanket-exempt: a
// struct used as a method parameter inside an interface is an ordinary
// anonymous struct, not a type-set term, and is flagged.
type MethodParamInInterface interface {
	M(struct { // want `anonymous struct with fields; define a named type`
		m int
	})
}
