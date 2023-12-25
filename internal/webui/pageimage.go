package webui

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hansmi/dossier/internal/httperr"
	"github.com/hansmi/dossier/pkg/renderformat"
)

func (s *server) handlePageImage(w http.ResponseWriter, r *http.Request) error {
	pageNumber, err := strconv.Atoi(chi.URLParam(r, "num"))
	if err != nil {
		return httperr.New(http.StatusBadRequest, err)
	}

	if err := r.ParseForm(); err != nil {
		return httperr.New(http.StatusBadRequest, err)
	}

	f := &renderformat.PNG{
		Output: w,
		Width:  100,
	}

	if widthStr := r.FormValue("width"); widthStr != "" {
		if f.Width, err = strconv.Atoi(widthStr); err != nil {
			return httperr.New(http.StatusBadRequest, err)
		}
	}

	doc, fi, err := s.openDocument(r.Context())
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Last-Modified", fi.ModTime().Format(http.TimeFormat))

	if fp := r.FormValue("docfp"); fp != "" {
		w.Header().Set("Cache-Control", "private, max-age=86400")
	}

	return doc.RenderPageUsing(r.Context(), pageNumber, f)
}
