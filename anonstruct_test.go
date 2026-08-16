package anonstruct_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	goyze "github.com/gomatic/go-yze"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"

	anonstruct "github.com/gomatic/yze-go-anonstruct"
)

// fixtureFile is a path relative to testdata/src.
type fixtureFile string

// sourceToken is the exact text a position is required to point at.
type sourceToken string

// wantMessage is the diagnostic text the doc comment obliges this analyzer to
// emit, spelled here in full rather than as a pattern: the `// want` comments
// match an unanchored regexp, so a message that grew a prescription nobody
// implemented would still match them.
const wantMessage = "anonymous struct with fields; define a named type"

func TestAnonymousStructsAreReported(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), anonstruct.Analyzer, "a", "internal/shape", "contract")
}

// TestDiagnosticIsTheContractedTextAtTheContractedToken pins what neither
// instrument can see. testdata.stickler's Finding is {rule, path, line} and
// analysistest matches a line against an unanchored regexp, so neither can tell
// the contracted message from one carrying a suffix, and neither reads the
// column at all. The position required is the `struct` keyword that opens the
// type literal — the token the author has to change — and it is resolved from
// the fixture source rather than copied from a run.
func TestDiagnosticIsTheContractedTextAtTheContractedToken(t *testing.T) {
	results := analysistest.Run(t, analysistest.TestData(), anonstruct.Analyzer, "contract")
	require.Len(t, results, 1)

	diagnostics := results[0].Action.Diagnostics
	require.Len(t, diagnostics, 1, "the contract fixture holds exactly one anonymous struct")
	assert.Equal(t, wantMessage, diagnostics[0].Message)

	wantLine, wantColumn := positionOf(t, "contract/contract.go", "struct {")
	got := results[0].Action.Package.Fset.Position(diagnostics[0].Pos)
	assert.Equal(t, wantLine, got.Line)
	assert.Equal(t, wantColumn, got.Column, "the diagnostic points at the struct keyword, not at the brace after it")
}

// TestRegistrationIsWellFormed pins every field of the registration rather than
// only its validity. Categories are what the runner files a finding under and
// URL is where it sends the author to read the rule; Validate() returning nil
// says nothing about either, and the corpus Finding carries neither.
func TestRegistrationIsWellFormed(t *testing.T) {
	assert.NoError(t, anonstruct.Registration.Validate())
	assert.Equal(t, "yze/anonstruct", anonstruct.Registration.RuleID())
	assert.Same(t, anonstruct.Analyzer, anonstruct.Registration.Analyzer)
	assert.Equal(t, goyze.AnalyzerName("anonstruct"), anonstruct.Registration.Name)
	assert.Equal(t, []goyze.Category{"types", "structure"}, anonstruct.Registration.Categories)
	assert.Equal(t, goyze.HelpURL("https://docs.gomatic.dev/yze/anonstruct"), anonstruct.Registration.URL)
}

// positionOf resolves a 1-based line and column for the first occurrence of
// token in a fixture, so a position assertion states which token it requires
// instead of repeating a number some earlier run produced.
func positionOf(t *testing.T, file fixtureFile, token sourceToken) (int, int) {
	t.Helper()

	source, err := os.ReadFile(filepath.Join(analysistest.TestData(), "src", string(file)))
	require.NoError(t, err)

	offset := strings.Index(string(source), string(token))
	require.GreaterOrEqual(t, offset, 0, "fixture %s does not contain %q", file, token)

	head := string(source[:offset])

	return strings.Count(head, "\n") + 1, len(head) - strings.LastIndex(head, "\n")
}
