package cliutil

import (
	"flag"
	"fmt"
	"strings"

	"github.com/hansmi/dossier/pkg/geometry"
)

type LengthUnitVar struct {
	p *geometry.LengthUnit
}

var _ flag.Getter = (*LengthUnitVar)(nil)

func NewLengthUnitVar(p *geometry.LengthUnit, def geometry.LengthUnit) *LengthUnitVar {
	*p = def

	return &LengthUnitVar{p}
}

func (v *LengthUnitVar) Usage(usage string) string {
	var names []string

	geometry.VisitLengthUnits(func(u geometry.LengthUnit) {
		names = append(names, u.Name())
	})

	return fmt.Sprintf("%s\n(supported: %s)", strings.TrimSpace(usage), names)
}

func (v *LengthUnitVar) String() string {
	if v.p == nil || *v.p == nil {
		return "<nil>"
	}

	return (*v.p).Name()
}

func (v *LengthUnitVar) Get() any {
	return *v.p
}

func (v *LengthUnitVar) Set(s string) error {
	var unit geometry.LengthUnit

	geometry.VisitLengthUnits(func(u geometry.LengthUnit) {
		if u.Name() == s {
			unit = u
		}
	})

	if unit == nil {
		return fmt.Errorf("unsupported length unit %q", s)
	}

	*v.p = unit

	return nil
}
