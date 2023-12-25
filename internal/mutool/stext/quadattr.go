package stext

import (
	"encoding/xml"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/hansmi/dossier/pkg/geometry"
)

type quadAttr geometry.Rect

var _ xml.UnmarshalerAttr = (*quadAttr)(nil)

func (a *quadAttr) UnmarshalXMLAttr(attr xml.Attr) error {
	fields := strings.Fields(attr.Value)

	if len(fields) != 8 {
		return fmt.Errorf("quad requires 8 values, got %d: %q", len(fields), attr.Value)
	}

	var err error
	var ul, ur, ll, lr struct {
		x, y float64
	}

	for idx, dest := range []*float64{
		&ul.x, &ul.y,
		&ur.x, &ur.y,
		&ll.x, &ll.y,
		&lr.x, &lr.y,
	} {
		*dest, err = strconv.ParseFloat(fields[idx], 64)
		if err != nil {
			return err
		}
	}

	// Convert quad to rect
	x0 := math.Min(math.Min(ul.x, ur.x), math.Min(ll.x, lr.x))
	y0 := math.Min(math.Min(ul.y, ur.y), math.Min(ll.y, lr.y))
	x1 := math.Max(math.Max(ul.x, ur.x), math.Max(ll.x, lr.x))
	y1 := math.Max(math.Max(ul.y, ur.y), math.Max(ll.y, lr.y))

	*a = quadAttr(geometry.RectFromPoints(x0, y0, x1, y1))

	return nil
}
