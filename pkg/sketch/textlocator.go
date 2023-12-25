package sketch

import (
	"regexp"

	"github.com/hansmi/dossier"
	"github.com/hansmi/dossier/pkg/content"
	"github.com/hansmi/dossier/pkg/geometry"
)

type textLocator struct {
	line            bool
	pattern         *regexp.Regexp
	boundsFromMatch bool
}

func newTextLocatorFromProto(pbnode interface {
	GetRegex() string
	GetBoundsFromMatch() bool
}, line bool) (*textLocator, error) {
	var err error

	l := &textLocator{
		line:            line,
		boundsFromMatch: pbnode.GetBoundsFromMatch(),
	}

	if l.pattern, err = regexp.Compile(pbnode.GetRegex()); err != nil {
		return nil, err
	}

	return l, nil
}

func (l *textLocator) locate(cb documentPage, bounds geometry.Rect) (func(*Node), error) {
	var result func(*Node)

	visit := func(elem content.TextElement) error {
		if !bounds.Contains(elem.Bounds()) {
			return nil
		}

		text := elem.Text()

		if m := evaluateMatch(l.pattern, text); m != nil {
			bounds := elem.Bounds()

			if l.boundsFromMatch {
				g0 := m.MustGroup(0)
				bounds = elem.RangeBounds(g0.Start, g0.End)
			}

			result = func(n *Node) {
				n.valid = true
				n.bounds = bounds
				n.text = &text
				n.textMatch = m
			}

			return dossier.ErrStopVisitation
		}

		return nil
	}

	var visitor dossier.PageElementVisitorFunc

	if l.line {
		visitor = dossier.AsPageElementVisitor(func(elem content.Line) error {
			return visit(elem)
		})
	} else {
		visitor = dossier.AsPageElementVisitor(func(elem content.Block) error {
			return visit(elem)
		})
	}

	if err := cb.VisitElementsIntersecting(bounds, visitor); err != nil {
		return nil, err
	}

	return result, nil
}
