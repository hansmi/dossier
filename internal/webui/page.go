package webui

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hansmi/dossier/internal/httperr"
	"github.com/hansmi/dossier/internal/webui/template"
	"github.com/hansmi/dossier/pkg/pagerange"
)

func (s *server) handlePage(w http.ResponseWriter, r *http.Request) error {
	pageNumber, err := strconv.Atoi(chi.URLParam(r, "num"))
	if err != nil {
		return httperr.New(http.StatusBadRequest, err)
	}

	doc, _, err := s.openDocument(r.Context())
	if err != nil {
		return err
	}

	fp, err := doc.Fingerprint()
	if err != nil {
		return err
	}

	pr, err := pagerange.Single(pageNumber)

	if err != nil {
		return httperr.New(http.StatusNotFound, fmt.Errorf("page %d not found: %w", pageNumber, err))
	}

	pages, err := doc.ParsePages(r.Context(), pr)
	if err != nil {
		return err
	}

	if !(len(pages) > 0 && pages[0].Number() == pageNumber) {
		return httperr.New(http.StatusNotFound, fmt.Errorf("page %d not found", pageNumber))
	}

	var messages []string

	page := pages[0]
	data := template.PageData{
		DocFingerprint: fp,
		Page:           page,
	}

	if cfg, err := s.compileSketch(); err != nil {
		messages = append(messages, fmt.Sprintf("Sketch: %v", err))
	} else if report, err := cfg.AnalyzePage(page); err != nil {
		messages = append(messages, fmt.Sprintf("Processing document: %v", err))
	} else {
		var nodes []template.SketchNodeData

		for idx, i := range report.Nodes() {
			nodes = append(nodes, template.SketchNodeData{
				ID:   fmt.Sprintf("sketch_node_%d_info", idx),
				Node: i,
			})
		}

		data.SketchNodes = nodes
	}

	return template.Base(template.BaseData{
		HeadTitle: fmt.Sprintf("Page %d of %s", page.Number(), filepath.Base(doc.Path())),
		Scripts: []string{
			"/static/page.js",
		},
		Messages: messages,
		Content:  template.PageContent(data),
		Sidebar:  template.PageSidebar(data),
	}).Render(r.Context(), w)
}
