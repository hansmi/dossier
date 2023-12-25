package mutool

import (
	"fmt"
	"strconv"

	"github.com/hansmi/dossier/pkg/pagerange"
)

func formatPageRange(r pagerange.Range) string {
	if r.Lower == pagerange.Last {
		return "N"
	}

	var upper string

	if r.Upper == pagerange.Last {
		upper = "N"
	} else {
		upper = strconv.Itoa(r.Upper)
	}

	return fmt.Sprintf("%d-%s", r.Lower, upper)
}
