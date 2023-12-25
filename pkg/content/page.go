package content

import (
	"github.com/hansmi/dossier/pkg/geometry"
)

type Page interface {
	// 1-based page number.
	Number() int

	// Physical page size.
	Size() geometry.Size

	Elements() []Element
}
