// Package shape sits under a directory segment named "internal", the segment a
// path-keyed exemption reaches for first. The rule is about the type, so the
// same struct is reported wherever the file lives, and the corpus varies the
// path so that a rule keyed on it cannot pass.
package shape

// Pair is the named type the rule prescribes, and is not reported.
type Pair struct {
	Left  int
	Right int
}

// pair holds the anonymous spelling of the same shape and is reported.
var pair = struct { // want `anonymous struct with fields; define a named type`
	Left  int
	Right int
}{}

// Sum reads both so neither declaration is unused.
func Sum() int { return pair.Left + pair.Right + Pair{}.Left }
