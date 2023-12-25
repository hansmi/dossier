package sketch

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/dossier/internal/testutil"
	"github.com/hansmi/dossier/proto/sketchpb"
)

func TestDetermineNodeOrder(t *testing.T) {
	for _, tc := range []struct {
		name    string
		pbnodes []*sketchpb.Node
		want    []int
		wantErr error
	}{
		{name: "empty"},
		{
			name: "one node",
			pbnodes: []*sketchpb.Node{
				testutil.MustUnmarshalTextproto(t, `
name: "only"
search_areas {
  top { abs {} }
  right { abs {} }
  bottom { abs {} }
  left { abs {} }
}
block_text {}
`, &sketchpb.Node{}),
			},
			want: []int{0},
		},
		{
			name: "first depends on second",
			pbnodes: []*sketchpb.Node{
				testutil.MustUnmarshalTextproto(t, `
name: "first"
search_areas {
  top { rel { node: "second" feature: BOTTOM_RIGHT } }
  right { abs {} }
  bottom { abs {} }
  left { abs {} }
}
block_text {}
`, &sketchpb.Node{}),
				testutil.MustUnmarshalTextproto(t, `
name: "second"
search_areas {
  top { abs {} }
  right { abs {} }
  bottom { abs {} }
  left { abs {} }
}
block_text {}
`, &sketchpb.Node{}),
			},
			want: []int{1, 0},
		},
		{
			name: "three nodes",
			pbnodes: []*sketchpb.Node{
				testutil.MustUnmarshalTextproto(t, `
name: "text"
search_areas {
  top { rel { node: "subtitle" feature: BOTTOM_LEFT } }
  right { abs {} }
  bottom { abs {} }
  left { abs {} }
}
block_text {}
`, &sketchpb.Node{}),
				testutil.MustUnmarshalTextproto(t, `
name: "title"
search_areas {
  top { abs {} }
  right { abs {} }
  bottom { abs {} }
  left { abs {} }
}
block_text {}
`, &sketchpb.Node{}),
				testutil.MustUnmarshalTextproto(t, `
name: "subtitle"
search_areas {
  top { rel { node: "title" feature: BOTTOM_LEFT } }
  right { abs {} }
  bottom { abs {} }
  left { abs {} }
}
block_text {}
`, &sketchpb.Node{}),
			},
			want: []int{1, 2, 0},
		},
		{
			name: "recursive reference",
			pbnodes: []*sketchpb.Node{
				testutil.MustUnmarshalTextproto(t, `
name: "first"
search_areas {
  top { abs {} }
  right { rel { node: "second" } }
  bottom { abs {} }
  left { abs {} }
}
block_text {}
`, &sketchpb.Node{}),
				testutil.MustUnmarshalTextproto(t, `
name: "second"
search_areas {
  top { abs {} }
  right { abs {} }
  bottom { rel { node: "first" } }
  left { abs {} }
}
block_text {}
`, &sketchpb.Node{}),
			},
			wantErr: ErrBadConfig,
		},
		{
			name: "node depends on itself",
			pbnodes: []*sketchpb.Node{
				testutil.MustUnmarshalTextproto(t, `
name: "foo"
search_areas {
  top { abs {} }
  right { rel { node: "foo" } }
  bottom { abs {} }
  left { abs {} }
}
block_text {}
`, &sketchpb.Node{}),
			},
			wantErr: ErrBadConfig,
		},
		{
			name: "reused name",
			pbnodes: []*sketchpb.Node{
				testutil.MustUnmarshalTextproto(t, `
name: "aaa"
search_areas {
  top { abs {} }
  right { abs {} }
  bottom { abs {} }
  left { abs {} }
}
block_text {}
`, &sketchpb.Node{}),
				testutil.MustUnmarshalTextproto(t, `
name: "aaa"
search_areas {
  top { abs {} }
  right { abs {} }
  bottom { abs {} }
  left { abs {} }
}
block_text {}
`, &sketchpb.Node{}),
			},
			wantErr: ErrBadConfig,
		},
		{
			name: "reference not found",
			pbnodes: []*sketchpb.Node{
				testutil.MustUnmarshalTextproto(t, `
name: "first"
search_areas {
  top { abs {} }
  right { abs {} }
  bottom { abs {} }
  left { rel { node: "missing" } }
}
block_text {}
`, &sketchpb.Node{}),
			},
			wantErr: ErrBadConfig,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var nodes []*sketchNode

			for _, pb := range tc.pbnodes {
				n, err := sketchNodeFromProto(pb)
				if err != nil {
					t.Fatal(err)
				}

				nodes = append(nodes, n)
			}

			got, err := determineNodeOrder(nodes)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("determineNodeOrder() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
