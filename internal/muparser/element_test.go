package muparser

import (
	"testing"

	"github.com/hansmi/dossier/internal/mutool/stext"
	"github.com/hansmi/dossier/internal/testfiles"
)

func loadTestDocument(t *testing.T, name string) *stext.Document {
	t.Helper()

	r, err := testfiles.All.Open(name)
	if err != nil {
		t.Fatal(err)
	}

	defer r.Close()

	doc, err := stext.DocumentFromXML(r)
	if err != nil {
		t.Fatal(err)
	}

	return doc
}

func extractContent(t *testing.T, name string) ([]stext.Block, []stext.Line) {
	t.Helper()

	doc := loadTestDocument(t, name)

	var blocks []stext.Block
	var lines []stext.Line

	for _, p := range doc.Pages {
		blocks = append(blocks, p.Blocks...)

		for _, b := range p.Blocks {
			lines = append(lines, b.Lines...)
		}
	}

	return blocks, lines
}

func extractBlocks(t *testing.T, name string) []stext.Block {
	blocks, _ := extractContent(t, name)
	return blocks
}

func extractLines(t *testing.T, name string) []stext.Line {
	_, lines := extractContent(t, name)
	return lines
}
