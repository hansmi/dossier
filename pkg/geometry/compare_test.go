package geometry

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMakeRectRowColumnCompare(t *testing.T) {
	compareRect := MakeRectRowColumnCompare(TopToBottom, LeftToRight)

	for _, tc := range []struct {
		a, b Rect
		want int
	}{
		{},
		{a: RectFromXYWH(-1, -2, 2, 2), b: RectFromXYWH(-1, -2, 2, 2)},
		{a: RectFromXYWH(2, 5, 3, 5), b: RectFromXYWH(2, 4, 4, 7)},
		{
			a:    RectFromXYWH(2, 5, 3, 5),
			b:    RectFromXYWH(6, 9, 4, 7),
			want: -1,
		},
		{
			a:    RectFromXYWH(2, 5, 3, 5),
			b:    RectFromXYWH(6, 1, 4, 7),
			want: +1,
		},
		{
			a:    RectFromXYWH(7, 5, 3, 5),
			b:    RectFromXYWH(2, 4, 4, 7),
			want: +1,
		},
	} {
		t.Run(fmt.Sprint(tc), func(t *testing.T) {
			got := compareRect(tc.a, tc.b)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("CompareRect() diff (-want +got):\n%s", diff)
			}

			gotReverse := compareRect(tc.b, tc.a)

			if diff := cmp.Diff(-tc.want, gotReverse); diff != "" {
				t.Errorf("Reverse CompareRect() diff (-want +got):\n%s", diff)
			}

		})
	}
}
