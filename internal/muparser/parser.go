package muparser

import (
	"context"
	"io"

	"github.com/hansmi/dossier/internal/mutool/stext"
	"github.com/hansmi/dossier/pkg/content"
	"github.com/hansmi/dossier/pkg/pagerange"
	"github.com/hansmi/dossier/pkg/renderformat"
)

// TODO: Support non-PDF file formats.
var SupportedContentTypes = []string{
	"application/pdf",
}

type ToolWrapper interface {
	Validate(context.Context, string) error
	StructuredText(context.Context, string, pagerange.Range) (*stext.Document, error)
	Draw(context.Context, string, int, renderformat.Renderer) error
}

func convertPages(pages []stext.Page) ([]content.Page, error) {
	result := make([]content.Page, len(pages))

	for idx, cur := range pages {
		p, err := newPage(cur)
		if err != nil {
			return nil, err
		}

		result[idx] = p
	}

	return result, nil
}

// ReadPagesFromXML reads structured text from an XML file written by mutool.
func ReadPagesFromXML(r io.Reader) ([]content.Page, error) {
	doc, err := stext.DocumentFromXML(r)
	if err != nil {
		return nil, err
	}

	return convertPages(doc.Pages)
}

type Parser struct {
	path string
	tool ToolWrapper
}

// New creates a new mutool-based parser. mutool requires a regular and
// seekable file.
func New(path string, tool ToolWrapper) *Parser {
	return &Parser{
		path: path,
		tool: tool,
	}
}

func (p *Parser) Validate(ctx context.Context) error {
	return p.tool.Validate(ctx, p.path)
}

// ParsePages uses mutool to parse a file and returns the page contents.
func (p *Parser) ParsePages(ctx context.Context, r pagerange.Range) ([]content.Page, error) {
	doc, err := p.tool.StructuredText(ctx, p.path, r)
	if err != nil {
		return nil, err
	}

	return convertPages(doc.Pages)
}

func (p *Parser) RenderPage(ctx context.Context, pageNum int, r renderformat.Renderer) error {
	return p.tool.Draw(ctx, p.path, pageNum, r)
}
