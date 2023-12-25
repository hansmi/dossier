package stext

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"github.com/hansmi/dossier/pkg/geometry"
	"go.uber.org/multierr"
	"golang.org/x/net/html/charset"
)

type Document struct {
	Name  string `xml:"name,attr"`
	Pages []Page `xml:"page"`
}

// DocumentFromXML unmarshals a document from the XML structure used for
// mutool's "stext" output format.
func DocumentFromXML(r io.Reader) (*Document, error) {
	doc := &Document{}

	if err := doc.Unmarshal(r); err != nil {
		return nil, fmt.Errorf("parsing XML: %w", err)
	}

	return doc, nil
}

// DocumentFromXMLFile is the same as [DocumentFromXML], but reads the contents
// from a file.
func DocumentFromXMLFile(path string) (_ *Document, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer multierr.AppendInvoke(&err, multierr.Close(f))

	return DocumentFromXML(f)
}

func (d *Document) Unmarshal(r io.Reader) error {
	dec := xml.NewDecoder(r)
	dec.Strict = true
	dec.CharsetReader = charset.NewReaderLabel

	return dec.Decode(&struct {
		*Document
		XMLName xml.Name `xml:"document"`
	}{Document: d})
}

type Page struct {
	ID     string          `xml:"id,attr"`
	Width  geometry.Length `xml:"width,attr"`
	Height geometry.Length `xml:"height,attr"`
	Blocks []Block         `xml:"block"`
}

type Block struct {
	BBox  geometry.Rect
	Lines []Line `xml:"line"`
}

var _ xml.Unmarshaler = (*Block)(nil)

func (b *Block) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type plain Block
	var elem struct {
		*plain
		BBox rectAttr `xml:"bbox,attr"`
	}
	elem.plain = (*plain)(b)

	if err := d.DecodeElement(&elem, &start); err != nil {
		return err
	}

	b.BBox = geometry.Rect(elem.BBox)

	return nil
}

type Line struct {
	BBox      geometry.Rect
	FontSpans []FontSpan `xml:"font"`
}

var _ xml.Unmarshaler = (*Line)(nil)

func (l *Line) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type plain Line
	var elem struct {
		*plain
		BBox rectAttr `xml:"bbox,attr"`
	}
	elem.plain = (*plain)(l)

	if err := d.DecodeElement(&elem, &start); err != nil {
		return err
	}

	l.BBox = geometry.Rect(elem.BBox)

	return nil
}

type FontSpan struct {
	FontName string          `xml:"name,attr"`
	FontSize geometry.Length `xml:"size,attr"`
	Chars    []Char          `xml:"char"`
}

type Char struct {
	C      rune
	Bounds geometry.Rect
}

func (c *Char) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type plain Char
	var elem struct {
		*plain
		C    runeAttr `xml:"c,attr"`
		Quad quadAttr `xml:"quad,attr"`
	}
	elem.plain = (*plain)(c)

	if err := d.DecodeElement(&elem, &start); err != nil {
		return err
	}

	c.C = rune(elem.C)
	c.Bounds = geometry.Rect(elem.Quad)

	return nil
}
