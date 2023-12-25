package sketch

import (
	"errors"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestMapOrFirstError(t *testing.T) {
	errTest := errors.New("test error")

	for _, tc := range []struct {
		name    string
		items   []string
		fn      func(string) (int, error)
		wantErr error
		want    []int
	}{
		{name: "empty"},
		{
			name:  "one",
			items: []string{"aaa"},
			fn: func(s string) (int, error) {
				return len(s), nil
			},
			want: []int{3},
		},
		{
			name:  "three items",
			items: []string{"first", "second", "third"},
			fn: func(s string) (int, error) {
				return len(s), nil
			},
			want: []int{5, 6, 5},
		},
		{
			name:  "error on fifth item",
			items: []string{"a", "b", "c", "d", "e", "f", "g"},
			fn: func(s string) (int, error) {
				if s == "e" {
					return 0, errTest
				}

				return int(s[0]), nil
			},
			wantErr: errTest,
		},
		{
			name: "hundreds of items",
			items: (func() (result []string) {
				for i := 0; i < 700; i++ {
					result = append(result, strconv.Itoa(i))
				}
				return
			})(),
			fn: func(s string) (int, error) {
				return strconv.Atoi(s)
			},
			want: (func() (result []int) {
				for i := 0; i < 700; i++ {
					result = append(result, i)
				}
				return
			})(),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := mapOrFirstError(tc.items, tc.fn)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("analyzeInParallel() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
