package sketch

import (
	"context"
	"fmt"
	"slices"

	"github.com/hansmi/dossier"
	"github.com/hansmi/dossier/internal/flexrect"
	"github.com/hansmi/dossier/internal/sketcherror"
	"github.com/hansmi/dossier/pkg/pagerange"
	"github.com/hansmi/dossier/proto/sketchpb"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/encoding/prototext"
)

type Sketch struct {
	tags        []string
	nodes       []*sketchNode
	searchOrder []int
}

func Compile(pb *sketchpb.Sketch) (*Sketch, error) {
	var err error

	s := &Sketch{
		nodes: make([]*sketchNode, 0, len(pb.GetNodes())),
	}

	if s.tags, err = validateTags(pb.GetTags()); err != nil {
		return nil, multierr.Combine(sketcherror.ErrBadConfig, err)
	}

	for _, pn := range pb.GetNodes() {
		n, err := sketchNodeFromProto(pn)
		if err != nil {
			return nil, fmt.Errorf("node %s: %w", pn.GetName(), err)
		}

		s.nodes = append(s.nodes, n)
	}

	if order, err := determineNodeOrder(s.nodes); err != nil {
		return nil, err
	} else {
		s.searchOrder = order
	}

	return s, nil
}

func CompileFromTextproto(b []byte) (*Sketch, error) {
	var msg sketchpb.Sketch

	if err := prototext.Unmarshal(b, &msg); err != nil {
		return nil, err
	}

	return Compile(&msg)
}

func CompileFromTextprotoString(s string) (*Sketch, error) {
	return CompileFromTextproto([]byte(s))
}

func (s *Sketch) Tags() []string {
	return s.tags
}

func (s *Sketch) AnalyzePage(p *dossier.Page) (*PageReport, error) {
	r := newPageReport(p)

	callbacks := struct {
		documentPage
		flexrect.Callbacks
	}{
		documentPage: p,
		Callbacks:    r,
	}

	for _, i := range s.searchOrder {
		node := s.nodes[i]

		match, err := node.search(&callbacks)
		if err != nil {
			return nil, fmt.Errorf("search for node %q on page %d: %w", node.name, r.Number(), err)
		}

		r.appendNode(match)
	}

	return r, nil
}

func (s *Sketch) AnalyzePages(pages []*dossier.Page) ([]*PageReport, error) {
	return mapOrFirstError(pages, s.AnalyzePage)
}

func (s *Sketch) AnalyzeDocument(ctx context.Context, doc *dossier.Document, r pagerange.Range) (*DocumentReport, error) {
	pages, err := doc.ParsePages(ctx, r)
	if err != nil {
		return nil, err
	}

	result := &DocumentReport{
		tags: slices.Clone(s.tags),
	}

	if result.pages, err = s.AnalyzePages(pages); err != nil {
		return nil, err
	}

	return result, nil
}
