package muparser

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hansmi/dossier/internal/mutool/stext"
	"github.com/hansmi/dossier/pkg/geometry"
)

func TestBlock(t *testing.T) {
	type rangeBounds struct {
		start, end int
		want       geometry.Rect
	}

	type test struct {
		name            string
		block           stext.Block
		wantText        string
		wantBounds      geometry.Rect
		wantRangeBounds []rangeBounds
	}

	tests := []test{
		{name: "empty"},
	}

	for idx, b := range extractBlocks(t, "lorem-mixed.xml") {
		tests = append(tests, test{
			name:  fmt.Sprintf("lorem_mixed_block%d", idx),
			block: b,
			wantText: map[int]string{
				0: "Lorem ipsum dolor sit amet, consectetur adipisici elit, sed eiusmod tempor incidunt \nut labore et dolore magna aliqua.",
				1: "At vero eos et accusam et justo duo dolores et ea\nrebum. Stet clita kasd gubergren, no sea ",
				2: "takimata sanctus est Lorem ipsum dolor sit \namet.",
			}[idx],
			wantBounds: map[int]geometry.Rect{
				0: geometry.RectFromCentimeters(2, 2, 19, 3.11),
				1: geometry.RectFromCentimeters(2, 4.59, 10.3, 5.5),
				2: geometry.RectFromCentimeters(10.8, 4.59, 18.1, 5.5),
			}[idx],
			wantRangeBounds: map[int][]rangeBounds{
				0: {
					{0, 118, geometry.RectFromCentimeters(2, 2, 19, 3.11)},
					{69, 83, geometry.RectFromCentimeters(14.9, 2, 18.7, 2.62)},
					{89, 99, geometry.RectFromCentimeters(3.02, 2.62, 5.57, 3.11)},
				},
				1: {
					{0, 91, geometry.RectFromCentimeters(2, 4.59, 10.3, 5.5)},
				},
				2: {
					{0, 49, geometry.RectFromCentimeters(10.8, 4.59, 18.1, 5.5)},
				},
			}[idx],
		})
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			block := newBlock(tc.block)

			if diff := cmp.Diff(tc.wantText, block.Text()); diff != "" {
				t.Errorf("Text diff (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tc.wantBounds, block.Bounds(), geometry.EquateLength()); diff != "" {
				t.Errorf("Bounds diff (-want +got):\n%s", diff)
			}

			for _, rb := range tc.wantRangeBounds {
				got := block.RangeBounds(rb.start, rb.end)

				if diff := cmp.Diff(rb.want, got, geometry.EquateLength()); diff != "" {
					t.Errorf("RangeBounds(%d, %d) diff (-want +got):\n%s", rb.start, rb.end, diff)
				}
			}

			t.Run("RangeBounds panics", func(t *testing.T) {
				defer func() {
					if err := recover(); err == nil {
						t.Errorf("RangeBounds should panic")
					}
				}()

				block.RangeBounds(0, len(block.Text())+1)
			})
		})
	}
}
