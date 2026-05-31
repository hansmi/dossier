package muparser

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/hansmi/dossier/internal/mutool/stext"
	"github.com/hansmi/dossier/internal/ref"
	"github.com/hansmi/dossier/pkg/content"
	"github.com/hansmi/dossier/pkg/geometry"
)

// Line is a single line of text without newlines.
type Line struct {
	bounds geometry.Rect
	chars  []stext.Char
	text   *string
}

var _ content.Line = (*Line)(nil)

func newLine(m stext.Line) *Line {
	result := &Line{
		bounds: m.BBox,
	}

	for _, span := range m.FontSpans {
		result.chars = append(result.chars, span.Chars...)
	}

	return result
}

func (*Line) Kind() content.Line {
	return nil
}

func (l *Line) Bounds() geometry.Rect {
	return l.bounds
}

func (l *Line) Text() string {
	if l.text == nil {
		var buf strings.Builder

		buf.Grow(len(l.chars))

		for _, c := range l.chars {
			buf.WriteRune(c.C)
		}

		l.text = ref.Ref(buf.String())
	}

	return *l.text
}

// RangeBounds returns the rectangular bounds enclosing the characters between
// the byte offsets start and end.
func (l *Line) RangeBounds(start, end int) geometry.Rect {
	raiseOutOfRange := func() {
		panic(fmt.Errorf("bounds [%d:%d] out of range", start, end))
	}

	if start < 0 || end < start {
		raiseOutOfRange()
	}

	runeStart := -1
	runeEnd := -1
	offset := 0

	// Map byte offsets to rune indices.
	for idx, c := range l.chars {
		size := utf8.RuneLen(c.C)

		// Find the starting boundary.
		if runeStart < 0 && start < offset+size {
			runeStart = idx
		}

		// Find the ending boundary.
		if runeEnd < 0 && end <= offset+size {
			if end == offset {
				runeEnd = idx
			} else {
				runeEnd = idx + 1
			}
		}

		offset += size

		if runeStart >= 0 && runeEnd >= 0 {
			// Both ends have been found.
			break
		}
	}

	if runeStart < 0 {
		// Happens if start equals the total byte length. In Go, slicing
		// exactly at the end (e.g., slice[len:len]) safely returns an empty
		// slice.
		runeStart = len(l.chars)
	}

	// Enforce the upper length boundary.
	if runeEnd < 0 {
		// The requested end exceeded the total bytes. The only valid exception
		// is an empty slice where end == 0.
		if len(l.chars) == 0 && end == 0 {
			runeEnd = 0
		} else {
			raiseOutOfRange()
		}
	}

	chars := l.chars[runeStart:runeEnd]
	bounds := chars[0].Bounds

	if len(chars) > 0 {
		for _, c := range chars[1:] {
			bounds = bounds.Union(c.Bounds)
		}
	}

	return bounds
}
