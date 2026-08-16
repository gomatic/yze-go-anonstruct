package a

import "testing"

// TestTable holds the table-driven anonymous struct this rule exists to leave
// alone. It is the identical shape to exportedFields in a.go, which IS
// reported, so silence here can only be the file's role and not the type's
// shape — and the pair is what pins the scope: delete the test-file guard and
// this file starts reporting.
func TestTable(t *testing.T) {
	for _, tc := range []struct {
		Name string
		Want int
	}{{Name: "one", Want: 1}} {
		if tc.Want != 1 {
			t.Errorf("%s: want 1", tc.Name)
		}
	}
}
