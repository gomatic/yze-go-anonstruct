package a

// Named is a proper named struct type and must not be flagged.
type Named struct {
	a int
}

// AliasAnon gives an anonymous struct a name without defining a new type. The
// standard is about names, and the alias supplies one, so it is not flagged.
type AliasAnon = struct {
	x int
}

// ParenDefined parenthesizes the type literal of a definition. Parentheses are
// legal here, gofmt -s leaves them alone, and the declaration still names the
// struct, so it is not flagged.
type ParenDefined (struct {
	p int
})

// ParenParenDefined nests the parentheses, proving the declaration is looked
// for through however many of them there are.
type ParenParenDefined ((struct {
	pp int
}))

// ParenAlias parenthesizes an alias's type literal; it names the struct too.
type ParenAlias = (struct {
	pa int
})

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
// forbids a DEFINED type there, but an alias is accepted — `type S = struct{ x
// int }` makes `~S` legal — so a compliant rewrite exists and the term is
// flagged like any other anonymous struct.
type TildeConstraint interface {
	~struct{ x int } // want `anonymous struct with fields; define a named type`
}

// BareTermConstraint embeds a struct directly as an interface type-set term. A
// defined type is legal in that position, so the rewrite exists and the term is
// flagged.
type BareTermConstraint interface {
	struct{ b int } // want `anonymous struct with fields; define a named type`
}

// UnionConstraint unions struct terms in an interface type set. Both arms have
// a compliant rewrite and both are flagged.
type UnionConstraint interface {
	~struct{ u int } | struct{ w int } // want `anonymous struct with fields; define a named type` `anonymous struct with fields; define a named type`
}

// genericFunc constrains its type parameter directly with a struct type. A
// defined type is legal there, so it is flagged.
func genericFunc[P struct{ p int }]() { // want `anonymous struct with fields; define a named type`
	var _ P
}

// genericFuncTilde constrains its type parameter with a ~ struct term; an alias
// is legal after ~, so it is flagged.
func genericFuncTilde[P ~struct{ q int }]() { // want `anonymous struct with fields; define a named type`
	var _ P
}

// GenericType constrains its type parameter in a generic type declaration. The
// constraint accepts a defined type, so the struct is flagged.
type GenericType[P struct{ g int }] struct { // want `anonymous struct with fields; define a named type`
	v P
}

// NamedConstraint proves the prescribed rewrite is accepted where the analyzer
// prescribes it: a defined type as a bare type-set term. Nothing from here to
// remediedType is flagged, and the package compiles, which is what makes the
// remedy actionable rather than advice.
type NamedConstraint interface {
	Named
}

// AliasedForTilde is the alias spelling ~ accepts.
type AliasedForTilde = struct {
	t int
}

// TildeNamed applies ~ to the alias, which is the rewrite TildeConstraint's
// diagnostic prescribes.
type TildeNamed interface {
	~AliasedForTilde
}

// remediedFunc takes the defined type as a direct type-parameter constraint.
func remediedFunc[P Named]() { var _ P }

// RemediedType takes the defined type as a generic type's constraint.
type RemediedType[P Named] struct {
	v P
}

var _ = remediedFunc[Named]

var _ RemediedType[Named]

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

// twoFields carries more than one field, so a rule keyed on the number of
// fields cannot pass by exempting everything above one.
func twoFields() {
	var v struct { // want `anonymous struct with fields; define a named type`
		x int
		y int
	}
	_ = v
}

// threeFields carries three, so the same widening cannot pass at two either.
func threeFields() {
	var v struct { // want `anonymous struct with fields; define a named type`
		x int
		y int
		z int
	}
	_ = v
}

// embeddedOnly carries a single embedded field, which has no name of its own.
// A field count is not what makes a struct anonymous, and neither is a name on
// the field, so this is reported like any other.
func embeddedOnly() {
	var v struct { // want `anonymous struct with fields; define a named type`
		Named
	}
	_ = v
}

// exportedFields carries exported fields only, proving the rule is not keyed on
// the case of the field names.
var exportedFields = struct { // want `anonymous struct with fields; define a named type`
	X int
	Y int
}{}

// A declaration names the struct only when the struct IS the declared type. The
// six below each name something ELSE that merely contains an anonymous struct —
// a slice, a map, a func, a channel, a pointer, an array — and every one is
// still reported. They vary the kind of node sitting between the TypeSpec and
// the struct, which is the dimension the rest of this file holds constant: any
// widening of the walk that looks past that node for a TypeSpec silences all of
// them at once, and nothing else here would notice.

// Rows names a slice, not the element struct.
type Rows []struct { // want `anonymous struct with fields; define a named type`
	r int
}

// Index names a map, not the value struct.
type Index map[string]struct { // want `anonymous struct with fields; define a named type`
	i int
}

// Callback names a func type, not its parameter struct.
type Callback func(struct { // want `anonymous struct with fields; define a named type`
	c int
})

// Pipe names a channel, not its element struct.
type Pipe chan struct { // want `anonymous struct with fields; define a named type`
	p int
}

// PtrTo names a pointer, not its pointee struct.
type PtrTo *struct { // want `anonymous struct with fields; define a named type`
	q int
}

// Arr names an array, not its element struct.
type Arr [3]struct { // want `anonymous struct with fields; define a named type`
	s int
}
