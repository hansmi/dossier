package sketch

import (
	"errors"
	"fmt"
	"sort"

	"go.uber.org/multierr"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var errInvalidTags = errors.New("invalid tags")

func validateTags(tags []string) ([]string, error) {
	tags = slices.Clone(tags)

	sort.Strings(tags)

	duplicates := map[string]struct{}{}

	var errEmpty, errDuplicates error

	for idx, i := range tags {
		if i == "" {
			errEmpty = errors.New("empty tags are forbidden")
		} else if idx > 0 && tags[idx-1] == i {
			duplicates[i] = struct{}{}
		}
	}

	if names := maps.Keys(duplicates); len(names) > 0 {
		sort.Strings(names)
		errDuplicates = fmt.Errorf("duplicated tags %q", names)
	}

	if err := multierr.Combine(errEmpty, errDuplicates); err != nil {
		return nil, fmt.Errorf("%w: %s", errInvalidTags, err.Error())
	}

	return tags, nil
}
