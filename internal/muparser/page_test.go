package muparser

import (
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/dossier/internal/mutool/stext"
	"github.com/hansmi/dossier/pkg/geometry"
)

func TestNewPage(t *testing.T) {
	for _, tc := range []struct {
		name    string
		input   stext.Page
		want    *Page
		wantErr error
	}{
		{
			name:    "defaults",
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name: "empty",
			input: stext.Page{
				ID:     "page15354",
				Width:  45 * geometry.Cm,
				Height: 67 * geometry.Cm,
			},
			want: &Page{
				num: 15354,
				size: geometry.Size{
					Width:  45 * geometry.Cm,
					Height: 67 * geometry.Cm,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := newPage(tc.input)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			opts := cmp.Options{
				cmp.AllowUnexported(Page{}),
				cmpopts.EquateEmpty(),
				geometry.EquateLength(),
			}

			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Errorf("newPage() diff (-want +got):\n%s", diff)
			}
		})
	}
}
