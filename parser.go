package dossier

import (
	"context"

	"github.com/hansmi/dossier/pkg/content"
	"github.com/hansmi/dossier/pkg/pagerange"
	"github.com/hansmi/dossier/pkg/renderformat"
)

type Parser interface {
	// Validate whether the data can be successfully parsed.
	Validate(context.Context) error

	// PageCount() int

	ParsePages(context.Context, pagerange.Range) ([]content.Page, error)

	RenderPage(context.Context, int, renderformat.Renderer) error
}
