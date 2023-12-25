package dossier

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/gabriel-vasile/mimetype"
	"github.com/hansmi/dossier/pkg/pagerange"
	"github.com/hansmi/dossier/pkg/renderformat"
	lru "github.com/hashicorp/golang-lru/v2"
	"golang.org/x/sys/unix"
)

var ErrUnsupportedFileFormat = errors.New("unsupported file format")

const pageCacheSize = 10

type DocumentOption func(*Document)

type DocumentParserFactory func(path, contentType string) (Parser, error)

// Configure a custom factory function to create parser instances. When left
// unconfigured an appropriate parser is automatically chosen based on the
// source file's content.
func WithDocumentParserFactory(f DocumentParserFactory) DocumentOption {
	return func(doc *Document) {
		doc.parserFactory = f
	}
}

// Use a fixed parser for all documents without considering the content type.
func WithStaticDocumentParser(p Parser) DocumentOption {
	return WithDocumentParserFactory(func(_, _ string) (Parser, error) {
		return p, nil
	})
}

type Document struct {
	path          string
	parserFactory DocumentParserFactory

	mu          sync.Mutex
	contentType string
	parser      Parser
	pageCache   *lru.Cache[int, *Page]
}

// NewDocument constructs a new document. The file must not be modified while
// it's being used. Operations may open and close the file multiple times.
func NewDocument(path string, opts ...DocumentOption) *Document {
	doc := &Document{
		path: path,
	}

	if pageCache, err := lru.New[int, *Page](pageCacheSize); err != nil {
		panic(err)
	} else {
		doc.pageCache = pageCache
	}

	for _, opt := range opts {
		opt(doc)
	}

	if doc.parserFactory == nil {
		doc.parserFactory = MuPdfParserFactory{}.Create
	}

	return doc
}

// Path returns the file path given to [NewDocument].
func (d *Document) Path() string {
	return d.path
}

func (d *Document) getContentType() (string, error) {
	if d.contentType == "" {
		if ct, err := mimetype.DetectFile(d.path); err != nil {
			return "", err
		} else {
			d.contentType = ct.String()
		}
	}

	return d.contentType, nil
}

// ContentType determines and returns the MIME content-type of the source file.
// The returned string may contain parameters (e.g. charset). Use
// [mime.ParseMediaType] or similar to parse the type.
func (d *Document) ContentType() (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.getContentType()
}

func (d *Document) getParser() (Parser, error) {
	if d.parser == nil {
		contentType, err := d.getContentType()
		if err != nil {
			return nil, err
		}

		if parser, err := d.parserFactory(d.path, contentType); err != nil {
			return nil, err
		} else if parser == nil {
			return nil, fmt.Errorf("%w: %s", ErrUnsupportedFileFormat, contentType)
		} else {
			d.parser = parser
		}
	}

	return d.parser, nil
}

// Fingerprint returns a best-effort file version identifier in the form of an
// opaque, non-empty string. While a changed fingerprint is indicative of
// a modified file, the fingerprint may also change for an unchanged file.
func (d *Document) Fingerprint() (string, error) {
	// Parsers are very likely to follow symlinks (muPDF does).
	fi, err := os.Stat(d.path)
	if err != nil {
		return "", err
	}

	numbers := []uint64{
		uint64(fi.Size()),
		uint64(fi.Mode().Type()),
	}

	if st, ok := fi.Sys().(*unix.Stat_t); ok && st != nil {
		numbers = append(numbers, st.Dev, st.Ino)
	}

	var buf bytes.Buffer

	binary.Write(&buf, binary.LittleEndian, numbers)

	if mtime, err := fi.ModTime().UTC().GobEncode(); err != nil {
		return "", err
	} else {
		buf.Write(mtime)
	}

	buf.WriteByte('\x00')
	buf.WriteString(filepath.Clean(d.path))
	buf.WriteByte('\x00')
	buf.WriteString(fi.Name())

	digest := sha256.Sum256(buf.Bytes())

	return base64.RawURLEncoding.EncodeToString(digest[:]), nil
}

// Validate is a simple check whether the document can be read and parsed.
func (d *Document) Validate(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	p, err := d.getParser()
	if err != nil {
		return err
	}

	return p.Validate(ctx)
}

// ParsePages uses the underlying document parser to read and parse pages. The
// returned slice may contain fewer or more pages than requested by the given
// range, depending on what the document actually contains and the parser's
// behaviour. Page numbers can be determined via [Page.Number].
func (d *Document) ParsePages(ctx context.Context, r pagerange.Range) ([]*Page, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	parser, err := d.getParser()
	if err != nil {
		return nil, err
	}

	var result []*Page

	if r.Lower != pagerange.Last {
		// Best-effort cache lookup starting at the lower end of the requested
		// range.
		for checked, cacheSize := 0, d.pageCache.Len(); checked < cacheSize && r.Lower <= r.Upper; checked++ {
			page, ok := d.pageCache.Get(r.Lower)
			if !ok || page == nil {
				break
			}

			r.Lower++
			result = append(result, page)
		}
	}

	if r.Lower == pagerange.Last || r.Lower <= r.Upper {
		pages, err := parser.ParsePages(ctx, r)
		if err != nil {
			return nil, err
		}

		for _, parsed := range pages {
			p, err := newPage(d, parsed)
			if err != nil {
				return nil, fmt.Errorf("page %d: %w", parsed.Number(), err)
			}

			d.pageCache.Add(p.Number(), p)

			result = append(result, p)
		}
	}

	return result, nil
}

// RenderPageUsing writes a single page using the given renderer, e.g. as a PNG
// image via [renderformat.PNG].
func (d *Document) RenderPageUsing(ctx context.Context, num int, r renderformat.Renderer) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	parser, err := d.getParser()
	if err != nil {
		return err
	}

	return parser.RenderPage(ctx, num, r)
}
