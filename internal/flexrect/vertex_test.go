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

func TestVertexFromProto(t *testing.T) {
	for _, tc := range []struct {
		name       string
		input      *sketchpb.FlexRect_Vertex
		want       genericVertex
		wantErr    error
		wantString string
	}{
		{
			name:    "empty",
			input:   &sketchpb.FlexRect_Vertex{},
			wantErr: sketcherror.ErrIncompleteConfig,
		},
		{
			name:  "absolute",
			input: testutil.MustUnmarshalTextproto(t, `abs { left { cm: 3 } top { in: 2 } }`, &sketchpb.FlexRect_Vertex{}),
			want: &absoluteVertex{
				name: "absolute",
				pos: geometry.Point{
					Left: 3 * geometry.Cm,
					Top:  2 * geometry.Inch,
				},
			},
			wantString: "absolute (position (3cm, 5.08cm))",
		},
		{
			name: "relative minimal",
			input: testutil.MustUnmarshalTextproto(t, `
rel {
	node: "foo"
	feature: BOTTOM_RIGHT
}
`, &sketchpb.FlexRect_Vertex{}),
			want: &relativeVertex{
				name: "relative minimal",
				feature: NodeFeature{
					name:    "foo",
					feature: sketchpb.NodeFeature_BOTTOM_RIGHT,
				},
			},
			wantString: `relative minimal (feature "foo:BOTTOM_RIGHT")`,
		},
		{
			name: "relative",
			input: testutil.MustUnmarshalTextproto(t, `
rel {
	node: "title"
	feature: TOP_LEFT
	offset: { width: { in: 5 } height: { cm: 2 } }
}
`, &sketchpb.FlexRect_Vertex{}),
			want: &relativeVertex{
				name: "relative",
				feature: NodeFeature{
					name:    "title",
					feature: sketchpb.NodeFeature_TOP_LEFT,
				},
				offset: geometry.Size{
					Width:  5 * geometry.Inch,
					Height: 2 * geometry.Cm,
				},
			},
			wantString: `relative (feature "title:TOP_LEFT", offset (12.7cm, 2cm))`,
		},
		{
			name:    "relative no node name",
			input:   testutil.MustUnmarshalTextproto(t, `rel {}`, &sketchpb.FlexRect_Vertex{}),
			wantErr: sketcherror.ErrIncompleteConfig,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := vertexFromProto(tc.input, tc.name)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got,
					cmp.AllowUnexported(absoluteVertex{}, relativeVertex{}, NodeFeature{}),
					geometry.EquateLength(),
				); diff != "" {
					t.Errorf("vertexFromProto diff (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tc.wantString, got.String()); diff != "" {
					t.Errorf("String() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestVertexPosition(t *testing.T) {
	cb := fakeCallbacks{
		features: map[NodeFeature]geometry.Point{
			{
				name:    "title",
				feature: sketchpb.NodeFeature_BOTTOM_RIGHT,
			}: geometry.Point{
				Left: 10 * geometry.Cm,
				Top:  20 * geometry.Cm,
			},
		},
	}

	for _, tc := range []struct {
		name    string
		input   genericVertex
		want    geometry.Point
		wantErr error
	}{
		{
			name: "absolute",
			input: &absoluteVertex{
				pos: geometry.Point{
					Left: 22 * geometry.Cm,
					Top:  33 * geometry.Cm,
				},
			},
			want: geometry.Point{
				Left: 22 * geometry.Cm,
				Top:  33 * geometry.Cm,
			},
		},
		{
			name: "relative",
			input: &relativeVertex{
				feature: NodeFeature{
					name:    "title",
					feature: sketchpb.NodeFeature_BOTTOM_RIGHT,
				},
				offset: geometry.Size{
					Width:  12 * geometry.Cm,
					Height: 34 * geometry.Cm,
				},
			},
			want: geometry.Point{
				Left: (10 + 12) * geometry.Cm,
				Top:  (20 + 34) * geometry.Cm,
			},
		},
		{
			name: "relative with unknown node",
			input: &relativeVertex{
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
