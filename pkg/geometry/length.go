package geometry

import (
	"bytes"
	"fmt"
	"math"
	"sort"
	"strconv"

	"github.com/hansmi/dossier/proto/geometrypb"
	"golang.org/x/exp/maps"
)

var DefaultLengthUnit = Cm

// Length is a measure of distance in points (1/72th of an inch).
type Length float64

var _ fmt.Stringer = Length(0)
var _ LengthUnit = Length(0)

const (
	Pt Length = 1

	Inch = 72 * Pt
	In   = Inch

	Centimeter = Inch / 2.54
	Cm         = Centimeter

	Millimeter = Centimeter / 10
	Mm         = Millimeter
)

type knownLengthUnit struct {
	name string
	set  func(*geometrypb.Length, Length)
}

var knownLengthUnits = map[Length]*knownLengthUnit{
	Pt: {
		name: "pt",
		set: func(pb *geometrypb.Length, v Length) {
			pb.Value = &geometrypb.Length_Pt{Pt: v.Pt()}
		},
	},
	Millimeter: {
		name: "mm",
		set: func(pb *geometrypb.Length, v Length) {
			pb.Value = &geometrypb.Length_Mm{Mm: v.Mm()}
		},
	},
	Centimeter: {
		name: "cm",
		set: func(pb *geometrypb.Length, v Length) {
			pb.Value = &geometrypb.Length_Cm{Cm: v.Cm()}
		},
	},
	Inch: {
		name: "in",
		set: func(pb *geometrypb.Length, v Length) {
			pb.Value = &geometrypb.Length_In{In: v.Inch()}
		},
	},
}

func getKnownLengthUnit(l Length) *knownLengthUnit {
	u, ok := knownLengthUnits[l]
	if !ok {
		panic(fmt.Sprintf("Unrecognized unit value %[1]s (%[1]f)", l))
	}

	return u
}

func VisitLengthUnits(visitor func(LengthUnit)) {
	units := maps.Keys(knownLengthUnits)

	sort.Slice(units, func(a, b int) bool {
		return knownLengthUnits[units[a]].name < knownLengthUnits[units[b]].name
	})

	for _, u := range units {
		visitor(u)
	}
}

// LengthFromProto constructs a Length from a protobuf message. Nil messages
// are treated as lengths of zero.
func LengthFromProto(l *geometrypb.Length) (Length, error) {
	switch v := l.GetValue().(type) {
	case nil:
		return 0, nil

	case *geometrypb.Length_Cm:
		return Length(v.Cm) * Cm, nil

	case *geometrypb.Length_Mm:
		return Length(v.Mm) * Mm, nil

	case *geometrypb.Length_In:
		return Length(v.In) * In, nil

	case *geometrypb.Length_Pt:
		return Length(v.Pt) * Pt, nil

	default:
		return 0, fmt.Errorf("unknown length unit %T", v)
	}
}

// Name returns the name of the length unit. Panics if the length is not
// a base value.
func (l Length) Name() string {
	return getKnownLengthUnit(l).name
}

func (l Length) lengthToString(value Length) string {
	name := getKnownLengthUnit(l).name

	if value == 0 {
		return "0"
	}

	v := float64(value / l)

	var buf []byte

	// Attempt to find a short representation without scientific notification.
	// There is no way to disable scientific format in combination with
	// trimming trailing zeros.
	for prec := 3; prec <= 6; prec++ {
		buf = strconv.AppendFloat(buf[0:0], v, 'g', prec, 64)

		if bytes.IndexByte(buf, 'e') < 0 {
			break
		}
	}

	return string(append(buf, name...))
}

func (l Length) lengthToProto(value Length) *geometrypb.Length {
	pb := &geometrypb.Length{}

	getKnownLengthUnit(l).set(pb, value)

	return pb
}

func (l Length) AsProto(unit LengthUnit) *geometrypb.Length {
	return unit.lengthToProto(l)
}

// UnitString returns a string representation of the length using the given
// unit.
func (l Length) UnitString(unit LengthUnit) string {
	return unit.lengthToString(l)
}

// String returns a string representation of the length using the default unit
// set in [DefaultLengthUnit].
func (l Length) String() string {
	return l.UnitString(DefaultLengthUnit)
}

// Mm returns the distance in points (1/72 inch).
func (l Length) Pt() float64 {
	return float64(l / Pt)
}

// Mm returns the distance in inch (2.54 cm).
func (l Length) Inch() float64 {
	return float64(l / Inch)
}

// Mm returns the distance in centimeters.
func (l Length) Cm() float64 {
	return float64(l / Cm)
}

// Mm returns the distance in millimeters.
func (l Length) Mm() float64 {
	return float64(l / Mm)
}

// Abs returns the absolute value.
func (l Length) Abs() Length {
	return Length(math.Abs(float64(l)))
}

// Round returns the closest multiple of nearest.
func (l Length) Round(nearest Length) Length {
	return Length(math.Round(float64(l/nearest))) * nearest
}

// Mul returns the scalar product of l and factor.
func (l Length) Mul(factor float64) Length {
	return Length(factor) * l
}

// Min returns the smaller of l and other.
func (l Length) Min(other Length) Length {
	if other < l {
		return other
	}

	return l
}

// Max returns the larger of l and other.
func (l Length) Max(other Length) Length {
	if other > l {
		return other
	}

	return l
}

func nextLengthFromProto(err error, dest *Length, pb *geometrypb.Length) error {
	if err == nil {
		*dest, err = LengthFromProto(pb)
	}

	return err
}
