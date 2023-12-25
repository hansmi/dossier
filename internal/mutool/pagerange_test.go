package mutool

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hansmi/dossier/pkg/pagerange"
)

func TestFormatPageRange(t *testing.T) {
	for _, tc := range []struct {
		input pagerange.Range
		want  string
	}{
		{want: "0-0"},
		{input: pagerange.MustSingle(1), want: "1-1"},
		{input: pagerange.MustNew(1, 22), want: "1-22"},
		{input: pagerange.MustNew(10, 33), want: "10-33"},
		{input: pagerange.MustNew(pagerange.Last, pagerange.Last), want: "N"},
		{input: pagerange.MustNew(10, pagerange.Last), want: "10-N"},
		{input: pagerange.All, want: "1-N"},
	} {
		t.Run(tc.input.String(), func(t *testing.T) {
			got := formatPageRange(tc.input)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("formatPageRange() diff (-want +got):\n%s", diff)
			}
		})
	}
}
