package sketch

import (
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestEvaluateMatch(t *testing.T) {
	for _, tc := range []struct {
		name   string
		expr   string
		text   string
		want   []TextMatchGroup
		checks func(*testing.T, *TextMatch)
	}{
		{
			name: "empty",
			want: []TextMatchGroup{{}},
		},
		{
			name: "hello world",
			expr: `(?i)hello world (?P<num>\d+) (\w+) (2\d*)`,
			text: "hello world 123 aldksjf 21",
			want: []TextMatchGroup{
				{End: 26, Text: "hello world 123 aldksjf 21"},
				{Name: "num", Start: 12, End: 15, Text: "123"},
				{Start: 16, End: 23, Text: "aldksjf"},
				{Start: 24, End: 26, Text: "21"},
			},
			checks: func(t *testing.T, m *TextMatch) {
				for _, mg := range []*TextMatchGroup{m.Group(1), m.Named("num")} {
					if diff := cmp.Diff("123", mg.Text); diff != "" {
						t.Errorf("Match text diff (-want +got):\n%s", diff)
					}
				}
			},
		},
		{
			name: "optional group without match",
			expr: `text(optional)?`,
			text: "text",
			want: []TextMatchGroup{
				{End: 4, Text: "text"},
				{Name: "", Start: -1, End: -1},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := evaluateMatch(regexp.MustCompile(tc.expr), tc.text)

			if got == nil {
				t.Fatalf("evaluateMatch() returned nil")
			}

			if diff := cmp.Diff(tc.want, got.Groups(), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Match groups diff (-want +got):\n%s", diff)
			}

			for _, mg := range []*TextMatchGroup{got.Group(0), got.Named("")} {
				if diff := cmp.Diff(tc.want[0], *mg, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("Match group diff (-want +got):\n%s", diff)
				}
			}

			if tc.checks != nil {
				tc.checks(t, got)
			}
		})
	}
}
