package flexrect

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/dossier/internal/sketcherror"
	"github.com/hansmi/dossier/internal/testutil"
	"github.com/hansmi/dossier/pkg/geometry"
	"github.com/hansmi/dossier/proto/sketchpb"
)

func TestEdgeFromProto(t *testing.T) {
	for _, tc := range []struct {
		name       string
		input      *sketchpb.FlexRect_Edge
		want       genericEdge
		wantErr    error
		wantString string
	}{
		{
			name:    "empty",
			input:   &sketchpb.FlexRect_Edge{},
			wantErr: sketcherror.ErrIncompleteConfig,
		},
		{
			name:  "absolute",
			input: testutil.MustUnmarshalTextproto(t, `abs { cm: 3 }`, &sketchpb.FlexRect_Edge{}),
			want: &absoluteEdge{
				name:     "absolute",
				distance: 3 * geometry.Cm,
			},
			wantString: "absolute (distance 3cm)",
		},
		{
			name: "relative minimal",
			input: testutil.MustUnmarshalTextproto(t, `
rel {
	node: "other"
	feature: TOP_RIGHT
}
`, &sketchpb.FlexRect_Edge{}),
			want: &relativeEdge{
				name: "relative minimal",
				feature: NodeFeature{
					name:    "other",
					feature: sketchpb.NodeFeature_TOP_RIGHT,
				},
			},
			wantString: `relative minimal (feature "other:TOP_RIGHT")`,
		},
		{
			name: "relative",
			input: testutil.MustUnmarshalTextproto(t, `
rel {
	node: "title"
	feature: BOTTOM_LEFT
	offset: { in: 3 }
}
`, &sketchpb.FlexRect_Edge{}),
			want: &relativeEdge{
				name: "relative",
				feature: NodeFeature{
					name:    "title",
					feature: sketchpb.NodeFeature_BOTTOM_LEFT,
				},
				offset: 3 * 2.54 * geometry.Cm,
			},
			wantString: `relative (feature "title:BOTTOM_LEFT", offset 7.62cm)`,
		},
		{
			name:    "relative no node name",
			input:   testutil.MustUnmarshalTextproto(t, `rel {}`, &sketchpb.FlexRect_Edge{}),
			wantErr: sketcherror.ErrIncompleteConfig,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := edgeFromProto(tc.input, tc.name, nil)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got,
					cmp.AllowUnexported(absoluteEdge{}, relativeEdge{}, NodeFeature{}),
					geometry.EquateLength(),
				); diff != "" {
					t.Errorf("edgeFromProto diff (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tc.wantString, got.String()); diff != "" {
					t.Errorf("String() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestEdgePosition(t *testing.T) {
	cb := fakeCallbacks{
		features: map[NodeFeature]geometry.Point{
			{
				name:    "title",
				feature: sketchpb.NodeFeature_BOTTOM_RIGHT,
			}: geometry.Point{
				Left: 4 * geometry.Cm,
				Top:  3 * geometry.Cm,
			},
		},
	}

	for _, tc := range []struct {
		name    string
		input   genericEdge
		want    geometry.Length
		wantErr error
	}{
		{
			name: "absolute",
			input: &absoluteEdge{
				distance: 22 * geometry.Cm,
			},
			want: 22 * geometry.Cm,
		},
		{
			name: "relative",
			input: &relativeEdge{
				feature: NodeFeature{
					name:    "title",
					feature: sketchpb.NodeFeature_BOTTOM_RIGHT,
				},
				offset:  12 * geometry.Cm,
				extract: pointLeft,
			},
			want: (4 + 12) * geometry.Cm,
		},
		{
			name: "relative with unknown node",
			input: &relativeEdge{
				feature: NodeFeature{name: "unknown"},
			},
			wantErr: errUnknownNode,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.input.Position(&cb)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, geometry.EquateLength()); diff != "" {
					t.Errorf("Position() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
