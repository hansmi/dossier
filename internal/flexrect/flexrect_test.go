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

func TestFromProto(t *testing.T) {
	cb := fakeCallbacks{
		features: map[NodeFeature]geometry.Point{
			{
				name:    "title",
				feature: sketchpb.NodeFeature_BOTTOM_LEFT,
			}: geometry.Point{
				Left: 3 * geometry.Cm,
				Top:  8 * geometry.Cm,
			},
		},
	}

	for _, tc := range []struct {
		name    string
		input   *sketchpb.FlexRect
		wantErr error
		want    geometry.Rect
	}{
		{
			name:    "empty",
			input:   &sketchpb.FlexRect{},
			wantErr: sketcherror.ErrIncompleteConfig,
		},
		{
			name:    "incomplete edge",
			input:   testutil.MustUnmarshalTextproto(t, `top {}`, &sketchpb.FlexRect{}),
			wantErr: sketcherror.ErrIncompleteConfig,
		},
		{
			name:    "incomplete vertex",
			input:   testutil.MustUnmarshalTextproto(t, `top_left {}`, &sketchpb.FlexRect{}),
			wantErr: sketcherror.ErrIncompleteConfig,
		},
		{
			name: "absolute edges",
			input: testutil.MustUnmarshalTextproto(t, `
top { abs { cm: 3 } }
right { abs { cm: 4 } }
bottom { abs { cm: 5 } }
left { abs { cm: 1 } }
`, &sketchpb.FlexRect{}),
			want: geometry.RectFromCentimeters(1, 3, 4, 5),
		},
		{
			name: "absolute vertices",
			input: testutil.MustUnmarshalTextproto(t, `
top_left { abs { left { cm: 1 } top { cm: 3 } } }
bottom_right { abs { left { cm: 4 } top { cm: 5 } } }
`, &sketchpb.FlexRect{}),
			want: geometry.RectFromCentimeters(1, 3, 4, 5),
		},
		{
			name: "top left with size",
			input: testutil.MustUnmarshalTextproto(t, `
top_left { abs { left { cm: 1 } top { cm: 3 } } }
width { cm: 3 }
height { cm: 2 }
`, &sketchpb.FlexRect{}),
			want: geometry.RectFromCentimeters(1, 3, 4, 5),
		},
		{
			name: "top right with size",
			input: testutil.MustUnmarshalTextproto(t, `
top_right { abs { left { cm: 4 } top { cm: 5 } } }
width { cm: 3 }
height { cm: 2 }
`, &sketchpb.FlexRect{}),
			want: geometry.RectFromCentimeters(1, 5, 4, 7),
		},
		{
			name: "bottom right with size",
			input: testutil.MustUnmarshalTextproto(t, `
bottom_right { abs { left { cm: 4 } top { cm: 5 } } }
width { cm: 3 }
height { cm: 2 }
`, &sketchpb.FlexRect{}),
			want: geometry.RectFromCentimeters(1, 3, 4, 5),
		},
		{
			name: "bottom left with size",
			input: testutil.MustUnmarshalTextproto(t, `
bottom_left { abs { left { cm: 4 } top { cm: 5 } } }
width { cm: 3 }
height { cm: 2 }
`, &sketchpb.FlexRect{}),
			want: geometry.RectFromCentimeters(4, 3, 7, 5),
		},
		{
			name: "top left relative to title",
			input: testutil.MustUnmarshalTextproto(t, `
top_left {
	rel {
		node: "title"
		feature: BOTTOM_LEFT
		offset {
			width: { mm: 5 }
			height { mm: 1 }
		}
	}
}
width { cm: 10 }
height { cm: 3 }
`, &sketchpb.FlexRect{}),
			want: geometry.RectFromCentimeters(3.5, 8.1, 13.5, 11.1),
		},
		{
			name: "bottom left relative to title",
			input: testutil.MustUnmarshalTextproto(t, `
bottom_left {
	rel {
		node: "title"
		feature: BOTTOM_LEFT
	}
}
width { cm: 9 }
height { cm: 7 }
`, &sketchpb.FlexRect{}),
			want: geometry.RectFromCentimeters(3, 1, 12, 8),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := FromProto(tc.input)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if gotRect, err := got.Resolve(&cb); err != nil {
					t.Errorf("Resolve() failed: %v", err)
				} else if diff := cmp.Diff(tc.want, gotRect, geometry.EquateLength()); diff != "" {
					t.Errorf("Resolve() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
