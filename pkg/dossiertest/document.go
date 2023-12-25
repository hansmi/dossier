package dossiertest

import (
	"os"
	"testing"

	"github.com/hansmi/dossier"
	"github.com/hansmi/dossier/pkg/parsertest"
)

func NewEmptyDocument(t *testing.T) *dossier.Document {
	t.Helper()

	return dossier.NewDocument(os.DevNull,
		dossier.WithStaticDocumentParser(&parsertest.SimpleParser{}),
	)
}
