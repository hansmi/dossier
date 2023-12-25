package geometry

import (
	"errors"
	"fmt"

	"github.com/hansmi/dossier/proto/geometrypb"
)

var ErrRectInvalid = errors.New("rectangle is invalid")

type Rect struct {
	// X-coordinate of the upper-left corner.
	Left Length

	// Y-coordinate of the upper-left corner.
	Top Length

	// X-coordinate of the lower-right corner.
	Right Length

	// Y-coordinate of the lower-right corner.
	Bottom Length
}

var _ fmt.Stringer = Rect{}

// RectFromPoints contructs a rectangle from coordinates in points.
func RectFromPoints(left, top, right, bottom float64) Rect {
	return Rect{
		Left:   Pt.Mul(left),
		Top:    Pt.Mul(top),
		Right:  Pt.Mul(right),
		Bottom: Pt.Mul(bottom),
	}
}

// RectFromCentimeters constructs a rectangle from coordinates in centimeters.
func RectFromCentimeters(left, top, right, bottom float64) Rect {
	return Rect{
		Left:   Cm.Mul(left),
		Top:    Cm.Mul(top),
		Right:  Cm.Mul(right),
		Bottom: Cm.Mul(bottom),
	}
}

func RectFromXYWH(left, top, width, height Length) Rect {
	return Rect{
		Left:   left,
		Top:    top,
		Right:  left + width,
		Bottom: top + height,
	}
}

func RectFromProto(pb *geometrypb.Rect) (Rect, error) {
	var result Rect
	var err error

	if pb != nil {
		err = nextLengthFromProto(err, &result.Top, pb.Top)
		err = nextLengthFromProto(err, &result.Right, pb.Right)
		err = nextLengthFromProto(err, &result.Bottom, pb.Bottom)
		err = nextLengthFromProto(err, &result.Left, pb.Left)
	}

	return result, err
}

func (r Rect) String() string {
	return fmt.Sprintf("%s-%s", r.TopLeft(), r.BottomRight())
}

func (r Rect) AsProto(unit LengthUnit) *geometrypb.Rect {
	return &geometrypb.Rect{
		Top:    r.Top.AsProto(unit),
		Right:  r.Right.AsProto(unit),
		Bottom: r.Bottom.AsProto(unit),
		Left:   r.Left.AsProto(unit),
	}
}

// Normalizes the rectangle so both the width and height are increasing.
func (r Rect) Normalize() Rect {
	if r.Right < r.Left {
		r.Left, r.Right = r.Right, r.Left
	}

	if r.Bottom < r.Top {
		r.Top, r.Bottom = r.Bottom, r.Top
	}

	return r
}

func (r Rect) Validate() error {
	if r.Left <= r.Right && r.Top <= r.Bottom {
		return nil
	}

	return fmt.Errorf("%w: %s", ErrRectInvalid, r.String())
}

func (r Rect) Width() Length {
	return r.Right - r.Left
}

func (r Rect) Height() Length {
	return r.Bottom - r.Top
}

func (r Rect) Center() Point {
	return Point{
		Left: r.Left + r.Width()/2,
		Top:  r.Top + r.Height()/2,
	}
}

func (r Rect) TopLeft() Point {
	return Point{Left: r.Left, Top: r.Top}
}

func (r Rect) TopRight() Point {
	return Point{Left: r.Right, Top: r.Top}
}

func (r Rect) BottomLeft() Point {
	return Point{Left: r.Left, Top: r.Bottom}
}

func (r Rect) BottomRight() Point {
	return Point{Left: r.Right, Top: r.Bottom}
}

func (r Rect) IsEmpty() bool {
	return r.Left == r.Right || r.Top == r.Bottom
}

// Test whether the rectangle is fully contained within another.
func (r Rect) Inside(other Rect) bool {
	return other.Contains(r)
}

// Test whether the rectangle fully contains another.
func (r Rect) Contains(other Rect) bool {
	return (r.Left <= other.Left &&
		other.Left <= other.Right &&
		other.Right <= r.Right &&
		r.Top <= other.Top &&
		other.Top <= other.Bottom &&
		other.Bottom <= r.Bottom)
}

// Union merges two rectangles. The union is the smallest rectangle which
// contains both r and other.
func (r Rect) Union(other Rect) Rect {
	r.Left = r.Left.Min(other.Left)
	r.Top = r.Top.Min(other.Top)
	r.Right = r.Right.Max(other.Right)
	r.Bottom = r.Bottom.Max(other.Bottom)
	return r
}
