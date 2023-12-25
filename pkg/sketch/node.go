package sketch

import (
	"github.com/hansmi/dossier/pkg/geometry"
	"github.com/hansmi/dossier/proto/reportpb"
	"github.com/hansmi/dossier/proto/sketchpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Node struct {
	s           *sketchNode
	valid       bool
	bounds      geometry.Rect
	searchAreas []geometry.Rect
	text        *string
	textMatch   *TextMatch
}

func (n *Node) Name() string {
	return n.s.name
}

func (n *Node) Tags() []string {
	return n.s.tags
}

func (n *Node) Valid() bool {
	return n.valid
}

func (n *Node) Bounds() geometry.Rect {
	return n.bounds
}

func (n *Node) SearchAreas() []geometry.Rect {
	return n.searchAreas
}

func (n *Node) FeaturePosition(feature sketchpb.NodeFeature) (geometry.Point, error) {
	if !n.valid {
		return geometry.Point{}, ErrNodePositionUnknown
	}

	return n.s.featurePosition(n.bounds, feature)
}

func (n *Node) Text() string {
	if n.text != nil {
		return *n.text
	}

	return ""
}

func (n *Node) TextMatch() *TextMatch {
	return n.textMatch
}

func (n *Node) AsProto(unit geometry.LengthUnit) *reportpb.Node {
	pb := &reportpb.Node{
		Name:  n.s.name,
		Valid: n.Valid(),
		Tags:  n.s.tags,
	}

	for _, area := range n.searchAreas {
		pb.SearchAreas = append(pb.SearchAreas, area.AsProto(unit))
	}

	if pb.GetValid() {
		pb.Bounds = n.bounds.AsProto(unit)

		if n.text != nil {
			pb.Text = wrapperspb.String(*n.text)
		}

		if n.textMatch != nil {
			for _, g := range n.textMatch.Groups() {
				pb.TextMatchGroups = append(pb.TextMatchGroups, g.AsProto())
			}
		}
	}

	return pb
}
