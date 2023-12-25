package renderformat

import "io"

type PNG struct {
	Width  int
	Height int
	Output io.Writer
}

var _ Renderer = (*PNG)(nil)

func (r *PNG) String() string {
	return "PNG"
}
