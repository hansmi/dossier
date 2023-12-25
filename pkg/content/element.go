package content

import "github.com/hansmi/dossier/pkg/geometry"

type Element interface {
	// Bounds returns the boundary of the element relative to the page.
	Bounds() geometry.Rect
}

type TextElement interface {
	Element

	// Text returns the complete, unmodified text of the element (including
	// leading or trailing space). Multi-line elements separate lines using
	// newline characters (`\n`).
	Text() string

	// RangeBounds returns the boundary of the characters from start to end. Panics
	// if the range starts or ends outside the element text.
	RangeBounds(start, end int) geometry.Rect
}

type Block interface {
	TextElement

	Kind() Block

	// Lines returns the lines contained within the block.
	Lines() []Line
}

type Line interface {
	TextElement

	Kind() Line
}
