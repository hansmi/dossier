package mutool

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hansmi/dossier/internal/mutool/stext"
	"github.com/hansmi/dossier/pkg/pagerange"
	"github.com/hansmi/dossier/pkg/renderformat"
	"go.uber.org/multierr"
)

type Options struct {
	// Command and optional arguments for running mutool. Defaults to "mutool".
	MutoolCommand []string

	// Command and optional arguments for running xmllint. Defaults to
	// "xmllint".
	XmllintCommand []string
}

type Wrapper struct {
	mutool  mutoolInvoker
	xmllint xmllintInvoker
}

func New(opts Options) *Wrapper {
	if len(opts.MutoolCommand) == 0 {
		opts.MutoolCommand = []string{"mutool"}
	}

	if len(opts.XmllintCommand) == 0 {
		opts.XmllintCommand = []string{"xmllint"}
	}

	return &Wrapper{
		mutool:  &mutoolCommand{opts.MutoolCommand},
		xmllint: &xmllintCommand{opts.XmllintCommand},
	}
}

func (w *Wrapper) CheckCommands(ctx context.Context) error {
	err := multierr.Combine(
		w.mutool.CheckCommand(ctx),
		w.xmllint.CheckCommand(ctx),
	)
	if err != nil {
		return fmt.Errorf("commands are not working properly: %w", err)
	}

	return nil
}

// Validate checks whether the given file can be parsed successfully.
func (w *Wrapper) Validate(ctx context.Context, path string) error {
	// TODO: Run command with known output and check for that (LC_ALL may have
	// to be set).
	err := w.mutool.Show(ctx, showArgs{input: path})

	if err != nil {
		return fmt.Errorf("validating document %q: %w", path, err)
	}

	return nil
}

func (w *Wrapper) structuredText(ctx context.Context, path string, r pagerange.Range) (_ *stext.Document, err error) {
	tmpdir, tmpdirCleanup, err := withTempdir()
	if err != nil {
		return nil, err
	}

	defer multierr.AppendFunc(&err, tmpdirCleanup)

	stextFile := filepath.Join(tmpdir, "stext.xml")
	recoveredFile := filepath.Join(tmpdir, "recovered.xml")

	if err := w.mutool.Draw(ctx, drawArgs{
		input:     path,
		pageRange: formatPageRange(r),

		output: stextFile,
		format: "stext",
	}); err != nil {
		return nil, err
	}

	var syntaxErr *xml.SyntaxError

	if doc, err := stext.DocumentFromXMLFile(stextFile); err == nil {
		return doc, nil
	} else if !errors.As(err, &syntaxErr) {
		return nil, err
	}

	// mutool can produce invalid XML output, e.g. with NUL bytes encoded into
	// attributes. Some of these outputs can be recovered using xmllint.
	if err := w.xmllint.Recover(ctx, recoverArgs{
		input:  stextFile,
		output: recoveredFile,
	}); err != nil {
		return nil, err
	}

	// Try again after repairing errors.
	return stext.DocumentFromXMLFile(recoveredFile)
}

func (w *Wrapper) StructuredText(ctx context.Context, path string, r pagerange.Range) (_ *stext.Document, err error) {
	doc, err := w.structuredText(ctx, path, r)
	if err != nil {
		return nil, fmt.Errorf("extraction of structured text from %q: %w", path, err)
	}

	return doc, nil
}

// Draw produces an image of a single page from a document.
func (w *Wrapper) Draw(ctx context.Context, path string, pageNum int, r renderformat.Renderer) error {
	a := drawArgs{
		input:     path,
		pageRange: strconv.Itoa(pageNum),
		output:    "-",
	}

	switch r := r.(type) {
	case *renderformat.PNG:
		a.format = "png"
		a.width = r.Width
		a.height = r.Height
		a.stdout = r.Output
	default:
		return fmt.Errorf("%w: render format %q is not supported", os.ErrInvalid, r.String())
	}

	if err := w.mutool.Draw(ctx, a); err != nil {
		return fmt.Errorf("drawing page %d of %q: %w", pageNum, path, err)
	}

	return nil
}
