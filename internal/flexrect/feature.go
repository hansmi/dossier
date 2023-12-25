package flexrect

import (
	"fmt"

	"github.com/hansmi/dossier/internal/sketcherror"
	"github.com/hansmi/dossier/pkg/geometry"
	"github.com/hansmi/dossier/proto/sketchpb"
)

type NodeFeature struct {
	name    string
	feature sketchpb.NodeFeature
}

func newNodeFeature(pb interface {
	GetNode() string
	GetFeature() sketchpb.NodeFeature
}) (NodeFeature, error) {
	f := NodeFeature{
		name:    pb.GetNode(),
		feature: pb.GetFeature(),
	}

	if f.name == "" {
		return NodeFeature{}, fmt.Errorf("%w: missing node name", sketcherror.ErrIncompleteConfig)
	}

	return f, nil
}

func (f *NodeFeature) String() string {
	return fmt.Sprintf("%s:%s", f.name, f.feature.String())
}

func (f *NodeFeature) NodeName() string {
	return f.name
}

func (f *NodeFeature) Feature() sketchpb.NodeFeature {
	return f.feature
}

func (f *NodeFeature) compare(other NodeFeature) int {
	if f.name < other.name {
		return -1
	} else if f.name > other.name {
		return +1
	}

	if f.feature < other.feature {
		return -1
	} else if f.feature > other.feature {
		return +1
	}

	return 0
}

func (f *NodeFeature) get(cb Callbacks) (geometry.Point, error) {
	return cb.NodeFeaturePosition(f.name, f.feature)
}
