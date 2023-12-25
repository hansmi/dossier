package geometry

type RectCompareFunc func(a, b Rect) int

// MakeRectRowColumnCompare returns a function comparing the positions of two
// rectangles. Rectangles are organized into lines and columns with the
// orientations controlled by v and h. The comparer's return value is zero if
// the rectangles overlap at least 50% on one of the axes, -1 if a is to the
// left or above b, +1 for the opposite.
func MakeRectRowColumnCompare(v VerticalDirection, h HorizontalDirection) RectCompareFunc {
	vr := map[VerticalDirection]int{
		TopToBottom: 1,
		BottomToTop: -1,
	}[v]

	hr := map[HorizontalDirection]int{
		LeftToRight: 1,
		RightToLeft: -1,
	}[h]

	return func(a, b Rect) int {
		ac := a.Center()
		bc := b.Center()

		if ac.Top < b.Top || bc.Top > a.Bottom {
			return -vr
		}

		if bc.Top < a.Top || ac.Top > b.Bottom {
			return vr
		}

		if ac.Left < b.Left || bc.Left > a.Right {
			return -hr
		}

		if bc.Left < a.Left || ac.Left > b.Right {
			return hr
		}

		return 0
	}
}
