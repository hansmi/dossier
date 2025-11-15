package webui

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hansmi/dossier"
	"github.com/hansmi/dossier/internal/httperr"
	"github.com/hansmi/dossier/pkg/sketch"
)

//go:embed static/*.css static/*.js
var staticFiles embed.FS

type serverOptions struct {
	maxConcurrent int
	maxPages      int
	sketchPath    string
	documentPath  string
}

type server struct {
	opts serverOptions
}

func newServer(opts serverOptions) (*server, error) {
	if opts.maxConcurrent < 1 {
		opts.maxConcurrent = 1
	}

	s := &server{
		opts: opts,
	}

	return s, nil
}

func (s *server) makeRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Throttle(s.opts.maxConcurrent))

	r.Get(`/`, httperr.WrapHandler(httperr.HandlerFunc(s.handleOverview)))
	r.Get(`/page/{num:\d+}`, httperr.WrapHandler(httperr.HandlerFunc(s.handlePage)))
	r.Get(`/page/{num:\d+}/image`, httperr.WrapHandler(httperr.HandlerFunc(s.handlePageImage)))

	r.Handle("/static/*", http.FileServerFS(staticFiles))

	return r
}

func (s *server) compileSketch() (*sketch.Sketch, error) {
	content, err := os.ReadFile(s.opts.sketchPath)
	if err != nil {
		return nil, fmt.Errorf("reading sketch file: %w", err)
	}

	cfg, err := sketch.CompileFromTextproto(content)
	if err != nil {
		return nil, fmt.Errorf("compiling sketch: %w", err)
	}

	return cfg, nil
}

func (s *server) openDocument(ctx context.Context) (*dossier.Document, fs.FileInfo, error) {
	fi, err := os.Lstat(s.opts.documentPath)
	if err != nil {
		return nil, nil, err
	}

	return dossier.NewDocument(s.opts.documentPath), fi, nil
}
