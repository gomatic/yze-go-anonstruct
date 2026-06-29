package anonstruct_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"

	anonstruct "github.com/gomatic/yze-anonstruct"
)

func TestAnonymousStructsAreReported(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), anonstruct.Analyzer, "a")
}

func TestRegistrationIsWellFormed(t *testing.T) {
	assert.NoError(t, anonstruct.Registration.Validate())
	assert.Equal(t, "yze/anonstruct", anonstruct.Registration.RuleID())
	assert.Same(t, anonstruct.Analyzer, anonstruct.Registration.Analyzer)
}
