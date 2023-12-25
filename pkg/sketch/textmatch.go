package sketch

import (
	"fmt"
	"regexp"

	"github.com/hansmi/dossier/internal/ref"
	"github.com/hansmi/dossier/proto/reportpb"
	"golang.org/x/exp/slices"
)

type TextMatchGroup struct {
	// Capture group name. Empty if no name is set using (?P<...>).
	Name string

	// Zero-based start and end offset of the group in the original string.
	Start int
	End   int

	// Text captured by the group.
	Text string
}

func (g *TextMatchGroup) AsProto() *reportpb.TextMatchGroup {
	return &reportpb.TextMatchGroup{
		Name:  g.Name,
		Start: int32(g.Start),
		End:   int32(g.End),
		Text:  g.Text,
	}
}

// TextMatch captures information about a regular expression match in a string.
type TextMatch struct {
	expr    *regexp.Regexp
	text    string
	indexes []int
	names   []string
}

// evaluateMatch tries to match a regular expression pattern in a string.
// Returns nil if no match is found.
func evaluateMatch(pat *regexp.Regexp, text string) *TextMatch {
	indexes := pat.FindStringSubmatchIndex(text)

	if len(indexes) == 0 {
		return nil
	}

	return &TextMatch{
		expr:    pat,
		text:    text,
		indexes: indexes,
		names:   pat.SubexpNames(),
	}
}

// Pattern returns the pattern source text for the matched regular expression.
func (m *TextMatch) Pattern() string {
	return m.expr.String()
}

func (m *TextMatch) group(idx int) TextMatchGroup {
	g := TextMatchGroup{
		Name:  m.names[idx],
		Start: m.indexes[idx*2],
		End:   m.indexes[idx*2+1],
	}
	if g.Start >= 0 && g.End >= 0 {
		g.Text = m.text[g.Start:g.End]
	}

	return g
}

// Groups returns a slice with all subgroups of the match.
func (m *TextMatch) Groups() []TextMatchGroup {
	groups := make([]TextMatchGroup, len(m.indexes)/2)

	for idx := range groups {
		groups[idx] = m.group(idx)
	}

	return groups
}

// Group returns a subgroup by index or nil if it's not captured.
func (m *TextMatch) Group(idx int) *TextMatchGroup {
	if idx >= len(m.indexes)/2 {
		return nil
	}

	return ref.Ref(m.group(idx))
}

// MustGroup returns a subgroup by index or panics if it's not captured.
func (m *TextMatch) MustGroup(idx int) TextMatchGroup {
	g := m.Group(idx)
	if g == nil {
		panic(fmt.Sprintf("Group %d not captured by %q", idx, m.expr.String()))
	}

	return *g
}

// Group returns a subgroup by name or nil if it's not captured.
func (m *TextMatch) Named(name string) *TextMatchGroup {
	var idx int

	if name != "" {
		idx = slices.Index(m.names, name)
		if idx < 0 {
			return nil
		}
	}

	return ref.Ref(m.group(idx))
}

// Group returns a subgroup by name or panics if it's not captured.
func (m *TextMatch) MustNamed(name string) TextMatchGroup {
	g := m.Named(name)
	if g == nil {
		panic(fmt.Sprintf("Group %q not captured by %q", name, m.expr.String()))
	}

	return *g
}
