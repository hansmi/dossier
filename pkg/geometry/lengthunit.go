package geometry

import "github.com/hansmi/dossier/proto/geometrypb"

type LengthUnit interface {
	Name() string
	lengthToString(Length) string
	lengthToProto(Length) *geometrypb.Length
}

// Unit rounding all length values for a friendlier presentation.
type RoundedLength struct {
	// Unit for formatting values. Defaults to [Nearest] if unset.
	Unit LengthUnit

	// All values are rounded to the nearest multiple of this length.
	Nearest Length
}

var _ LengthUnit = (*RoundedLength)(nil)

func (r RoundedLength) unit() LengthUnit {
	if r.Unit == nil {
		return r.Nearest
	}

	return r.Unit
}

func (r RoundedLength) Name() string {
	return r.unit().Name()
}

func (r RoundedLength) round(value Length) Length {
	return value.Round(r.Nearest)
}

func (r RoundedLength) lengthToString(value Length) string {
	return r.unit().lengthToString(r.round(value))
}

func (r RoundedLength) lengthToProto(value Length) *geometrypb.Length {
	return r.unit().lengthToProto(r.round(value))
}
