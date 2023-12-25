package muparser

import (
	"strings"

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

func (l *Line) RangeBounds(start, end int) geometry.Rect {
	chars := l.chars[start:end]
	bounds := chars[0].Bounds

	if len(chars) > 1 {
		for _, c := range chars[1:] {
			bounds = bounds.Union(c.Bounds)
		}
	}

	return bounds
}
