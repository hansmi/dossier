package pagerange

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNew(t *testing.T) {
	for _, tc := range []struct {
		name    string
		lower   int
		upper   int
		want    Range
		wantErr error
		wantStr string
	}{
		{name: "zero", wantErr: os.ErrInvalid},
		{name: "upper before lower", lower: 100, wantErr: os.ErrInvalid},
		{name: "upper before last", lower: Last, upper: 100, wantErr: os.ErrInvalid},

		{name: "all", lower: 1, upper: Last, want: Range{1, Last}, wantStr: "1-(last)"},
		{name: "first", lower: 1, upper: 1, want: Range{1, 1}, wantStr: "1"},
		{name: "second", lower: 2, upper: 2, want: Range{2, 2}, wantStr: "2"},
		{name: "range", lower: 10, upper: 19, want: Range{10, 19}, wantStr: "10-19"},
		{name: "range to last", lower: 30, upper: Last, want: Range{30, Last}, wantStr: "30-(last)"},
		{name: "last", lower: Last, upper: Last, want: Range{Last, Last}, wantStr: "(last)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := New(tc.lower, tc.upper)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("New() diff (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tc.wantStr, got.String()); diff != "" {
					t.Errorf("String() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
