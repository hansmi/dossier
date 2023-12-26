package flexrect

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/dossier/internal/sketcherror"
	"github.com/hansmi/dossier/proto/sketchpb"
	"golang.org/x/exp/slices"
)

func TestNodeFeature(t *testing.T) {
	for _, tc := range []struct {
		name       string
		input      *sketchpb.RelativePosition1D
		wantErr    error
		wantString string
	}{
		{
			name:    "nil",
			wantErr: sketcherror.ErrIncompleteConfig,
		},
		{
			name: "feature unspecified",
			input: &sketchpb.RelativePosition1D{
				Node: "foo",
			},
			wantString: "foo:NODE_FEATURE_UNSPECIFIED",
		},
		{
			name: "bottom right",
			input: &sketchpb.RelativePosition1D{
				Node:    "node",
				Feature: sketchpb.NodeFeature_BOTTOM_RIGHT,
			},
			wantString: "node:BOTTOM_RIGHT",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := newNodeFeature(tc.input)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.wantString, got.String()); diff != "" {
					t.Errorf("String() diff (-want +got):\n%s", diff)
				}
			}

			if gotCompare := got.compare(got); gotCompare != 0 {
				t.Errorf("compare() against itself returned %d, want 0", gotCompare)
			}
		})
	}
}

func TestNodeFeatureCompare(t *testing.T) {
	items := []NodeFeature{
		{},
		{name: "a", feature: sketchpb.NodeFeature_BOTTOM_RIGHT},
		{name: "b"},
		{name: "b", feature: sketchpb.NodeFeature_BOTTOM_LEFT},
		{name: "b", feature: sketchpb.NodeFeature_BOTTOM_RIGHT},
		{name: "c", feature: sketchpb.NodeFeature_TOP_LEFT},
		{name: "c", feature: sketchpb.NodeFeature_TOP_RIGHT},
	}

	for i := 0; i < (2 * len(items)); i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tmp := slices.Clone(items)

			rand.Shuffle(len(tmp), func(a, b int) {
				tmp[a], tmp[b] = tmp[b], tmp[a]
			})

			slices.SortStableFunc(tmp, func(a, b NodeFeature) int {
				return a.compare(b)
			})

			if diff := cmp.Diff(items, tmp, cmp.AllowUnexported(NodeFeature{})); diff != "" {
				t.Errorf("compare() diff (-want +got):\n%s", diff)
			}
		})
	}
}
