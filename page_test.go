package dossier

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hansmi/dossier/pkg/content"
	"github.com/hansmi/dossier/pkg/geometry"
)

type fakePage struct {
	num   int
	size  geometry.Size
	elems []content.Element
}

func (p *fakePage) Number() int {
	return p.num
}

func (p *fakePage) Size() geometry.Size {
	return p.size
}

func (p *fakePage) Elements() []content.Element {
	return p.elems
}

type fakeLine struct {
	bounds geometry.Rect
	text   string
}

var _ content.Line = (*fakeLine)(nil)

func (l *fakeLine) Kind() content.Line {
	return nil
}

func (l *fakeLine) Bounds() geometry.Rect {
	return l.bounds
}

func (l *fakeLine) Text() string {
	return l.text
}

func (l *fakeLine) RangeBounds(start, end int) geometry.Rect {
	return geometry.Rect{}
}

type lineVisitor []string

func (v *lineVisitor) visit(e content.Element) error {
	if l, ok := e.(content.Line); ok {
		*v = append(*v, l.Text())
	}

	return nil
}

func TestVisitElements(t *testing.T) {
	for _, tc := range []struct {
		name      string
		page      content.Page
		wantTexts []string
		check     func(*testing.T, *Page)
	}{
		{
			name: "empty",
			page: &fakePage{},
		},
		{
			name: "lines",
			page: &fakePage{
				elems: []content.Element{
					&fakeLine{
						text:   "first",
						bounds: geometry.RectFromCentimeters(1, 2, 3, 4),
					},
					&fakeLine{
						text:   "inside second",
						bounds: geometry.RectFromCentimeters(8, 10.5, 12, 17),
					},
					&fakeLine{text: "zero"},
					&fakeLine{
						text:   "second",
						bounds: geometry.RectFromCentimeters(0, 10, 15, 20),
					},
				},
			},
			wantTexts: []string{"first", "inside second", "zero", "second"},
			check: func(t *testing.T, p *Page) {
				var got lineVisitor

				if err := p.VisitElementsIntersecting(geometry.RectFromCentimeters(3, 9, 20, 25), got.visit); err != nil {
					t.Errorf("VisitElements() failed: %v", err)
				}

				want := []string{"inside second", "second"}

				if diff := cmp.Diff(want, []string(got)); diff != "" {
					t.Errorf("Matches diff (-want +got):\n%s", diff)
				}
			},
		},
		{
			name:      "corners",
			page:      mustReadPages(t, "corners.xml")[0],
			wantTexts: []string{"TL", "BL", "TR", "BR"},
		},
		{
			name: "lorem-mixed",
			page: mustReadPages(t, "lorem-mixed.xml")[0],
			wantTexts: []string{
				"Lorem ipsum dolor sit amet, consectetur adipisici elit, sed eiusmod tempor incidunt ",
				"ut labore et dolore magna aliqua.",
				"At vero eos et accusam et justo duo dolores et ea",
				"rebum. Stet clita kasd gubergren, no sea ",
				"takimata sanctus est Lorem ipsum dolor sit ",
				"amet.",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p, err := newPage(nil, tc.page)
			if err != nil {
				t.Fatalf("newPage() failed: %v", err)
			}

			var got lineVisitor

			if err := p.VisitElements(got.visit); err != nil {
				t.Errorf("VisitElements() failed: %v", err)
			}

			if diff := cmp.Diff(tc.wantTexts, []string(got)); diff != "" {
				t.Errorf("Matches diff (-want +got):\n%s", diff)
			}

			if tc.check != nil {
				tc.check(t, p)
			}
		})
	}
}
