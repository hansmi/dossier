//go:build tools

package tooldeps

// https://github.com/golang/go/issues/48429
// https://go.dev/cl/495555
// https://www.jvt.me/posts/2022/06/15/go-tools-dependency-management/
import (
	_ "github.com/a-h/templ/cmd/templ"
)
