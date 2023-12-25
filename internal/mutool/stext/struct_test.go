package stext

import (
	"encoding/xml"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/dossier/pkg/geometry"
)

func TestChar(t *testing.T) {
	for _, tc := range []struct {
		name    string
		input   string
		wantErr error
		want    Char
	}{
		{
			name:    "empty",
			wantErr: io.EOF,
		},
		{
			name:  "element only",
			input: `<char/>`,
		},
		{
			name:  "full",
			input: `<char quad="140.74666 211.81932 147.69416 211.81932 140.74666 223.46102 147.69416 223.46102" x="140.74666" y="221.1024" color="#000000" c="R"/>`,
			want: Char{
				C: 'R',
				Bounds: geometry.Rect{
					Left:   140.75 * geometry.Pt,
					Top:    211.81 * geometry.Pt,
					Right:  147.69 * geometry.Pt,
					Bottom: 223.46 * geometry.Pt,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got Char

			err := xml.Unmarshal([]byte(tc.input), &got)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, geometry.EquateLength()); diff != "" {
					t.Errorf("Element diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
