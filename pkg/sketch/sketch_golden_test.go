package sketch

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hansmi/aurum"
	"github.com/hansmi/dossier"
	"github.com/hansmi/dossier/internal/muparser"
	"github.com/hansmi/dossier/internal/testfiles"
	"github.com/hansmi/dossier/internal/testutil"
	"github.com/hansmi/dossier/pkg/geometry"
	"github.com/hansmi/dossier/pkg/pagerange"
	"github.com/hansmi/dossier/pkg/parsertest"
	"github.com/hansmi/dossier/proto/sketchpb"
)

func init() {
	aurum.Init()
}

func readTestDocument(t *testing.T, name string) *dossier.Document {
	t.Helper()

	f, err := testfiles.All.Open(name)
	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	pages, err := muparser.ReadPagesFromXML(f)
	if err != nil {
		t.Fatalf("ReadPagesFromXML() failed: %v", err)
	}

	return dossier.NewDocument(
		testutil.MustWriteFile(t, filepath.Join(t.TempDir(), "empty"), nil),
		dossier.WithStaticDocumentParser(&parsertest.SimpleParser{
			Pages: pages,
		}))
}

func TestSketchGolden(t *testing.T) {
	g := aurum.Golden{
		Dir:   "./testdata/layout-golden",
		Codec: &aurum.TextProtoCodec{},
	}

	unit := geometry.RoundedLength{
		Unit:    geometry.Pt,
		Nearest: geometry.Pt,
	}

	for _, tc := range []struct {
		name     string
		sketch   *sketchpb.Sketch
		document string
	}{
		{
			name:     "empty",
			document: "multipage.xml",
			sketch:   &sketchpb.Sketch{},
		},
		{
			name:     "multimatch",
			document: "multipage.xml",
			sketch: testutil.MustUnmarshalTextproto(t, `
tags: "sketch tag"
nodes: {
  name: "text"
  search_areas {
    top_left {
      abs { left: { } top: {} }
    }
	bottom_right {
	  abs { left: { cm: 20 } top: { cm: 30 } }
	}
  }
  line_text: { regex: "(?i)\\b(?:Lorem|Second|World)\\b" }
}
`, &sketchpb.Sketch{}),
		},
		{
			name:     "corners",
			document: "corners.xml",
			sketch: testutil.MustUnmarshalTextproto(t, `
nodes: {
  name: "tl"
  search_areas {
    top_left {
      abs: { left: { mm: 1 } top: { mm: 1 } }
    }
    width: { cm: 3 }
    height: { cm: 2 }
  }
  block_text: {
    regex: "(?i)^\\s*tl\\b"
  }
  tags: "top left"
}

nodes: {
  name: "tr"
  search_areas {
    top_left {
      rel: {
        node: "tl"
        feature: TOP_RIGHT
        offset: {
          width: { cm: 1 }
          height: { cm: 0 }
        }
      }
    }
    width: { cm: 3 }
    height: { cm: 4 }
  }

  line_text: {
    regex: "(?i)\\btr\\b"
  }

  tags: "top right"
}
`, &sketchpb.Sketch{}),
		},
		{
			name:     "acme-invoice",
			document: "acme-invoice-11321-19.xml",
			sketch: testutil.MustUnmarshalTextproto(t,
				testutil.MustReadFileString(t, os.DirFS("."), "testdata/acme-invoice.textproto"),
				&sketchpb.Sketch{}),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			doc := readTestDocument(t, tc.document)

			sketch, err := Compile(tc.sketch)
			if err != nil {
				t.Fatalf("Compile() failed: %v", err)
			}

			got, err := sketch.AnalyzeDocument(context.Background(), doc, pagerange.All)
			if err != nil {
				t.Fatalf("FullReport() failed: %v", err)
			}

			g.Assert(t, tc.name, got.AsProto(unit))
		})
	}
}
