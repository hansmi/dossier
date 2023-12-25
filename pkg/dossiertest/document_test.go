package dossiertest

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hansmi/dossier/pkg/pagerange"
	"github.com/hansmi/dossier/pkg/renderformat"
)

func TestNewEmptyDocument(t *testing.T) {
	doc := NewEmptyDocument(t)

	if got, err := doc.ContentType(); err != nil {
		t.Errorf("ContentType() failed: %v", err)
	} else if diff := cmp.Diff("text/plain", got); diff != "" {
		t.Errorf("Content-type diff (-want +got):\n%s", diff)
	}

	if err := doc.Validate(context.Background()); err != nil {
		t.Errorf("Validate() failed: %v", err)
	}

	if _, err := doc.ParsePages(context.Background(), pagerange.All); err != nil {
		t.Errorf("ParsePages() failed: %v", err)
	}

	if err := doc.RenderPageUsing(context.Background(), 0, &renderformat.PNG{}); err != nil {
		t.Errorf("RenderPageUsing() failed: %v", err)
	}
}
