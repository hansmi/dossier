package muparser

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/dossier/internal/mutool/stext"
	"github.com/hansmi/dossier/pkg/content"
	"github.com/hansmi/dossier/pkg/geometry"
	"github.com/hansmi/dossier/pkg/pagerange"
	"github.com/hansmi/dossier/pkg/renderformat"
)

var errUnimplemented = errors.New("not implemented")
var errValidation = errors.New("validation error")
var errDraw = errors.New("draw failed")

type fakeTool struct {
	validation func() error
	stext      func() (*stext.Document, error)
	draw       func() error
}

func (t *fakeTool) Validate(context.Context, string) error {
	if t.validation == nil {
		return errUnimplemented
	}

	return t.validation()
}

func (t *fakeTool) StructuredText(context.Context, string, pagerange.Range) (*stext.Document, error) {
	if t.stext == nil {
		return nil, errUnimplemented
	}

	return t.stext()
}

func (t *fakeTool) Draw(context.Context, string, int, renderformat.Renderer) error {
	if t.draw == nil {
		return errUnimplemented
	}

	return t.draw()
}

func TestValidate(t *testing.T) {
	for _, tc := range []struct {
		name    string
		tool    ToolWrapper
		wantErr error
	}{
		{
			name: "success",
			tool: &fakeTool{
				validation: func() error {
					return nil
				},
			},
		},
		{
			name: "failure",
			tool: &fakeTool{
				validation: func() error {
					return errValidation
				},
			},
			wantErr: errValidation,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			p := New(filepath.Join(t.TempDir(), "unused"), tc.tool)

			err := p.Validate(ctx)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParsePages(t *testing.T) {
	for _, tc := range []struct {
		name    string
		tool    ToolWrapper
		want    []content.Page
		wantErr error
	}{
		{
			name: "success",
			tool: &fakeTool{
				stext: func() (*stext.Document, error) {
					return &stext.Document{
						Name: "unused",
						Pages: []stext.Page{{
							ID:     "page 1",
							Width:  10 * geometry.Cm,
							Height: 20 * geometry.Cm,
						}},
					}, nil
				},
			},
			want: []content.Page{&Page{
				num: 1,
				size: geometry.Size{
					Width:  10 * geometry.Cm,
					Height: 20 * geometry.Cm,
				},
			}},
		},
		{
			name: "parse fails",
			tool: &fakeTool{
				stext: func() (*stext.Document, error) {
					return nil, errDraw
				},
			},
			wantErr: errDraw,
		},
		{
			name: "multipage",
			tool: &fakeTool{
				stext: func() (*stext.Document, error) {
					return loadTestDocument(t, "multipage.xml"), nil
				},
			},
			want: []content.Page{
				&Page{
					num: 1,
					size: geometry.Size{
						Width:  14.8 * geometry.Cm,
						Height: 21 * geometry.Cm,
					},
				},
				&Page{
					num: 2,
					size: geometry.Size{
						Width:  14.8 * geometry.Cm,
						Height: 21 * geometry.Cm,
					},
				},
				&Page{
					num: 3,
					size: geometry.Size{
						Width:  21 * geometry.Cm,
						Height: 14.8 * geometry.Cm,
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			p := New(filepath.Join(t.TempDir(), "unused"), tc.tool)

			got, err := p.ParsePages(ctx, pagerange.All)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				opts := cmp.Options{
					cmp.AllowUnexported(Page{}),
					cmpopts.EquateEmpty(),
					geometry.EquateLength(),
					cmpopts.IgnoreFields(Page{}, "elements"),
				}

				if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
					t.Errorf("ParsePages() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestRenderPage(t *testing.T) {
	for _, tc := range []struct {
		name     string
		tool     ToolWrapper
		renderer renderformat.Renderer
		wantErr  error
	}{
		{
			name: "success",
			tool: &fakeTool{
				draw: func() error {
					return nil
				},
			},
		},
		{
			name: "draw fails",
			tool: &fakeTool{
				draw: func() error {
					return errDraw
				},
			},
			wantErr: errDraw,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			p := New(filepath.Join(t.TempDir(), "unused"), tc.tool)

			err := p.RenderPage(ctx, 0, tc.renderer)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}
		})
	}
}
