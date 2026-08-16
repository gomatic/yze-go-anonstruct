// Package contract carries one anonymous struct whose exact diagnostic text and
// exact position an in-repo test pins. The expectation comments elsewhere in
// this corpus cannot: they match a line against an unanchored regexp, so they
// see neither the column the diagnostic points at nor a message that has grown
// a tail.
package contract

// Value is the anonymous struct the contract test measures.
var Value = struct { // want `anonymous struct with fields; define a named type`
	Key string
}{}
