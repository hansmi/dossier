package flexrect

import (
	"github.com/hansmi/dossier/pkg/geometry"
	"github.com/hansmi/dossier/proto/sketchpb"
)

type Callbacks interface {
	NodeFeaturePosition(name string, feature sketchpb.NodeFeature) (geometry.Point, error)
}
