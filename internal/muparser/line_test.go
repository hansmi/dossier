package muparser

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hansmi/dossier/internal/mutool/stext"
	"github.com/hansmi/dossier/pkg/geometry"
)

func TestLine(t *testing.T) {
	type test struct {
		name            string
		line            stext.Line
		wantText        string
		wantBounds      geometry.Rect
		wantBoundsFirst geometry.Rect
		wantBoundsLast  geometry.Rect
	}

	tests := []test{
		{name: "empty"},
	}

	for idx, l := range extractLines(t, "corners.xml") {
		tests = append(tests, test{
			name: fmt.Sprintf("corners_line%d", idx),
			line: l,
			wantText: map[int]string{
				0: "TL",
				1: "BL",
				2: "TR",
				3: "BR",
			}[idx],
			wantBounds: map[int]geometry.Rect{
				0: geometry.RectFromCentimeters(1, 0.93, 1.4, 1.34),
				1: geometry.RectFromCentimeters(1, 7.47, 1.4, 7.88),
				2: geometry.RectFromCentimeters(4.7, 0.93, 5.21, 1.34),
				3: geometry.RectFromCentimeters(4.7, 7.47, 5.21, 7.88),
			}[idx],
			wantBoundsFirst: map[int]geometry.Rect{
				0: geometry.RectFromCentimeters(1, 0.93, 1.22, 1.34),
				1: geometry.RectFromCentimeters(1, 7.47, 1.22, 7.88),
				2: geometry.RectFromCentimeters(4.7, 0.93, 4.96, 1.34),
				3: geometry.RectFromCentimeters(4.7, 7.47, 4.96, 7.88),
			}[idx],
			wantBoundsLast: map[int]geometry.Rect{
				0: geometry.RectFromCentimeters(1.22, 0.93, 1.4, 1.34),
				1: geometry.RectFromCentimeters(1.22, 7.47, 1.4, 7.88),
				2: geometry.RectFromCentimeters(4.96, 0.93, 5.21, 1.34),
				3: geometry.RectFromCentimeters(4.96, 7.47, 5.21, 7.88),
			}[idx],
		})
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			line := newLine(tc.line)

			if diff := cmp.Diff(tc.wantText, line.Text()); diff != "" {
				t.Errorf("Text diff (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tc.wantBounds, line.Bounds(), geometry.EquateLength()); diff != "" {
				t.Errorf("Bounds diff (-want +got):\n%s", diff)
			}

			if line.Text() != "" {
				for _, st := range []struct {
					start int
					end   int
					want  geometry.Rect
				}{
					{0, len(line.Text()), tc.wantBounds},
					{0, 1, tc.wantBoundsFirst},
					{len(line.Text()) - 1, len(line.Text()), tc.wantBoundsLast},
				} {
					got := line.RangeBounds(st.start, st.end)

					if diff := cmp.Diff(st.want, got, geometry.EquateLength()); diff != "" {
						t.Errorf("RangeBounds(%d, %d) diff (-want +got):\n%s", st.start, st.end, diff)
					}
				}
			}

			t.Run("RangeBounds panics", func(t *testing.T) {
				defer func() {
					if err := recover(); err == nil {
						t.Errorf("RangeBounds should panic")
					}
				}()

				line.RangeBounds(0, len(line.Text())+1)
			})
		})
	}
}
