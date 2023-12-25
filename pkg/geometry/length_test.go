package geometry

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/dossier/proto/geometrypb"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestLengthFromProto(t *testing.T) {
	for _, tc := range []struct {
		name    string
		input   *geometrypb.Length
		want    Length
		wantErr error
	}{
		{
			name: "nil",
		},
		{
			name:  "empty",
			input: &geometrypb.Length{},
		},
		{
			name: "centimeter",
			input: &geometrypb.Length{
				Value: &geometrypb.Length_Cm{Cm: 22},
			},
			want: 22 * Cm,
		},
		{
			name: "millimeter",
			input: &geometrypb.Length{
				Value: &geometrypb.Length_Mm{Mm: 123},
			},
			want: 123 * Mm,
		},
		{
			name: "inch",
			input: &geometrypb.Length{
				Value: &geometrypb.Length_In{In: 33},
			},
			want: 33 * In,
		},
		{
			name: "point",
			input: &geometrypb.Length{
				Value: &geometrypb.Length_Pt{Pt: 44},
			},
			want: 44 * Pt,
		},
		{
			name: "negative",
			input: &geometrypb.Length{
				Value: &geometrypb.Length_Pt{Pt: -1},
			},
			want: -1,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := LengthFromProto(tc.input)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, EquateLength()); diff != "" {
					t.Errorf("LengthFromProto() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestLengthString(t *testing.T) {
	for _, tc := range []struct {
		value Length
		want  string
	}{
		{
			want: "0",
		},
		{
			value: 1 * Pt,
			want:  "0.0353cm",
		},
		{
			value: Millimeter,
			want:  "0.1cm",
		},
		{
			value: 1.3 * Millimeter,
			want:  "0.13cm",
		},
		{
			value: 1.7777 * Millimeter,
			want:  "0.178cm",
		},
		{
			value: Inch,
			want:  "2.54cm",
		},
		{
			value: (-7 * Inch) + (3 * Cm),
			want:  "-14.8cm",
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

func TestLengthAsProto(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value Length
		unit  LengthUnit
		want  *geometrypb.Length
	}{
		{
			name: "empty cm",
			unit: Cm,
			want: &geometrypb.Length{
				Value: &geometrypb.Length_Cm{},
			},
		},
		{
			name: "empty pt",
			unit: Pt,
			want: &geometrypb.Length{
				Value: &geometrypb.Length_Pt{},
			},
		},
		{
			name:  "1m",
			value: 100 * Cm,
			unit:  Cm,
			want: &geometrypb.Length{
				Value: &geometrypb.Length_Cm{Cm: 100},
			},
		},
		{
			name:  "mm",
			value: Cm,
			unit:  Mm,
			want: &geometrypb.Length{
				Value: &geometrypb.Length_Mm{Mm: 10},
			},
		},
		{
			name:  "pt",
			value: 1.2 * Cm,
			unit:  Pt,
			want: &geometrypb.Length{
				Value: &geometrypb.Length_Pt{Pt: 34.016},
			},
		},
		{
			name:  "pt rounded to mm as cm",
			value: 12 * Pt,
			unit: &RoundedLength{
				Unit:    Cm,
				Nearest: Mm,
			},
			want: &geometrypb.Length{
				Value: &geometrypb.Length_Cm{Cm: 0.4},
			},
		},
		{
			name:  "negative mm rounded to inch as mm",
			value: -28 * Mm,
			unit: &RoundedLength{
				Unit:    Mm,
				Nearest: Inch,
			},
			want: &geometrypb.Length{
				Value: &geometrypb.Length_Mm{Mm: -25.4},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
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

func TestLengthRound(t *testing.T) {
	for _, tc := range []struct {
		name    string
		value   Length
		nearest Length
		want    Length
	}{
		{
			name:    "zero",
			nearest: Inch,
		},
		{
			name:    ".8in",
			nearest: Inch,
			value:   .8 * Inch,
			want:    Inch,
		},
		{
			name:    "-.8cm",
			nearest: Cm,
			value:   -.8 * Cm,
			want:    -Cm,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.value.Round(tc.nearest)

			if diff := cmp.Diff(tc.want, got, EquateLength()); diff != "" {
				t.Errorf("Round() diff (-want +got):\n%s", diff)
			}
		})
	}
}
