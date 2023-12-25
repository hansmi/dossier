package sketchiter

import (
	"context"
	"errors"

	"github.com/hansmi/dossier"
	"github.com/hansmi/dossier/pkg/pagerange"
	"github.com/hansmi/dossier/pkg/sketch"
)

var Done = errors.New("iteration done")

// PageIter iterates evaluates a sketch for each page within a document.
type PageIter struct {
	s    *sketch.Sketch
	doc  *dossier.Document
	cur  int
	done bool
}

func NewPageIter(s *sketch.Sketch, doc *dossier.Document) *PageIter {
	return &PageIter{
		s:   s,
		doc: doc,
		cur: 1,
	}
}

// Next returns the current page and moves the internal position to the next.
// [Done] is returned on all calls after the last page has been returned
// previously.
func (it *PageIter) Next(ctx context.Context) (*sketch.PageReport, error) {
	if !it.done {
		r, err := pagerange.Single(it.cur)
		if err != nil {
			return nil, err
		}

		pages, err := it.doc.ParsePages(ctx, r)
		if err != nil {
			return nil, err
		}

		if len(pages) > 0 {
			report, err := it.s.AnalyzePage(pages[0])
			if err != nil {
				return nil, err
			}

			it.cur++

			return report, nil
		}

		it.done = true
	}

	return nil, Done
}
