package cliutil

import (
	"regexp"
	"testing"

	"github.com/hansmi/dossier/pkg/geometry"
)

func TestLengthUnit(t *testing.T) {
	var flagValue geometry.LengthUnit

	luv := NewLengthUnitVar(&flagValue, geometry.Inch)

	gotUsage := luv.Usage("foo bar")

	if wantUsage := regexp.MustCompile(`(?im)^foo bar\n.*\bsupported:.*\bcm\b`); !wantUsage.MatchString(gotUsage) {
		t.Errorf("Usage() result doesn't match %q: %q", wantUsage.String(), gotUsage)
	}
}
