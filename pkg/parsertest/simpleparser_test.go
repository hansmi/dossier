package parsertest

import (
	"context"
	"testing"

	"github.com/hansmi/dossier/pkg/pagerange"
)

func FuzzSimpleParserParsePages(f *testing.F) {
	for _, i := range []int{0, 1, 2, 100, pagerange.Last} {
		f.Add(i, i, i)
	}

	f.Fuzz(func(t *testing.T, count, start, end int) {
		var p SimpleParser

		if count > 1000 {
			count = 1000
		}

		for i := 0; i < count; i++ {
			p.Pages = append(p.Pages, nil)
		}

		if r, err := pagerange.New(start, end); err == nil {
			p.ParsePages(context.Background(), r)
		}
	})
}
