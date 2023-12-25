package webui

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/hansmi/dossier/internal/webui/template"
	"github.com/hansmi/dossier/pkg/pagerange"
)

func (s *server) handleOverview(w http.ResponseWriter, r *http.Request) error {
	doc, fi, err := s.openDocument(r.Context())
	if err != nil {
		return err
	}

	var data template.OverviewContentData

	if fp, err := doc.Fingerprint(); err != nil {
		return err
	} else {
		data.DocFingerprint = fp
	}

	pr := pagerange.All

	if s.opts.maxPages > 0 {
		if pr, err = pagerange.New(1, s.opts.maxPages); err != nil {
			return err
		}
	}

	if pages, err := doc.ParsePages(r.Context(), pr); err != nil {
		return err
	} else {
		data.Pages = pages
	}

	mtime := fi.ModTime()

	return template.Base(template.BaseData{
		HeadTitle:    filepath.Base(doc.Path()),
		TopNavActive: template.TopNavOverview,
		Content:      template.OverviewContent(data),
		Sidebar: template.OverviewSidebar(template.OverviewSidebarData{
			Path:        doc.Path(),
			Size:        humanize.IBytes(uint64(fi.Size())),
			ModTime:     mtime.Format(time.RFC822),
			ModTimeFull: mtime.Format(time.RFC1123),
		}),
	}).Render(r.Context(), w)
}
