package template

import "github.com/a-h/templ"

type TopNavItem int

const (
	TopNavNone TopNavItem = iota
	TopNavOverview
)

type BaseData struct {
	HeadTitle    string
	Scripts      []string
	TopNavActive TopNavItem
	Messages     []string
	Sidebar      templ.Component
	Content      templ.Component
}
