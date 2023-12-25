package geometry

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/dossier/internal/testutil"
	"github.com/hansmi/dossier/proto/geometrypb"
	"google.golang.org/protobuf/testing/protocmp"
)

func validateRectangles(t *testing.T, all ...Rect) {
	t.Helper()

	for _, r := range all {
		if err := r.Validate(); err != nil {
			t.Errorf("%v: %v", r, err)
		}
	}
}

func TestRectFrom(t *testing.T) {
	for _, tc := range []struct {
		name       string
		got        Rect
		want       Rect
		wantString string
		wantEmpty  bool
		wantWidth  Length
		wantHeight Length
	}{
		{
			name:       "RectFromPoints empty",
			got:        RectFromPoints(0, 0, 0, 0),
			wantString: "(0, 0)-(0, 0)",
			wantEmpty:  true,
		},
		{
			name: "RectFromPoints",
			got:  RectFromPoints(72, 0.5*72, 2*72, 10*72),
			want: Rect{
				Left:   Inch,
				Top:    1.27 * Cm,
				Right:  2 * Inch,
				Bottom: 10 * Inch,
			},
			wantString: "(2.54cm, 1.27cm)-(5.08cm, 25.4cm)",
			wantWidth:  Inch,
			wantHeight: 9.5 * Inch,
		},
		{
			name: "RectFromPoints negative origin",
			got:  RectFromPoints(-100.0*72, -50.0*72, 200*72, 400*72),
			want: Rect{
				Left:   -100 * Inch,
				Top:    -50 * Inch,
				Right:  200 * Inch,
				Bottom: 400 * Inch,
			},
			wantString: "(-254cm, -127cm)-(508cm, 1016cm)",
			wantWidth:  300 * Inch,
			wantHeight: 450 * Inch,
		},
		{
			name: "RectFromCentimeters",
			got:  RectFromCentimeters(1, 8, 3, 11),
			want: Rect{
				Left:   Cm.Mul(1),
				Top:    8 * Cm,
				Right:  Cm.Mul(3),
				Bottom: 11 * Cm,
			},
			wantString: "(1cm, 8cm)-(3cm, 11cm)",
			wantWidth:  2 * Cm,
			wantHeight: 3 * Cm,
		},
		{
			name: "RectFromXYWH",
			got:  RectFromXYWH(1*Mm, 2*Mm, 3*Cm, 4*Cm),
			want: Rect{
				Left:   1 * Mm,
				Top:    2 * Mm,
				Right:  31 * Mm,
				Bottom: 42 * Mm,
			},
			wantString: "(0.1cm, 0.2cm)-(3.1cm, 4.2cm)",
			wantWidth:  3 * Cm,
			wantHeight: 4 * Cm,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			validateRectangles(t, tc.got, tc.want)

			if diff := cmp.Diff(tc.want, tc.got, EquateLength()); diff != "" {
				t.Errorf("Rect diff (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tc.wantString, tc.got.String()); diff != "" {
				t.Errorf("String diff (-want +got):\n%s", diff)
			}

			if got := tc.got.IsEmpty(); got != tc.wantEmpty {
				t.Errorf("%v.IsEmpty() = %v, want %v", tc.got, got, tc.wantEmpty)
			}

			if diff := cmp.Diff(tc.wantWidth, tc.got.Width(), EquateLength()); diff != "" {
				t.Errorf("Rect width diff (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tc.wantHeight, tc.got.Height(), EquateLength()); diff != "" {
				t.Errorf("Rect height diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestRectAsProto(t *testing.T) {
	for _, tc := range []struct {
		value Rect
		unit  LengthUnit
		want  *geometrypb.Rect
	}{
		{
			unit: Pt,
			want: testutil.MustUnmarshalTextproto(t,
				`top { pt: 0 } right { pt: 0 } bottom { pt: 0 } left { pt: 0 }`,
				&geometrypb.Rect{}),
		},
		{
			value: RectFromCentimeters(4, 1, 2, 3),
			unit:  Mm,
			want: testutil.MustUnmarshalTextproto(t,
				`top { mm: 10 } right { mm: 20 } bottom { mm: 30 } left { mm: 40 }`,
				&geometrypb.Rect{}),
		},
	} {
		t.Run(tc.value.String(), func(t *testing.T) {
			got := tc.value.AsProto(tc.unit)

			if diff := cmp.Diff(tc.want, got,
				protocmp.Transform(),
				cmpopts.EquateApprox(0.01, 0),
			); diff != "" {
				t.Errorf("AsProto() diff (-want +got):\n%s", diff)
			}

			if restored, err := RectFromProto(got); err != nil {
				t.Errorf("RectFromProto() failed: %v", err)
			} else if diff := cmp.Diff(tc.value, restored, EquateLength()); diff != "" {
				t.Errorf("RectFromProto() diff (-want +got):\n%s", diff)
			}

		})
	}
}

func TestRectContains(t *testing.T) {
	for _, tc := range []struct {
		input Rect
		other Rect
		want  bool
	}{
		{
			want: true,
		},
		{
			input: RectFromPoints(100, 100, 200, 200),
			other: RectFromPoints(120, 130, 140, 150),
			want:  true,
		},
		{
			other: RectFromPoints(12, 13, 14, 15),
		},
		{
			input: RectFromPoints(-100, -100, 100, 100),
			other: RectFromPoints(-101, -100, 100, 100),
		},
		{
			input: RectFromPoints(-100, -100, 100, 100),
			other: RectFromPoints(-100, -100, 100, 120),
		},
		{
			input: RectFromPoints(-10, -11, 12, 13),
			other: RectFromPoints(-1, -2, 3, 4),
			want:  true,
		},
	} {
		t.Run(tc.input.String(), func(t *testing.T) {
			validateRectangles(t, tc.input, tc.other)

			if got := tc.input.Contains(tc.other); got != tc.want {
				t.Errorf("%s.Contains(%s) = %v, want %v", tc.input, tc.other, got, tc.want)
			}

			if got := tc.other.Inside(tc.input); got != tc.want {
				t.Errorf("%s.Inside(%s) = %v, want %v", tc.other, tc.input, got, tc.want)
			}
		})
	}
}

func TestRectUnion(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input Rect
		other Rect
		want  Rect
	}{
		{
			name: "empty",
		},
		{
			name:  "input is larger",
			input: RectFromPoints(100, 100, 200, 200),
			other: RectFromPoints(120, 130, 140, 150),
			want:  RectFromPoints(100, 100, 200, 200),
		},
		{
			name:  "other is larger",
			input: RectFromPoints(120, 130, 140, 150),
			other: RectFromPoints(100, 100, 200, 200),
			want:  RectFromPoints(100, 100, 200, 200),
		},
		{
			name:  "negative and empty",
			input: RectFromPoints(-6, -5, -4, -3),
			want:  RectFromPoints(-6, -5, 0, 0),
		},
		{
			name:  "cross",
			input: RectFromPoints(40, 10, 50, 60),
			other: RectFromPoints(10, 40, 60, 50),
			want:  RectFromPoints(10, 10, 60, 60),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			validateRectangles(t, tc.input, tc.other, tc.want)

			got := tc.input.Union(tc.other)

			if diff := cmp.Diff(tc.want, got, EquateLength()); diff != "" {
				t.Errorf("Rect diff (-want +got):\n%s", diff)
			}

			validateRectangles(t, got)
		})
	}
}
