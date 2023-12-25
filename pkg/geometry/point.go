package geometry

import (
	"fmt"

	"github.com/hansmi/dossier/proto/geometrypb"
)

type Point struct {
	// Coordinates
	Left, Top Length
}

var _ fmt.Stringer = Point{}

func PointFromProto(pb *geometrypb.Point) (Point, error) {
	var result Point
	var err error

	if pb != nil {
		err = nextLengthFromProto(err, &result.Left, pb.Left)
		err = nextLengthFromProto(err, &result.Top, pb.Top)
	}

	return result, err
}

func (p Point) String() string {
	return fmt.Sprintf("(%s, %s)", p.Left.String(), p.Top.String())
}

func (p Point) AsProto(unit LengthUnit) *geometrypb.Point {
	return &geometrypb.Point{
		Left: p.Left.AsProto(unit),
		Top:  p.Top.AsProto(unit),
	}
}

func (p Point) Shift(s Size) Point {
	p.Left += s.Width
	p.Top += s.Height

	return p
}
