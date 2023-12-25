package testfiles

import (
	"embed"
)

//go:embed *.xml
//go:embed *.pdf
var All embed.FS
