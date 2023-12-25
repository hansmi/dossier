package template

import (
	"github.com/hansmi/dossier"
)

type OverviewContentData struct {
	DocFingerprint string
	Pages          []*dossier.Page
}

type OverviewSidebarData struct {
	Path        string
	Size        string
	ModTime     string
	ModTimeFull string
}
