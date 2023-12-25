package geometry

import (
	"fmt"

	"github.com/hansmi/dossier/proto/geometrypb"
)

type Size struct {
	Width, Height Length
}

var _ fmt.Stringer = Size{}

func SizeFromProto(pb *geometrypb.Size) (Size, error) {
	var result Size
	var err error

	if pb != nil {
		err = nextLengthFromProto(err, &result.Width, pb.Width)
		err = nextLengthFromProto(err, &result.Height, pb.Height)
	}

	return result, err
}

func (s Size) String() string {
	return fmt.Sprintf("(%s, %s)", s.Width.String(), s.Height.String())
}

func (s Size) AsProto(unit LengthUnit) *geometrypb.Size {
	return &geometrypb.Size{
		Width:  s.Width.AsProto(unit),
		Height: s.Height.AsProto(unit),
	}
}
