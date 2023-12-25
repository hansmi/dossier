package muparser

import (
	"fmt"
	"strings"

	"github.com/hansmi/dossier/internal/mutool/stext"
	"github.com/hansmi/dossier/internal/ref"
	"github.com/hansmi/dossier/pkg/content"
	"github.com/hansmi/dossier/pkg/geometry"
)

// Block is a holder for one or multiple lines of text.
type Block struct {
	bounds geometry.Rect
	lines  []content.Line
	text   *string
}

var _ content.Block = (*Block)(nil)

func newBlock(m stext.Block) *Block {
	b := &Block{
		bounds: m.BBox,
	}

	for _, l := range m.Lines {
		b.lines = append(b.lines, newLine(l))
	}

	return b
}

func (*Block) Kind() content.Block {
	return nil
}

func (b *Block) Bounds() geometry.Rect {
	return b.bounds
}

func (b *Block) Lines() []content.Line {
	return b.lines
}

func (b *Block) Text() string {
	if b.text == nil {
		var buf strings.Builder

		for idx, l := range b.lines {
			if idx > 0 {
				buf.WriteRune('\n')
			}

			buf.WriteString(l.Text())
		}

		b.text = ref.Ref(buf.String())
	}

	return *b.text
}

func (b *Block) RangeBounds(start, end int) geometry.Rect {
	if start < 0 || end < start {
		panic(fmt.Sprintf("range %d-%d is not valid", start, end))
	}

	var started bool
	var bounds geometry.Rect

	pos := 0

	for idx, l := range b.lines {
		text := l.Text()

		if idx > 0 {
			// New line
			pos++
		}

		first := start - pos
		if first < 0 {
			first = 0
		}

		last := end - pos
		if last > len(text) {
			last = len(text)
		}

		if first < last {
			if rbounds := l.RangeBounds(first, last); started {
				bounds = bounds.Union(rbounds)
			} else {
				started = true
				bounds = rbounds
			}
		}

		pos += len(text)

		if pos > end {
			break
		}
	}

	if start > pos || end > pos {
		panic(fmt.Sprintf("range %d-%d exceeds block length of %d", start, end, pos))
	}

	return bounds
}
