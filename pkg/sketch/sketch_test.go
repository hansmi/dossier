package sketch

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/dossier/pkg/dossiertest"
	"github.com/hansmi/dossier/pkg/pagerange"
)

func TestCompileFromTextprotoString(t *testing.T) {
	for _, tc := range []struct {
		name    string
		sketch  string
		wantErr error
	}{
		{
			name: "empty",
		},
		{
			name: "two nodes with dependency",
			sketch: `
nodes: {
  name: "label"
  search_areas {
    top_left {
      abs { left: { cm: 2 } top: { cm: 9 } }
    }
    width { cm: 8 }
    height { cm: 3 }
  }
  line_text: { regex: "(?mi)^\\s*Number\\b" }
}
nodes: {
  name: "value"
  search_areas {
    left { abs { cm: 6 } }
    top {
      rel {
        node: "label"
        feature: TOP_RIGHT
        offset: { cm: -0.5 }
      }
    }
    width { cm: 12 }
    height { cm: 19 }
  }
  line_text: { regex: "(?mi)^\\s*(\\d{8,12})\\s*$" }
}
`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CompileFromTextprotoString(tc.sketch)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				doc := dossiertest.NewEmptyDocument(t)

				if _, err := got.AnalyzeDocument(context.Background(), doc, pagerange.All); err != nil {
					t.Errorf("AnalyzeDocument() failed: %v", err)
				}
			}
		})
	}
}
