package template

//go:generate -command templ go run -mod=readonly github.com/a-h/templ/cmd/templ
//go:generate templ --version
//go:generate templ fmt .
//go:generate templ generate
