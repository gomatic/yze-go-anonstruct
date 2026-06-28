package anonstruct_test

import (
	"testing"

	anonstruct "github.com/gomatic/yze-go-anonstruct"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnonymousStructsAreReported(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), anonstruct.Analyzer, "a")
}

func TestRegistrationIsWellFormed(t *testing.T) {
	assert.NoError(t, anonstruct.Registration.Validate())
	assert.Equal(t, "yze/go/anonstruct", anonstruct.Registration.RuleID())
	assert.Same(t, anonstruct.Analyzer, anonstruct.Registration.Analyzer)
}
