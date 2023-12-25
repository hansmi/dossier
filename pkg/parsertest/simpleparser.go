package parsertest

import (
	"context"

	"github.com/hansmi/dossier/pkg/content"
	"github.com/hansmi/dossier/pkg/pagerange"
	"github.com/hansmi/dossier/pkg/renderformat"
)

type SimpleParser struct {
	Pages []content.Page

	ValidateErr error
	ParseErr    error
	RenderErr   error
}

func (p *SimpleParser) Validate(_ context.Context) error {
	return p.ValidateErr
}

func (p *SimpleParser) ParsePages(_ context.Context, r pagerange.Range) ([]content.Page, error) {
	if p.ParseErr != nil {
		return nil, p.ParseErr
	}

	if len(p.Pages) == 0 {
		return nil, nil
	}

	var lower, upper int

	if r.Lower == pagerange.Last {
		lower = len(p.Pages) - 1
		upper = len(p.Pages)
	} else {
		lower = r.Lower - 1

		if r.Upper == pagerange.Last {
			upper = len(p.Pages)
		} else {
			upper = r.Upper
		}

		if upper > len(p.Pages) {
			upper = len(p.Pages)
		}
	}

	if lower > len(p.Pages) {
		return nil, nil
	}

	if lower < 0 {
		lower = 0
	}

	return p.Pages[lower:upper], nil
}

func (p *SimpleParser) RenderPage(_ context.Context, num int, r renderformat.Renderer) error {
	return p.RenderErr
}
