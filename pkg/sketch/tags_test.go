package sketch

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestValidateTags(t *testing.T) {
	for _, tc := range []struct {
		name    string
		input   []string
		want    []string
		wantErr error
	}{
		{name: "empty"},
		{
			name:  "valid",
			input: []string{"zlast", "first", "second"},
			want:  []string{"first", "second", "zlast"},
		},
		{
			name:    "empty tag",
			input:   []string{"a", "", "b"},
			wantErr: errInvalidTags,
		},
		{
			name:    "duplicates",
			input:   []string{"a", "b", "another", "a", "c"},
			wantErr: errInvalidTags,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := validateTags(tc.input)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("validateTags() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
