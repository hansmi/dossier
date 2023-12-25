package muparser

import (
	"fmt"

	"github.com/hansmi/dossier/internal/mutool/stext"
	"github.com/hansmi/dossier/pkg/content"
	"github.com/hansmi/dossier/pkg/geometry"
)

// Page holds recognized elements on a single page.
type Page struct {
	num      int
	size     geometry.Size
	elements []content.Element
}

var _ content.Page = (*Page)(nil)

func newPage(p stext.Page) (*Page, error) {
	var num int

	if n, err := fmt.Sscanf(p.ID, "page%d\n", &num); !(n == 1 && err == nil) {
		return nil, fmt.Errorf("extracting page number from %q: %w", p.ID, err)
	}

	result := &Page{
		num: num,
		size: geometry.Size{
			Width:  p.Width,
			Height: p.Height,
		},
	}

	for _, block := range p.Blocks {
		b := newBlock(block)

		result.elements = append(result.elements, b)

		for _, l := range b.Lines() {
			result.elements = append(result.elements, l)
		}
	}

	return result, nil
}

func (p *Page) Number() int {
	return p.num
}

func (p *Page) Size() geometry.Size {
	return p.size
}

func (p *Page) Elements() []content.Element {
	return p.elements
}
