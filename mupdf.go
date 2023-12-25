package dossier

import (
	"context"

	"github.com/gabriel-vasile/mimetype"
	"github.com/hansmi/dossier/internal/muparser"
	"github.com/hansmi/dossier/internal/mutool"
)

type MuPdfParserFactory struct {
	// Command arguments to invoke MuPDF's "mutool" program. Leave empty to use
	// the default.
	MutoolCommand []string

	// Command arguments to invoke the "xmllint" program. Leave empty to use
	// the default.
	XmllintCommand []string
}

func (f MuPdfParserFactory) makeTool() *mutool.Wrapper {
	return mutool.New(mutool.Options{
		MutoolCommand:  f.MutoolCommand,
		XmllintCommand: f.XmllintCommand,
	})
}

func (f MuPdfParserFactory) Check(ctx context.Context) error {
	return f.makeTool().CheckCommands(ctx)
}

func (f MuPdfParserFactory) Create(path, contentType string) (Parser, error) {
	if !mimetype.EqualsAny(contentType, muparser.SupportedContentTypes...) {
		return nil, nil
	}

	return muparser.New(path, f.makeTool()), nil
}
