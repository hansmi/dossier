package main

import (
	"bytes"
	"context"
	"image/png"
	"mime"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hansmi/dossier"
	"github.com/hansmi/dossier/internal/testfiles"
	"github.com/hansmi/dossier/internal/testutil"
	"github.com/hansmi/dossier/pkg/pagerange"
	"github.com/hansmi/dossier/pkg/renderformat"
)

const envRequiredIntegrationTest = "DOSSIER_INTEGRATION_TEST_REQUIRED"

func checkRequirements(t *testing.T, pf *dossier.MuPdfParserFactory) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	if err := pf.Check(ctx); err != nil {
		fn := t.Fatalf

		if os.Getenv(envRequiredIntegrationTest) == "" {
			fn = t.Skipf
		}

		fn("Requirements check failed: %v", err)
	}
}

func TestIntegration(t *testing.T) {
	var parser dossier.MuPdfParserFactory

	checkRequirements(t, &parser)

	for _, tc := range []struct {
		name      string
		wantPages int
	}{
		{
			name:      "corners.pdf",
			wantPages: 1,
		},
		{
			name:      "corners-cropbox.pdf",
			wantPages: 1,
		},
		{
			name:      "lorem-mixed.pdf",
			wantPages: 1,
		},
		{
			name:      "multipage.pdf",
			wantPages: 3,
		},
		{
			name:      "acme-invoice-11321-19.pdf",
			wantPages: 1,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), tc.name)

			if content, err := testfiles.All.ReadFile(tc.name); err != nil {
				t.Errorf("ReadFile() failed: %v", err)
			} else {
				testutil.MustWriteFile(t, path, content)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			t.Cleanup(cancel)

			doc := dossier.NewDocument(path, dossier.WithDocumentParserFactory(parser.Create))

			if fp, err := doc.Fingerprint(); err != nil {
				t.Errorf("Fingerprint() failed: %v", err)
			} else if len(fp) < 30 {
				t.Errorf("Fingerprint() returned short value: %q", fp)
			}

			if contentType, err := doc.ContentType(); err != nil {
				t.Errorf("ContentType() failed: %v", err)
			} else if _, _, err := mime.ParseMediaType(contentType); err != nil {
				t.Errorf("ParseMediaType(%q) failed: %v", contentType, err)
			}

			if err := doc.Validate(ctx); err != nil {
				t.Errorf("Validate() failed: %v", err)
			}

			pages, err := doc.ParsePages(ctx, pagerange.All)
			if err != nil {
				t.Errorf("ParsePages() failed: %v", err)
			}

			if diff := cmp.Diff(tc.wantPages, len(pages)); diff != "" {
				t.Errorf("Page count diff (-want +got):\n%s", diff)
			}

			var buf bytes.Buffer

			for pageNum := 1; pageNum <= len(pages); pageNum++ {
				buf.Reset()

				if err := doc.RenderPageUsing(ctx, pageNum, &renderformat.PNG{
					Width:  400,
					Output: &buf,
				}); err != nil {
					t.Errorf("RenderPage() failed: %v", err)
				}

				if img, err := png.Decode(&buf); err != nil {
					t.Errorf("Decoding image failed: %v", err)
				} else if bounds := img.Bounds(); bounds.Dx() < 10 || bounds.Dy() < 10 {
					t.Errorf("Image smaller than 10x10: %v", bounds)
				}
			}
		})
	}
}
