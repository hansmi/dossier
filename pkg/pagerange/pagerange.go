package pagerange

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

const Last int = math.MaxInt

var All = MustNew(1, Last)

// Ranges are one-based and include the last page (e.g. {1,3} means page 1,
// 2 and 3). Upper may be [Last] to indicate all pages from lower to the last.
// If both values are [Last] only the last page is covered.
type Range struct {
	Lower, Upper int
}

// New creates a new page range.
func New(lower, upper int) (Range, error) {
	// switch {
	// case lower < 1:
	// 	return fmt.Errorf("%w: lower end must be >=1, got %d", os.ErrInvalid, lower)

	// case lower == Last && upper == Last:
	// }

	// if lower < 1 {
	// 	return Range{}, fmt.Errorf("%w: lower end must be >=1, got %d", os.ErrInvalid, lower)
	// }

	// if upper < lower {
	// 	return Range{}, fmt.Errorf("%w: upper end must be same as lower (%d) or larger, got %d", os.ErrInvalid, lower, upper)
	// }

	r := Range{lower, upper}

	if err := r.Validate(); err != nil {
		return Range{}, err
	}

	return r, nil
}

// MustNew is like [New] but panics if the range is invalid.
func MustNew(lower, upper int) Range {
	return must1(New(lower, upper))
}

// Single returns a range for a single page.
func Single(n int) (Range, error) {
	return New(n, n)
}

// MustSingle is like [Single] but panics if construction fails.
func MustSingle(n int) Range {
	return must1(Single(n))
}

// Validate checks whether the range is valid.
func (r Range) Validate() error {
	if r.Lower < 1 {
		return fmt.Errorf("%w: lower end must be >=1, got %d",
			os.ErrInvalid, r.Lower)
	}

	if r.Upper < r.Lower {
		return fmt.Errorf("%w: upper end must be same as lower (%d) or larger, got %d",
			os.ErrInvalid, r.Lower, r.Upper)
	}

	return nil
}

func (r Range) String() string {
	if r.Lower == Last {
		return "(last)"
	}

	if r.Lower == r.Upper {
		return strconv.Itoa(r.Lower)
	}

	var upper string

	if r.Upper == Last {
		upper = "(last)"
	} else {
		upper = strconv.Itoa(r.Upper)
	}

	return fmt.Sprintf("%d-%s", r.Lower, upper)
}
