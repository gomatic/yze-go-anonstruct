package a

// Named is a proper named struct type and must not be flagged.
type Named struct {
	a int
}

// Set uses an empty anonymous struct (idiomatic); it must not be flagged.
type Set map[string]struct{}

// withVar uses an anonymous struct with fields as a variable type and is flagged.
func withVar() {
	var v struct { // want `anonymous struct`
		x int
	}
	_ = v
}

// withParam takes an anonymous struct with fields and is flagged.
func withParam(p struct { // want `anonymous struct`
	y int
}) {
	_ = p
}

// Outer nests an anonymous struct in a field and is flagged at the inner struct.
type Outer struct {
	inner struct { // want `anonymous struct`
		z int
	}
}
