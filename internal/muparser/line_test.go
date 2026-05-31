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

	for idx, l := range extractLines(t, "unicode1.xml") {
		tests = append(tests, test{
			name: fmt.Sprintf("unicode1_line%d", idx),
			line: l,
			wantText: map[int]string{
				0:  "Iñtërnâtiônàlizætiøn",
				1:  "Heizölrückstoßabdämpfung",
				2:  "Quizdeltagerne spiste jordbær med fløde, mens cirkusklovnen Wolther spillede",
				3:  "på xylofon.",
				4:  "U+002D = \u002D",
				5:  "U+2010 = \u2010",
				6:  "U+2013 = \u2013",
				7:  "U+2014 = \u2014",
				8:  "U+2015 = \u2015",
				9:  "U+2053 = \u2053",
				10: "U+1F609 = \U0001F609",
				11: "U+1F641 = \U0001F641",
				12: "U+1F642 = \U0001F642",
			}[idx],
			wantBounds: map[int]geometry.Rect{
				0:  geometry.RectFromCentimeters(2, 2.06, 6.4, 2.41),
				1:  geometry.RectFromCentimeters(2, 3.06, 7.95, 3.48),
				2:  geometry.RectFromCentimeters(2, 3.98, 19, 4.47),
				3:  geometry.RectFromCentimeters(2, 4.49, 4.35, 4.96),
				4:  geometry.RectFromCentimeters(2, 5.54, 4.63, 5.86),
				5:  geometry.RectFromCentimeters(2, 6.52, 4.56, 6.84),
				6:  geometry.RectFromCentimeters(2, 7.51, 4.63, 7.83),
				7:  geometry.RectFromCentimeters(2, 8.5, 4.84, 8.82),
				8:  geometry.RectFromCentimeters(2, 9.49, 4.84, 9.9),
				9:  geometry.RectFromCentimeters(2, 10.5, 4.84, 10.8),
				10: geometry.RectFromCentimeters(2, 11.4, 5.24, 11.9),
				11: geometry.RectFromCentimeters(2, 12.4, 5.24, 12.9),
				12: geometry.RectFromCentimeters(2, 13.4, 5.24, 13.9),
			}[idx],
			wantBoundsFirst: map[int]geometry.Rect{
				0:  geometry.RectFromCentimeters(2, 2.06, 2.17, 2.39),
				1:  geometry.RectFromCentimeters(2, 3.06, 2.37, 3.48),
				2:  geometry.RectFromCentimeters(2, 3.98, 2.35, 4.47),
				3:  geometry.RectFromCentimeters(2, 4.64, 2.27, 4.96),
				4:  geometry.RectFromCentimeters(2, 5.54, 2.36, 5.86),
				5:  geometry.RectFromCentimeters(2, 6.52, 2.36, 6.84),
				6:  geometry.RectFromCentimeters(2, 7.51, 2.36, 7.83),
				7:  geometry.RectFromCentimeters(2, 8.5, 2.36, 8.82),
				8:  geometry.RectFromCentimeters(2, 9.49, 2.36, 9.9),
				9:  geometry.RectFromCentimeters(2, 10.5, 2.36, 10.8),
				10: geometry.RectFromCentimeters(2, 11.4, 2.36, 11.8),
				11: geometry.RectFromCentimeters(2, 12.4, 2.36, 12.8),
				12: geometry.RectFromCentimeters(2, 13.4, 2.36, 13.8),
			}[idx],
			wantBoundsLast: map[int]geometry.Rect{
				0:  geometry.RectFromCentimeters(6.12, 2.17, 6.4, 2.39),
				1:  geometry.RectFromCentimeters(7.68, 3.06, 7.95, 3.48),
				2:  geometry.RectFromCentimeters(18.7, 4.14, 19, 4.38),
				3:  geometry.RectFromCentimeters(4.21, 4.81, 4.35, 4.86),
				4:  geometry.RectFromCentimeters(4.48, 5.72, 4.63, 5.75),
				5:  geometry.RectFromCentimeters(4.41, 6.71, 4.56, 6.74),
				6:  geometry.RectFromCentimeters(4.41, 7.7, 4.63, 7.73),
				7:  geometry.RectFromCentimeters(4.41, 8.69, 4.84, 8.72),
				8:  geometry.RectFromCentimeters(4.41, 9.68, 4.84, 9.7),
				9:  geometry.RectFromCentimeters(4.41, 10.6, 4.84, 10.7),
				10: geometry.RectFromCentimeters(4.71, 11.4, 5.24, 11.8),
				11: geometry.RectFromCentimeters(4.71, 12.4, 5.24, 12.8),
				12: geometry.RectFromCentimeters(4.71, 13.4, 5.24, 13.8),
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
