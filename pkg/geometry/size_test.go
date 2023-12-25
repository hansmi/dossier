package geometry

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/dossier/internal/testutil"
	"github.com/hansmi/dossier/proto/geometrypb"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestSizeString(t *testing.T) {
	for _, tc := range []struct {
		value Size
		want  string
	}{
		{
			want: "(0, 0)",
		},
		{
			value: Size{1 * Pt, 2 * Pt},
			want:  "(0.0353cm, 0.0706cm)",
		},
		{
			value: Size{-1 * Mm, -2 * Cm},
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

func TestSizeAsProto(t *testing.T) {
	for _, tc := range []struct {
		value Size
		unit  LengthUnit
		want  *geometrypb.Size
	}{
		{
			unit: Pt,
			want: testutil.MustUnmarshalTextproto(t,
				`width { pt: 0 } height { pt: 0 }`,
				&geometrypb.Size{}),
		},
		{
			value: Size{1 * Pt, 2 * Pt},
			unit:  Mm,
			want: testutil.MustUnmarshalTextproto(t,
				`width { mm: 0.353 } height { mm: 0.706 }`,
				&geometrypb.Size{}),
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

func TestSizeFromProto(t *testing.T) {
	for _, tc := range []struct {
		name    string
		input   *geometrypb.Size
		want    Size
		wantErr error
	}{
		{
			name: "nil",
		},
		{
			name:  "empty",
			input: &geometrypb.Size{},
		},
		{
			name: "cm",
			input: testutil.MustUnmarshalTextproto(t,
				`width { cm: 123 } height { cm: 27572 }`,
				&geometrypb.Size{}),
			want: Size{Width: 123 * Cm, Height: 27572 * Cm},
		},
		{
			name: "negative",
			input: testutil.MustUnmarshalTextproto(t,
				`width { mm: -1.1 } height { mm: -3.4 }`,
				&geometrypb.Size{}),
			want: Size{Width: -1.1 * Mm, Height: -3.4 * Mm},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := SizeFromProto(tc.input)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, EquateLength()); diff != "" {
					t.Errorf("SizeFromProto() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
