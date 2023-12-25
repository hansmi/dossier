package mutool

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/dossier/internal/mutool/stext"
	"github.com/hansmi/dossier/internal/testfiles"
	"github.com/hansmi/dossier/internal/testutil"
	"github.com/hansmi/dossier/pkg/geometry"
	"github.com/hansmi/dossier/pkg/pagerange"
	"github.com/hansmi/dossier/pkg/renderformat"
)

var errTest = errors.New("test error")

type fakeMutool struct {
	show func(showArgs) error
	draw func(drawArgs) error
}

func (*fakeMutool) CheckCommand(_ context.Context) error {
	return nil
}

func (m *fakeMutool) Show(_ context.Context, a showArgs) error {
	return m.show(a)
}

func (m *fakeMutool) Draw(_ context.Context, a drawArgs) error {
	return m.draw(a)
}

type fakeXmllint struct {
	recover func(recoverArgs) error
}

func (*fakeXmllint) CheckCommand(_ context.Context) error {
	return nil
}

func (x *fakeXmllint) Recover(_ context.Context, a recoverArgs) error {
	return x.recover(a)
}

func TestWrapperValidate(t *testing.T) {
	emptyFile := filepath.Join(t.TempDir(), "empty")

	for _, tc := range []struct {
		name    string
		show    func(*testing.T, showArgs) error
		wantErr error
	}{
		{
			name: "success",
			show: func(t *testing.T, a showArgs) error {
				want := showArgs{
					input: emptyFile,
				}

				if diff := cmp.Diff(want, a, cmp.AllowUnexported(showArgs{})); diff != "" {
					t.Errorf("Args diff (-want +got):\n%s", diff)
				}

				return nil
			},
		},
		{
			name:    "error",
			show:    func(*testing.T, showArgs) error { return errTest },
			wantErr: errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)

			w := Wrapper{
				mutool: &fakeMutool{
					show: func(a showArgs) error {
						return tc.show(t, a)
					},
				},
			}

			err := w.Validate(ctx, emptyFile)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestWrapperStructuredText(t *testing.T) {
	emptyFile := filepath.Join(t.TempDir(), "empty")

	for _, tc := range []struct {
		name          string
		drawErr       error
		drawOutput    string
		wantRecover   bool
		recoverErr    error
		recoverOutput string
		wantErr       error
		want          *stext.Document
		wantOpts      cmp.Options
	}{
		{
			name:    "empty output",
			wantErr: io.EOF,
		},
		{
			name: "trivial document",
			drawOutput: `
				<?xml version="1.0"?>
				<document name="test.pdf">
				</document>
			`,
			want: &stext.Document{
				Name: "test.pdf",
			},
		},
		{
			name:    "draw error",
			drawErr: errTest,
			wantErr: errTest,
		},
		{
			name:          "bad xml",
			drawOutput:    `>document<`,
			wantRecover:   true,
			recoverOutput: `>more bad xml<`,
			wantErr:       cmpopts.AnyError,
		},
		{
			name:          "recovery from invalid entity",
			drawOutput:    `<document name="&#xffff;"></document>`,
			wantRecover:   true,
			recoverOutput: `<document />`,
			want:          &stext.Document{},
		},
		{
			name:       "real",
			drawOutput: testutil.MustReadFileString(t, testfiles.All, "corners.xml"),
			wantOpts: cmp.Options{
				cmpopts.IgnoreFields(stext.Block{}, "Lines"),
			},
			want: &stext.Document{
				Name: "corners.pdf",
				Pages: []stext.Page{
					{
						ID:     "page1",
						Width:  geometry.Pt * 175.748,
						Height: geometry.Pt * 249.448,
						Blocks: []stext.Block{
							{BBox: geometry.RectFromPoints(28.375, 26.354, 40.045, 38)},
							{BBox: geometry.RectFromPoints(27.365, 211.819, 39.795, 223.461)},
							{BBox: geometry.RectFromPoints(134.629, 26.354, 147.676, 38)},
							{BBox: geometry.RectFromPoints(133.887, 211.819, 147.694, 223.461)},
						},
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)

			recoverCalled := false

			w := Wrapper{
				mutool: &fakeMutool{
					draw: func(a drawArgs) error {
						testutil.MustWriteFileString(t, a.output, tc.drawOutput)

						return tc.drawErr
					},
				},

				xmllint: &fakeXmllint{
					recover: func(a recoverArgs) error {
						recoverCalled = true

						if !tc.wantRecover {
							t.Errorf("Xmllint recover called unexpectedly: %#v", a)
						}

						testutil.MustWriteFileString(t, a.output, tc.recoverOutput)

						return tc.recoverErr
					},
				},
			}

			got, err := w.StructuredText(ctx, emptyFile, pagerange.All)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				opts := append(cmp.Options{
					cmpopts.EquateEmpty(),
					geometry.EquateLength(),
				}, tc.wantOpts...)

				if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
					t.Errorf("Document diff (-want +got):\n%s", diff)
				}
			}

			if diff := cmp.Diff(tc.wantRecover, recoverCalled); diff != "" {
				t.Errorf("Recover invocation diff (-want +got):\n%s", diff)
			}
		})
	}
}

type unsupportedRenderer struct{}

func (unsupportedRenderer) String() string {
	return "unsupported"
}

func TestWrapperDraw(t *testing.T) {
	emptyFile := filepath.Join(t.TempDir(), "empty")

	for _, tc := range []struct {
		name       string
		pageNum    int
		renderer   renderformat.Renderer
		draw       func(*testing.T, drawArgs) error
		wantErr    error
		wantOutput string
	}{
		{
			name: "success",
			draw: func(t *testing.T, a drawArgs) error {
				if _, err := io.WriteString(a.stdout, "test output"); err != nil {
					t.Errorf("WriteString() failed: %v", err)
				}

				return nil
			},
			wantOutput: "test output",
		},
		{
			name: "error",
			draw: func(*testing.T, drawArgs) error {
				return errTest
			},
			wantErr: errTest,
		},
		{
			name:     "bad format",
			renderer: &unsupportedRenderer{},
			wantErr:  os.ErrInvalid,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)

			w := Wrapper{
				mutool: &fakeMutool{
					draw: func(a drawArgs) error {
						return tc.draw(t, a)
					},
				},
			}

			var out bytes.Buffer

			if tc.renderer == nil {
				tc.renderer = &renderformat.PNG{
					Output: &out,
				}
			}

			err := w.Draw(ctx, emptyFile, tc.pageNum, tc.renderer)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.wantOutput, out.String()); diff != "" {
					t.Errorf("Output diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
