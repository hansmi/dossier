package sketch

import (
	"errors"
	"fmt"

	"github.com/hansmi/dossier"
	"github.com/hansmi/dossier/internal/flexrect"
	"github.com/hansmi/dossier/internal/sketcherror"
	"github.com/hansmi/dossier/pkg/geometry"
	"github.com/hansmi/dossier/proto/sketchpb"
	"go.uber.org/multierr"
)

type documentPage interface {
	VisitElementsIntersecting(geometry.Rect, dossier.PageElementVisitorFunc) error
}

type sketchNodeLocator interface {
	locate(documentPage, geometry.Rect) (func(*Node), error)
}

type sketchNode struct {
	name        string
	searchAreas []*flexrect.FlexRect
	locator     sketchNodeLocator
	tags        []string
}

func sketchNodeFromProto(pbnode *sketchpb.Node) (*sketchNode, error) {
	var err error

	node := &sketchNode{
		name: pbnode.GetName(),
	}

	if node.tags, err = validateTags(pbnode.GetTags()); err != nil {
		return nil, multierr.Combine(sketcherror.ErrBadConfig, err)
	}

	switch m := pbnode.GetMatcher().(type) {
	case *sketchpb.Node_BlockText:
		node.locator, err = newTextLocatorFromProto(m.BlockText, false)

	case *sketchpb.Node_LineText:
		node.locator, err = newTextLocatorFromProto(m.LineText, true)

	default:
		err = fmt.Errorf("%w: node %q has unsupported match type %T", sketcherror.ErrBadConfig, node.name, m)
	}

	if err != nil {
		return nil, err
	}

	for _, pbArea := range pbnode.GetSearchAreas() {
		area, err := flexrect.FromProto(pbArea)
		if err != nil {
			return nil, err
		}

		node.searchAreas = append(node.searchAreas, area)
	}

	if len(node.searchAreas) < 1 {
		return nil, fmt.Errorf("%w: node %q requires at least one search area", sketcherror.ErrIncompleteConfig, node.name)
	}

	return node, nil
}

func (s *sketchNode) featurePosition(bounds geometry.Rect, feature sketchpb.NodeFeature) (geometry.Point, error) {
	switch feature {
	case sketchpb.NodeFeature_TOP_LEFT:
		return bounds.TopLeft(), nil
	case sketchpb.NodeFeature_TOP_RIGHT:
		return bounds.TopRight(), nil
	case sketchpb.NodeFeature_BOTTOM_LEFT:
		return bounds.BottomLeft(), nil
	case sketchpb.NodeFeature_BOTTOM_RIGHT:
		return bounds.BottomRight(), nil
	}

	return geometry.Point{}, fmt.Errorf("%w: node %q lacks feature %s", ErrNodeFeatureUnavailable, s.name, feature.String())
}

type sketchNodeSearchCallbacks interface {
	documentPage
	flexrect.Callbacks
}

func (s *sketchNode) search(cb sketchNodeSearchCallbacks) (*Node, error) {
	n := &Node{s: s}

	// Find candidate areas
	for _, area := range s.searchAreas {
		bounds, err := area.Resolve(cb)
		if err != nil {
			if errors.Is(err, ErrNodePositionUnknown) {
				continue
			}

			return nil, err
		}

		n.searchAreas = append(n.searchAreas, bounds)
	}

	// Search within valid areas
	for _, area := range n.searchAreas {
		if apply, err := s.locator.locate(cb, area); err != nil {
			return nil, err
		} else if apply != nil {
			apply(n)
			break
		}
	}

	return n, nil
}
