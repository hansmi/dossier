package dossier

import (
	"context"
	"errors"
	"math"

	rtree "github.com/dhconnelly/rtreego"
	"github.com/hansmi/dossier/pkg/content"
	"github.com/hansmi/dossier/pkg/geometry"
	"github.com/hansmi/dossier/pkg/renderformat"
)

func toRTreeRect(r geometry.Rect) (rtree.Rect, error) {
	return rtree.NewRectFromPoints(
		rtree.Point{r.Top.Pt(), r.Left.Pt()},
		rtree.Point{r.Bottom.Pt(), r.Right.Pt()},
	)
}

type spatialAdapter struct {
	elem   content.Element
	bounds rtree.Rect
}

var _ rtree.Spatial = (*spatialAdapter)(nil)

func newAdapter(e content.Element) (rtree.Spatial, error) {
	bounds, err := toRTreeRect(e.Bounds())
	if err != nil {
		return nil, err
	}

	return &spatialAdapter{
		elem:   e,
		bounds: bounds,
	}, nil
}

func (a *spatialAdapter) Bounds() rtree.Rect {
	return a.bounds
}

var ErrStopVisitation = errors.New("stop visitation")

type PageElementVisitorFunc func(content.Element) error

// AsPageElementVisitor returns a visitor function wrapper filtering for
// elements of type T.
func AsPageElementVisitor[T content.Element](fn func(T) error) PageElementVisitorFunc {
	return func(elem content.Element) error {
		if v, ok := elem.(T); ok {
			return fn(v)
		}

		return nil
	}
}

type Page struct {
	doc   *Document
	num   int
	size  geometry.Size
	elems []content.Element
	tree  *rtree.Rtree
}

func newPage(doc *Document, p content.Page) (*Page, error) {
	objects := make([]rtree.Spatial, 0, len(p.Elements()))

	for _, e := range p.Elements() {
		obj, err := newAdapter(e)
		if err != nil {
			return nil, err
		}

		objects = append(objects, obj)
	}

	return &Page{
		doc:   doc,
		num:   p.Number(),
		size:  p.Size(),
		elems: p.Elements(),
		tree:  rtree.NewTree(2, 2, 8, objects...),
	}, nil
}

// Document returns the source document for the page.
func (p *Page) Document() *Document {
	return p.doc
}

// 1-based page number.
func (p *Page) Number() int {
	return p.num
}

// Physical page size.
func (p *Page) Size() geometry.Size {
	return p.size
}

func (p *Page) visitElements(bounds rtree.Rect, visitor PageElementVisitorFunc) error {
	var err error

	p.tree.SearchIntersect(bounds, func(_ []rtree.Spatial, obj rtree.Spatial) (refuse, abort bool) {
		// The filter function may still be called even after it requested the
		// search to be aborted. This is an apparent bug in the rtreego
		// upstream code. The condition on err avoids invoking the handler in
		// such cases.
		if err == nil {
			err = visitor(obj.(*spatialAdapter).elem)
		}

		return true, (err != nil)
	})

	if errors.Is(err, ErrStopVisitation) {
		err = nil
	}

	return err
}

// VisitElements invokes the visitor function for all elements. The visitation
// continues until either all elements have been visited or the visitor
// function returns a non-nil error. [ErrStopVisitation] stops the search
// immediately without failing the overall search. The visitation order is
// undefined.
func (p *Page) VisitElements(visitor PageElementVisitorFunc) error {
	rbounds, err := rtree.NewRectFromPoints(
		rtree.Point{-math.MaxFloat64, -math.MaxFloat64},
		rtree.Point{math.MaxFloat64, math.MaxFloat64},
	)
	if err != nil {
		return err
	}

	return p.visitElements(rbounds, visitor)
}

// VisitElementsIntersecting is like [VisitElements] with the additional
// restriction that only elements within the specified bounds are visited.
func (p *Page) VisitElementsIntersecting(bounds geometry.Rect, visitor PageElementVisitorFunc) error {
	rbounds, err := toRTreeRect(bounds)
	if err != nil {
		return err
	}

	return p.visitElements(rbounds, visitor)
}

func (p *Page) RenderUsing(ctx context.Context, r renderformat.Renderer) error {
	return p.doc.RenderPageUsing(ctx, p.num, r)
}
