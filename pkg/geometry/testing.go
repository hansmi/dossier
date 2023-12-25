package geometry

import (
	"github.com/google/go-cmp/cmp"
)

// EquateLength returns a [cmp.Comparer] option checking whether two distances
// have a margin of less than one millimeter.
func EquateLength() cmp.Option {
	return EquateLengthApprox(Millimeter)
}

// EquateLengthApprox returns a [cmp.Comparer] option checking whether two
// distances are differ less than the specified margin.
func EquateLengthApprox(margin Length) cmp.Option {
	return cmp.FilterValues(func(x, y Length) bool {
		return true
	}, cmp.Comparer(func(x, y Length) bool {
		return (x - y).Abs() <= margin
	}))
}
