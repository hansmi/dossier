package sketch

import (
	"runtime"

	"github.com/sourcegraph/conc/stream"
)

func mapOrFirstError[T any, R any](items []T, fn func(T) (R, error)) ([]R, error) {
	switch len(items) {
	case 0:
		return nil, nil
	case 1:
		if r, err := fn(items[0]); err != nil {
			return nil, err
		} else {
			return []R{r}, nil
		}
	}

	var result = make([]R, len(items))
	var resultErr error

	s := stream.New().WithMaxGoroutines(runtime.GOMAXPROCS(0))

	for idx, i := range items {
		idx, i := idx, i
		s.Go(func() stream.Callback {
			r, err := fn(i)

			return func() {
				if err == nil {
					result[idx] = r
				} else if resultErr == nil {
					resultErr = err
				}
			}
		})
	}

	s.Wait()

	return result, resultErr
}
