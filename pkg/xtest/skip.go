package xtest

import (
	"os"
	"slices"
	"strings"
	"testing"
)

const envVar = "RUN_TEST"

const (
	Unit        = "unit"
	Integration = "integration"
	System      = "system"
	Validation  = "validation"
)

func SkipValidationIfRequired(t *testing.T) {
	SkipIfRequired(t, Validation)
}

func SkipSystemIfRequired(t *testing.T) {
	SkipIfRequired(t, System)
}

func SkipIntegrationIfRequired(t *testing.T) {
	SkipIfRequired(t, Integration)
}

func SkipUnitIfRequired(t *testing.T) {
	SkipIfRequired(t, Unit)
}

func SkipIfRequired(t *testing.T, test string) {
	v, ok := os.LookupEnv(envVar)
	if !ok {
		return
	}

	tests := strings.Split(v, ",")
	if !slices.Contains(tests, test) {
		t.Skipf("Set %s to %s run this test", envVar, test)
	}
}
