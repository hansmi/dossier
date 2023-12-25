package renderformat

import "fmt"

type Renderer interface {
	fmt.Stringer
}
