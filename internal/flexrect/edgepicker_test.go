package flexrect

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/dossier/internal/ref"
	"github.com/hansmi/dossier/internal/sketcherror"
	"github.com/hansmi/dossier/pkg/geometry"
)

func TestEdgePicker(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setup   func(*testing.T, *edgePicker)
		wantErr error
		wantPos *geometry.Length
	}{
		{
			name:    "input",
			setup:   func(*testing.T, *edgePicker) {},
			wantErr: sketcherror.ErrIncompleteConfig,
		},
		{
			name: "conflict",
			setup: func(_ *testing.T, p *edgePicker) {
				p.addEdge(&absoluteEdge{name: "e"}, 0)
				p.addVertex(&absoluteVertex{name: "v"}, pointTop, 0)
			},
			wantErr: sketcherror.ErrBadConfig,
		},
		{
			name: "edge only",
			setup: func(_ *testing.T, p *edgePicker) {
				p.addEdge(&absoluteEdge{name: "three points"}, 3*geometry.Pt)
			},
			wantPos: ref.Ref(3 * geometry.Pt),
		},
		{
			name: "negative distance from vertex",
			setup: func(_ *testing.T, p *edgePicker) {
				p.addVertex(&absoluteVertex{
					name: "right",
					pos:  geometry.Point{Left: 20 * geometry.Cm},
				}, pointLeft, -15*geometry.Cm)
			},
			wantPos: ref.Ref(5 * geometry.Cm),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var p edgePicker

			tc.setup(t, &p)

			got, err := p.pick()

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil && tc.wantPos != nil {
				if gotPos, err := got.Position(nil); err != nil {
					t.Errorf("Position() failed: %v", err)
				} else if diff := cmp.Diff(*tc.wantPos, gotPos, geometry.EquateLength()); diff != "" {
					t.Errorf("Position() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
