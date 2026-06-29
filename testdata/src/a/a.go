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
