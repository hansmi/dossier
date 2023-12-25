package flexrect

import "github.com/hansmi/dossier/pkg/geometry"

type pointDimensionFunc func(geometry.Point) geometry.Length

func pointLeft(p geometry.Point) geometry.Length {
	return p.Left
}

func pointTop(p geometry.Point) geometry.Length {
	return p.Top
}
