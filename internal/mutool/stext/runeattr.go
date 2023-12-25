package stext

import (
	"encoding/xml"
	"fmt"
)

type runeAttr rune

var _ xml.UnmarshalerAttr = (*runeAttr)(nil)

func (a *runeAttr) UnmarshalXMLAttr(attr xml.Attr) error {
	runes := []rune(attr.Value)

	if len(runes) > 1 {
		return fmt.Errorf("expected exactly one rune, got %q", attr.Value)
	}

	if len(runes) == 0 {
		*a = ' '
	} else {
		*a = runeAttr(runes[0])
	}

	return nil
}
