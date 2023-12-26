package dossier

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/dossier/internal/muparser"
	"github.com/hansmi/dossier/internal/testfiles"
	"github.com/hansmi/dossier/internal/testutil"
	"github.com/hansmi/dossier/pkg/content"
	"github.com/hansmi/dossier/pkg/geometry"
	"github.com/hansmi/dossier/pkg/pagerange"
	"github.com/hansmi/dossier/pkg/parsertest"
)

func mustReadPagesFromXML(t *testing.T, r io.Reader) []content.Page {
	t.Helper()

	pages, err := muparser.ReadPagesFromXML(r)
	if err != nil {
		t.Fatalf("ReadPagesFromXML() failed: %v", err)
	}

	return pages
}

func mustReadPages(t *testing.T, name string) []content.Page {
	t.Helper()

	f, err := testfiles.All.Open(name)
	if err != nil {
		t.Fatalf("Open() failed: %v", err)
	}

	defer f.Close()

	return mustReadPagesFromXML(t, f)
}

func TestDocumentFingerprint(t *testing.T) {
	for _, tc := range []struct {
		name    string
		path    string
		modify  bool
		wantErr error
	}{
		{
			name:    "missing file",
			path:    filepath.Join(t.TempDir(), "missing"),
			wantErr: os.ErrNotExist,
		},
		{
			name:   "regular file",
			path:   testutil.MustWriteFileString(t, filepath.Join(t.TempDir(), "f"), "content"),
			modify: true,
		},
		{
			name: "directory",
			path: t.TempDir(),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			doc := NewDocument(tc.path)

			fp, err := doc.Fingerprint()

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if len(fp) < 20 {
					t.Errorf("Fingerprint() returned short result: %q", fp)
				}

				if tc.modify {
					testutil.MustWriteFileString(t, tc.path, "modified content")

					if fpMod, err := doc.Fingerprint(); err != nil {
						t.Errorf("Fingerprint() on modified file failed: %v", err)
					} else if fp == fpMod {
						t.Errorf("Fingerprint unchanged after file modification; got %q", fpMod)
					}
				}
			}
		})
	}
}

func TestDocumentContentType(t *testing.T) {
	for _, tc := range []struct {
		name    string
		content string
		want    string
		wantErr error
	}{
		{
			name: "empty",
			want: "text/plain",
		},
		{
			name:    "pdf",
			content: testutil.MustReadFileString(t, testfiles.All, "corners.pdf"),
			want:    "application/pdf",
		},
		{
			name:    "xml",
			content: testutil.MustReadFileString(t, testfiles.All, "corners.xml"),
			want:    "text/xml",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := testutil.MustWriteFileString(t, filepath.Join(t.TempDir(), "content"), tc.content)

			ct, err := NewDocument(path).ContentType()

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil && !mimetype.EqualsAny(ct, tc.want) {
				t.Errorf("ContentType() returned %q, want %q", ct, tc.want)
			}
		})
	}
}

func TestDocumentParsePages(t *testing.T) {
	type pageLookupResult struct {
		num  int
		size geometry.Size
	}

	type pageLookup struct {
		r       pagerange.Range
		want    []pageLookupResult
		wantErr error
	}

	errTest := errors.New("test error")

	for _, tc := range []struct {
		name            string
		p               Parser
		wantValidateErr error
		wantPages       []pageLookup
	}{
		{
			name: "empty",
			p:    &parsertest.SimpleParser{},
			wantPages: []pageLookup{
				{r: pagerange.All},
				{r: pagerange.MustSingle(1)},
				{r: pagerange.MustNew(1, 2)},
			},
		},
		{
			name: "corners",
			p: &parsertest.SimpleParser{
				Pages: mustReadPages(t, "corners.xml"),
			},
			wantPages: func() []pageLookup {
				want := []pageLookupResult{
					{
						num: 1,
						size: geometry.Size{
							Width:  6.2 * geometry.Cm,
							Height: 8.8 * geometry.Cm,
						},
					},
				}

				return []pageLookup{
					{r: pagerange.All, want: want},
					{r: pagerange.MustSingle(1), want: want},
					{r: pagerange.MustSingle(2)},
					{r: pagerange.MustNew(2, 100)},
				}
			}(),
		},
		{
			name: "multipage",
			p: &parsertest.SimpleParser{
				Pages: mustReadPages(t, "multipage.xml"),
			},
			wantPages: func() []pageLookup {
				p1 := pageLookupResult{
					num: 1,
					size: geometry.Size{
						Width:  14.8 * geometry.Cm,
						Height: 21 * geometry.Cm,
					},
				}
				p2 := pageLookupResult{
					num:  2,
					size: p1.size,
				}
				p3 := pageLookupResult{
					num: 3,
					size: geometry.Size{
						Width:  21 * geometry.Cm,
						Height: 14.8 * geometry.Cm,
					},
				}

				return []pageLookup{
					{r: pagerange.All, want: []pageLookupResult{p1, p2, p3}},
					{r: pagerange.MustSingle(1), want: []pageLookupResult{p1}},
					{r: pagerange.MustSingle(2), want: []pageLookupResult{p2}},
					{r: pagerange.MustNew(2, 3), want: []pageLookupResult{p2, p3}},
					{r: pagerange.MustNew(2, 15), want: []pageLookupResult{p2, p3}},
					{r: pagerange.MustNew(2, pagerange.Last), want: []pageLookupResult{p2, p3}},
					{r: pagerange.MustSingle(3), want: []pageLookupResult{p3}},
					{r: pagerange.MustNew(3, 100), want: []pageLookupResult{p3}},
					{r: pagerange.MustSingle(4)},
					{r: pagerange.MustNew(100, 102)},
				}
			}(),
		},
		{
			name: "validation fails",
			p: &parsertest.SimpleParser{
				ValidateErr: errTest,
			},
			wantValidateErr: errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)

			d := NewDocument(
				testutil.MustWriteFileString(t, filepath.Join(t.TempDir(), "doc"), "content"),
				WithStaticDocumentParser(tc.p))

			err := d.Validate(ctx)

			if diff := cmp.Diff(tc.wantValidateErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			for _, l := range tc.wantPages {
				t.Run(l.r.String(), func(t *testing.T) {
					got, err := d.ParsePages(ctx, l.r)

					if diff := cmp.Diff(l.wantErr, err, cmpopts.EquateErrors()); diff != "" {
						t.Fatalf("Lookup error diff (-want +got):\n%s", diff)
					}

					if gotCount, wantCount := len(got), len(l.want); gotCount != wantCount {
						t.Fatalf("Pages() returned %d pages, want %d", gotCount, wantCount)
					}

					for idx, i := range l.want {
						p := got[idx]

						if diff := cmp.Diff(p.Number(), i.num); diff != "" {
							t.Errorf("%#v number diff (-want +got):\n%s", p, diff)
						}

						if diff := cmp.Diff(p.Size(), i.size, geometry.EquateLength()); diff != "" {
							t.Errorf("%#v size diff (-want +got):\n%s", p, diff)
						}
					}
				})
			}
		})
	}
}

func TestDocumentParsePagesCache(t *testing.T) {
	errTest := errors.New("test error")

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	const knownPageCount = 3

	parser := &parsertest.SimpleParser{
		Pages: mustReadPages(t, "multipage.xml"),
	}

	d := NewDocument(os.DevNull, WithStaticDocumentParser(parser))

	if err := d.Validate(ctx); err != nil {
		t.Errorf("Validate() failed: %v", err)
	}

	// Populate cache
	if pages, err := d.ParsePages(ctx, pagerange.All); err != nil {
		t.Errorf("ParsePages() failed: %v", err)
	} else if got, want := len(pages), knownPageCount; got != want {
		t.Errorf("ParsePages() returned %d pages, want %d", got, want)
	}

	parser.ParseErr = errTest

	// Pages are in cache
	for pr, wantNumbers := range map[pagerange.Range][]int{
		pagerange.MustSingle(1):              {1},
		pagerange.MustSingle(knownPageCount): {3},
		pagerange.MustNew(2, knownPageCount): {2, 3},
		pagerange.MustNew(1, knownPageCount): {1, 2, 3},
	} {
		pages, err := d.ParsePages(ctx, pr)
		if err != nil {
			t.Errorf("ParsePages() failed: %v", err)
		}

		var numbers []int

		for _, i := range pages {
			numbers = append(numbers, i.Number())
		}

		if diff := cmp.Diff(wantNumbers, numbers, cmpopts.EquateEmpty()); diff != "" {
			t.Errorf("Page number diff (-want +got):\n%s", diff)
		}
	}

	for _, pr := range []pagerange.Range{
		// Pages not cached previously
		pagerange.MustSingle(knownPageCount + 1),
		pagerange.MustSingle(100),

		// Last page can't be cached as the page count is unknown
		pagerange.MustSingle(pagerange.Last),
	} {
		if _, err := d.ParsePages(ctx, pr); !errors.Is(err, errTest) {
			t.Errorf("ParsePages(%v) = %v, want %v", pr, err, errTest)
		}
	}
}
