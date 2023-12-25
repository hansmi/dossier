package stext

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"github.com/hansmi/dossier/pkg/geometry"
)

type rectAttr geometry.Rect

var _ xml.UnmarshalerAttr = (*rectAttr)(nil)

func (a *rectAttr) UnmarshalXMLAttr(attr xml.Attr) error {
	fields := strings.Fields(attr.Value)

	if len(fields) != 4 {
		return fmt.Errorf("rect requires 4 values, got %d: %q", len(fields), attr.Value)
	}

	for idx, dest := range []*geometry.Length{&a.Left, &a.Top, &a.Right, &a.Bottom} {
		value, err := strconv.ParseFloat(fields[idx], 64)
		if err != nil {
			return err
		}

		*dest = geometry.Pt.Mul(value)
	}

	return nil
}
