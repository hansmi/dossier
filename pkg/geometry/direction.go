package geometry

type HorizontalDirection int
type VerticalDirection int

const (
	LeftToRight HorizontalDirection = iota
	RightToLeft

	TopToBottom VerticalDirection = iota
	BottomToTop
)
