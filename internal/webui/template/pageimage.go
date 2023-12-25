package template

import (
	"fmt"
	"math"
	"net/url"
	"strconv"

	"github.com/hansmi/dossier"
)

type PageImageData struct {
	DocFingerprint string

	Page       *dossier.Page
	Width      int
	ClassNames []string
	Alt        string
}

func (d PageImageData) imageData() ImageData {
	p := d.Page
	pageSize := p.Size()

	params := url.Values{}
	params.Set("width", strconv.Itoa(d.Width))

	if d.DocFingerprint != "" {
		params.Set("docfp", d.DocFingerprint)
	}

	return ImageData{
		Src:        fmt.Sprintf("/page/%d/image?%s", p.Number(), params.Encode()),
		Width:      d.Width,
		Height:     int(math.Ceil(float64(d.Width) * pageSize.Height.Pt() / pageSize.Width.Pt())),
		ClassNames: d.ClassNames,
		Alt:        d.Alt,
	}
}
