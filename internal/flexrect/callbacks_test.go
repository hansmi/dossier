package flexrect

import (
	"errors"
	"fmt"

	"github.com/hansmi/dossier/pkg/geometry"
	"github.com/hansmi/dossier/proto/sketchpb"
)

var errUnknownNode = errors.New("unknown node")

type fakeCallbacks struct {
	features map[NodeFeature]geometry.Point
}

func (c *fakeCallbacks) NodeFeaturePosition(name string, feature sketchpb.NodeFeature) (geometry.Point, error) {
	key := NodeFeature{
		name:    name,
		feature: feature,
	}

	pos, ok := c.features[key]
	if !ok {
		return geometry.Point{}, fmt.Errorf("%w: %#v", errUnknownNode, key)
	}

	return pos, nil
}
