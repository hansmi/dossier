package geometry

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/dossier/internal/testutil"
	"github.com/hansmi/dossier/proto/geometrypb"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestPointString(t *testing.T) {
	for _, tc := range []struct {
		value Point
		want  string
	}{
		{
			want: "(0, 0)",
		},
		{
			value: Point{1 * Pt, 2 * Pt},
			want:  "(0.0353cm, 0.0706cm)",
		},
		{
			value: Point{-1 * Mm, -2 * Cm},
			want:  "(-0.1cm, -2cm)",
		},
	} {
		t.Run(tc.want, func(t *testing.T) {
			got := tc.value.String()

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("String() diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPointAsProto(t *testing.T) {
	for _, tc := range []struct {
		value Point
		unit  LengthUnit
		want  *geometrypb.Point
	}{
		{
			unit: Pt,
			want: testutil.MustUnmarshalTextproto(t,
				`left { pt: 0 } top { pt: 0 }`,
				&geometrypb.Point{}),
		},
		{
			value: Point{1 * Pt, 2 * Pt},
			unit:  Cm,
			want: testutil.MustUnmarshalTextproto(t,
				`left { cm: 0.0353 } top { cm: 0.0706 }`,
				&geometrypb.Point{}),
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
		})
	}
}

func TestPointFromProto(t *testing.T) {
	for _, tc := range []struct {
		name    string
		input   *geometrypb.Point
		want    Point
		wantErr error
	}{
		{
			name: "nil",
		},
		{
			name:  "empty",
			input: &geometrypb.Point{},
		},
		{
			name: "positive",
			input: testutil.MustUnmarshalTextproto(t,
				`left { cm: 3 } top { mm: 17 }`,
				&geometrypb.Point{}),
			want: Point{Left: 3 * Cm, Top: 1.7 * Cm},
		},
		{
			name: "negative",
			input: testutil.MustUnmarshalTextproto(t,
				`left { mm: -9 } top { mm: -10.1 }`,
				&geometrypb.Point{}),
			want: Point{Left: -9 * Mm, Top: -10.1 * Mm},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := PointFromProto(tc.input)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, EquateLength()); diff != "" {
					t.Errorf("PointFromProto() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
