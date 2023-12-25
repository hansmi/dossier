package sketch

import (
	"fmt"

	"github.com/hansmi/dossier"
	"github.com/hansmi/dossier/pkg/geometry"
	"github.com/hansmi/dossier/proto/reportpb"
	"github.com/hansmi/dossier/proto/sketchpb"
)

type PageReport struct {
	num    int
	size   geometry.Size
	nodes  []*Node
	byName map[string]*Node
}

func newPageReport(p *dossier.Page) *PageReport {
	return &PageReport{
		num:    p.Number(),
		size:   p.Size(),
		byName: map[string]*Node{},
	}
}

func (p *PageReport) appendNode(n *Node) {
	p.nodes = append(p.nodes, n)
	p.byName[n.Name()] = n
}

// Number returns the 1-based page number.
func (p *PageReport) Number() int {
	return p.num
}

func (p *PageReport) Size() geometry.Size {
	return p.size
}

func (p *PageReport) Nodes() []*Node {
	return p.nodes
}

func (p *PageReport) NodeByName(name string) *Node {
	return p.byName[name]
}

func (p *PageReport) NodeFeaturePosition(name string, feature sketchpb.NodeFeature) (geometry.Point, error) {
	n := p.NodeByName(name)
	if n == nil {
		return geometry.Point{}, fmt.Errorf("%w: node %q not found", ErrNodePositionUnknown, name)
	}

	return n.FeaturePosition(feature)
}

func (p *PageReport) AsProto(unit geometry.LengthUnit) *reportpb.Page {
	pb := &reportpb.Page{
		Size:   p.size.AsProto(unit),
		Number: int32(p.num),
	}

	for _, n := range p.nodes {
		pb.Nodes = append(pb.Nodes, n.AsProto(unit))
	}

	return pb
}
