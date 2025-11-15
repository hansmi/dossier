package template

import (
	"cmp"
	"context"
	"fmt"
	"io"
	"math"
	"slices"
	"strings"
	"unicode"

	"github.com/a-h/templ"
	"github.com/hansmi/dossier"
	"github.com/hansmi/dossier/pkg/content"
	"github.com/hansmi/dossier/pkg/geometry"
	"github.com/hansmi/dossier/pkg/sketch"
)

type SketchNodeData struct {
	*sketch.Node

	ID string
}

type PageData struct {
	DocFingerprint string
	Page           *dossier.Page
	SketchNodes    []SketchNodeData
}

func (d *PageData) size() geometry.Size {
	return d.Page.Size()
}

func (d PageData) widthInCssPixels() int {
	return int(math.Ceil(96 * d.size().Width.Inch()))
}

func (d PageData) imageData() PageImageData {
	return PageImageData{
		DocFingerprint: d.DocFingerprint,
		Page:           d.Page,
		Width:          d.widthInCssPixels(),
		ClassNames:     []string{"user-select-none"},
	}
}

func (d PageData) overlays() []pageViewerOverlayData {
	var result []pageViewerOverlayData

	var counter int

	nextID := func() string {
		counter++
		return fmt.Sprintf("overlay%x", counter)
	}

	// Errors are ignored
	d.Page.VisitElements(func(elem content.Element) error {
		o := newPageViewerOverlayData(nextID())
		o.PageSize = d.size()
		o.Bounds = elem.Bounds()
		o.ModalTarget = "#dossier_page_node_dialog_template"

		o.DataAttr["node-bounds"] = toJSON(map[string]float64{
			"top-pt":    o.Bounds.Top.Pt(),
			"left-pt":   o.Bounds.Left.Pt(),
			"right-pt":  o.Bounds.Right.Pt(),
			"bottom-pt": o.Bounds.Bottom.Pt(),
		})

		var nodeKind string
		var className string
		var hasContent bool

		switch elem.(type) {
		case content.Block:
			nodeKind = "Block"
			className = "dossier_doc_block"
			hasContent = hasContent || len(elem.(content.Block).Lines()) > 0
		case content.Line:
			nodeKind = "Line"
			className = "dossier_doc_line"
		}

		if nodeKind != "" {
			o.DataAttr["node-kind"] = nodeKind
		}

		if className != "" {
			o.Classes = append(o.Classes, className)
		}

		if telem, ok := elem.(content.TextElement); ok {
			text := telem.Text()
			o.DataAttr["node-text"] = toJSON(text)
			hasContent = hasContent || strings.ContainsFunc(text, func(r rune) bool {
				return !unicode.IsSpace(r)
			})
		}

		if !hasContent {
			o.Classes = append(o.Classes, "dossier_doc_empty_element")
		}

		result = append(result, o)

		return nil
	})

	for _, node := range d.SketchNodes {
		if !node.Valid() {
			continue
		}

		o := newPageViewerOverlayData(nextID())
		o.Title = node.Name()
		o.PageSize = d.size()
		o.Bounds = node.Bounds()
		o.Classes = append(o.Classes, "dossier_sketch_node")
		o.Order = 100

		if node.ID != "" {
			o.DataAttr["info-id"] = node.ID
		}

		result = append(result, o)
	}

	for _, node := range d.SketchNodes {
		for _, area := range node.SearchAreas() {
			o := newPageViewerOverlayData(nextID())
			o.PageSize = d.size()
			o.Bounds = area
			o.Classes = append(o.Classes, "dossier_sketch_node_search_area")
			o.Order = 50

			if node.ID != "" {
				o.DataAttr["info-id"] = node.ID
			}

			result = append(result, o)
		}
	}

	var compareRect = geometry.MakeRectRowColumnCompare(geometry.TopToBottom, geometry.LeftToRight)

	// Give overlays a predictable order.
	slices.SortStableFunc(result, func(a, b pageViewerOverlayData) int {
		return cmp.Or(
			cmp.Compare(a.Order, b.Order),
			compareRect(a.Bounds, b.Bounds),
		)
	})

	return result
}

type pageViewerOverlayData struct {
	ID          string
	Title       string
	PageSize    geometry.Size
	Bounds      geometry.Rect
	Classes     []string
	DataAttr    map[string]string
	Order       int
	ModalTarget string
}

func newPageViewerOverlayData(id string) pageViewerOverlayData {
	return pageViewerOverlayData{
		ID:       id,
		DataAttr: map[string]string{},
	}
}

func (d pageViewerOverlayData) style() string {
	top := 100 * float32(d.Bounds.Top/d.PageSize.Height)
	left := 100 * float32(d.Bounds.Left/d.PageSize.Width)
	width := 100 * float32(d.Bounds.Width()/d.PageSize.Width)
	height := 100 * float32(d.Bounds.Height()/d.PageSize.Height)

	return fmt.Sprintf(`top: %.1f%%; left: %.1f%%; width: %.1f%%; height: %.1f%%;`, top, left, width, height)
}

func pageViewerOverlay(data pageViewerOverlayData) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		buf := templ.GetBuffer()
		defer templ.ReleaseBuffer(buf)

		classes := []any{
			"dossier_viewer_overlay",
			"position-absolute",
			"shadow-sm",
			"d-flex",
			"align-items-stretch",
		}

		for _, i := range data.Classes {
			classes = append(classes, i)
		}

		if err := templ.RenderCSSItems(ctx, buf, classes...); err != nil {
			return err
		}

		attr := templ.Attributes{
			"id":    data.ID,
			"class": templ.CSSClasses(classes).String(),
			"style": data.style(),
		}

		for k, v := range data.DataAttr {
			attr["data-"+k] = v
		}

		if data.Title != "" {
			attr["data-bs-toggle"] = "tooltip"
			attr["data-bs-container"] = toJSON("#" + data.ID)
			attr["data-bs-title"] = toJSON(data.Title)
		}

		buf.WriteString(`<div`)
		if err := templ.RenderAttributes(ctx, buf, attr); err != nil {
			return err
		}
		buf.WriteString(`>`)

		buf.WriteString(`<button type="button" class="flex-fill btn btn-sm btn-outline-primary overflow-hidden"`)
		if data.ModalTarget != "" {
			buf.WriteString(` data-bs-toggle="modal" data-bs-target="`)
			buf.WriteString(templ.EscapeString(data.ModalTarget))
			buf.WriteString(`"`)
		}
		buf.WriteString(`></button>`)

		buf.WriteString(`</div>`)

		_, err := buf.WriteTo(w)

		return err
	})
}
