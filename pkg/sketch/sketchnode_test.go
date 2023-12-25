package sketch

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/dossier"
	"github.com/hansmi/dossier/internal/sketcherror"
	"github.com/hansmi/dossier/internal/testutil"
	"github.com/hansmi/dossier/pkg/geometry"
	"github.com/hansmi/dossier/pkg/pagerange"
	"github.com/hansmi/dossier/proto/reportpb"
	"github.com/hansmi/dossier/proto/sketchpb"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestSketchNodeFromProto(t *testing.T) {
	for _, tc := range []struct {
		name     string
		input    *sketchpb.Node
		wantErr  error
		wantName string
		wantTags []string
	}{
		{
			name:    "empty",
			input:   &sketchpb.Node{},
			wantErr: sketcherror.ErrBadConfig,
		},
		{
			name: "missing search areas",
			input: testutil.MustUnmarshalTextproto(t, `
name: "nosearch"
block_text {}
`, &sketchpb.Node{}),
			wantErr: sketcherror.ErrIncompleteConfig,
		},
		{
			name: "minimal block",
			input: testutil.MustUnmarshalTextproto(t, `
name: "testblock"
search_areas {
  top { abs {} }
  right { abs {} }
  bottom { abs {} }
  left { abs {} }
}
block_text {}
`, &sketchpb.Node{}),
			wantName: "testblock",
		},
		{
			name: "minimal line",
			input: testutil.MustUnmarshalTextproto(t, `
name: "testline"
search_areas {
  top { abs {} }
  right { abs {} }
  bottom { abs {} }
  left { abs {} }
}
line_text {}
`, &sketchpb.Node{}),
			wantName: "testline",
		},
		{
			name: "sorted tags",
			input: testutil.MustUnmarshalTextproto(t, `
name: "sorttags"
search_areas {
  top { abs {} }
  right { abs {} }
  bottom { abs {} }
  left { abs {} }
}
tags: "ccc"
tags: "aaa"
tags: "bbb"
line_text {}
`, &sketchpb.Node{}),
			wantName: "sorttags",
			wantTags: []string{"aaa", "bbb", "ccc"},
		},
		{
			name: "invalid tags",
			input: testutil.MustUnmarshalTextproto(t, `
name: "tagnode"
search_areas {
  top { abs {} }
  right { abs {} }
  bottom { abs {} }
  left { abs {} }
}
tags: "dup"
tags: "key=value"
tags: ""
tags: "dup"
line_text {}
`, &sketchpb.Node{}),
			wantName: "tagnode",
			wantTags: []string{
				"",
				"dup",
				"key=value",
			},
			wantErr: ErrBadConfig,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := sketchNodeFromProto(tc.input)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.wantName, got.name, geometry.EquateLength()); diff != "" {
					t.Errorf("Node name diff (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tc.wantTags, got.tags, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("Node tags diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestSketchNodeFeaturePosition(t *testing.T) {
	for _, tc := range []struct {
		name    string
		input   *sketchpb.Node
		bounds  geometry.Rect
		feature sketchpb.NodeFeature
		want    map[sketchpb.NodeFeature]geometry.Point
		wantErr error
	}{
		{
			name: "feature unavailable",
			input: testutil.MustUnmarshalTextproto(t, `
name: "test26004"
search_areas {
	top { abs { cm: 5 } }
	left { abs { cm: 2 } }
	width { cm: 11 }
	height { cm: 22 }
}
block_text {}
`, &sketchpb.Node{}),
			want: map[sketchpb.NodeFeature]geometry.Point{
				sketchpb.NodeFeature_NODE_FEATURE_UNSPECIFIED: {},
			},
			wantErr: ErrNodeFeatureUnavailable,
		},
		{
			name: "",
			input: testutil.MustUnmarshalTextproto(t, `
name: "test24679"
search_areas {
	top { abs { cm: 10 } }
	left { abs { cm: 5 } }
	width { cm: 22 }
	height { cm: 33 }
}
block_text {}
`, &sketchpb.Node{}),
			bounds: geometry.RectFromCentimeters(10, 15, 21, 31),
			want: map[sketchpb.NodeFeature]geometry.Point{
				sketchpb.NodeFeature_TOP_LEFT: {
					Left: 10 * geometry.Cm,
					Top:  15 * geometry.Cm,
				},
				sketchpb.NodeFeature_TOP_RIGHT: {
					Left: 21 * geometry.Cm,
					Top:  15 * geometry.Cm,
				},
				sketchpb.NodeFeature_BOTTOM_LEFT: {
					Left: 10 * geometry.Cm,
					Top:  31 * geometry.Cm,
				},
				sketchpb.NodeFeature_BOTTOM_RIGHT: {
					Left: 21 * geometry.Cm,
					Top:  31 * geometry.Cm,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			n, err := sketchNodeFromProto(tc.input)
			if err != nil {
				t.Errorf("sketchNodeFromProto() failed: %v", err)
			}

			for feature, want := range tc.want {
				got, err := n.featurePosition(tc.bounds, feature)

				if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
					t.Errorf("Error diff (-want +got):\n%s", diff)
				}

				if err == nil {
					if diff := cmp.Diff(want, got, geometry.EquateLength()); diff != "" {
						t.Errorf("featurePosition() diff (-want +got):\n%s", diff)
					}
				}
			}
		})
	}
}

type fakeSearchCallbacks struct {
	doc                 *dossier.Document
	nodeFeaturePosition func(string, sketchpb.NodeFeature) (geometry.Point, error)
}

func (c *fakeSearchCallbacks) VisitElementsIntersecting(bounds geometry.Rect, visitor dossier.PageElementVisitorFunc) error {
	if c.doc == nil {
		return nil
	}

	pages, err := c.doc.ParsePages(context.Background(), pagerange.MustSingle(1))
	if err != nil {
		return err
	}

	err = pages[0].VisitElementsIntersecting(bounds, visitor)

	if errors.Is(err, dossier.ErrStopVisitation) {
		err = nil
	}

	return err
}

func (c *fakeSearchCallbacks) NodeFeaturePosition(name string, feature sketchpb.NodeFeature) (geometry.Point, error) {
	if c.nodeFeaturePosition == nil {
		return geometry.Point{}, ErrNodePositionUnknown
	}

	return c.nodeFeaturePosition(name, feature)
}

func TestSketchNodeSearch(t *testing.T) {
	unit := geometry.RoundedLength{
		Unit:    geometry.Pt,
		Nearest: geometry.Pt,
	}

	for _, tc := range []struct {
		name    string
		input   *sketchpb.Node
		cb      sketchNodeSearchCallbacks
		wantErr error
		want    *reportpb.Node
	}{
		{
			name: "nothing found",
			input: testutil.MustUnmarshalTextproto(t, `
name: "testblock"
search_areas {
  top { abs {} }
  right { abs {} }
  bottom { abs {} }
  left { abs {} }
}
block_text {}
`, &sketchpb.Node{}),
			want: testutil.MustUnmarshalTextproto(t, `
name: "testblock"
search_areas {
  top { pt: 0 }
  right { pt: 0 }
  bottom { pt: 0 }
  left { pt: 0 }
}
`, &reportpb.Node{}),
		},
		{
			name: "block match found",
			cb: &fakeSearchCallbacks{
				doc: readTestDocument(t, "multipage.xml"),
			},
			input: testutil.MustUnmarshalTextproto(t, `
name: "testblock"
search_areas {
  top { abs { cm: 1 } }
  right { abs { cm: 11 } }
  bottom { abs { cm: 8 } }
  left { abs { cm: 4 } }
}
search_areas {
  top { abs { cm: 5 } }
  right { abs { cm: 19 } }
  bottom { abs { cm: 20 } }
  left { abs { cm: 2 } }
}
search_areas {
  top { abs { cm: 31 } }
  right { abs { cm: 31 } }
  bottom { abs { cm: 38 } }
  left { abs { cm: 34 } }
}
block_text {
  regex: "(?i)(amet).*"
}
`, &sketchpb.Node{}),
			want: testutil.MustUnmarshalTextproto(t, `
name: "testblock"
valid: true
bounds {
  top { pt: 237 }
  right { pt: 345 }
  bottom { pt: 250 }
  left { pt: 74 }
}
search_areas {
  top { pt: 28 }
  right { pt: 312 }
  bottom { pt: 227 }
  left { pt: 113 }
}
search_areas {
  top: { pt: 142 }
  right: { pt: 539 }
  bottom: { pt: 567 }
  left: { pt: 57 }
}
search_areas {
  top: { pt: 879 }
  right: { pt: 964 }
  bottom: { pt: 1077 }
  left: { pt: 879 }
}
text: {
  value: "Lorem ipsum dolor sit amet, consectetur adipisici elit"
}
text_match_groups {
  start: 22
  end: 54
  text: "amet, consectetur adipisici elit"
}
text_match_groups {
  start: 22
  end: 26
  text: "amet"
}
`, &reportpb.Node{}),
		},
		{
			name: "bounds from match",
			cb: &fakeSearchCallbacks{
				doc: readTestDocument(t, "lorem-mixed.xml"),
			},
			input: testutil.MustUnmarshalTextproto(t, `
name: "x"
search_areas {
  top { abs { cm: 1 } }
  right { abs { cm: 30 } }
  bottom { abs { cm: 20 } }
  left { abs { cm: 2 } }
}
line_text {
  regex: "(?im)\\bsanctus\\b"
  bounds_from_match: true
}
`, &sketchpb.Node{}),
			want: testutil.MustUnmarshalTextproto(t, `
name: "x"
valid: true
bounds {
  top: { pt: 130 }
  right: { pt: 384 }
  bottom: { pt: 143 }
  left: { pt: 349 }
}
search_areas {
  top: { pt: 28 }
  right: { pt: 850 }
  bottom: { pt: 567 }
  left: { pt: 57 }
}
text: {
  value: "takimata sanctus est Lorem ipsum dolor sit "
}
text_match_groups {
  start: 9
  end: 16
  text: "sanctus"
}
`, &reportpb.Node{}),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			node, err := sketchNodeFromProto(tc.input)
			if err != nil {
				t.Fatalf("sketchNodeFromProto() failed: %v", err)
			}

			if tc.cb == nil {
				tc.cb = &fakeSearchCallbacks{}
			}

			got, err := node.search(tc.cb)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got.AsProto(unit), protocmp.Transform()); diff != "" {
					t.Errorf("search() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
