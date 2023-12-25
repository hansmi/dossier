package geometry

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestEquateLength(t *testing.T) {
	for _, tc := range []struct {
		name     string
		x, y     any
		opts     cmp.Options
		wantDiff bool
	}{
		{name: "no options, pt", x: 10 * Pt, y: 10 * Pt},
		{name: "no options, mm", x: 20 * Mm, y: 20 * Mm},
		{
			name:     "9mm difference",
			x:        Mm,
			y:        Cm,
			opts:     cmp.Options{EquateLength()},
			wantDiff: true,
		},
		{
			name: "twice the distance",
			x:    Pt,
			y:    2 * Pt,
			opts: cmp.Options{
				cmpopts.EquateApprox(0, 999),
				EquateLength(),
			},
		},
		{
			name: "one meter",
			x:    100 * Cm,
			y:    100 * Cm,
			opts: cmp.Options{EquateLengthApprox(Pt)},
		},
		{
			name:     "one meter, difference too large",
			x:        100 * Cm,
			y:        100.1 * Cm,
			opts:     cmp.Options{EquateLengthApprox(Pt)},
			wantDiff: true,
		},
		{
			name: "one meter, small difference",
			x:    100 * Cm,
			y:    100.1 * Cm,
			opts: cmp.Options{EquateLengthApprox(5 * Pt)},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			diff := cmp.Diff(tc.x, tc.y, tc.opts...)

			if tc.wantDiff != (diff != "") {
				t.Errorf("Diff(%v, %v) result:\n%s", tc.x, tc.y, diff)
			}
		})
	}
}
