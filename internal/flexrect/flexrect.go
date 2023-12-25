package flexrect

import (
	"fmt"

	"github.com/hansmi/dossier/internal/ref"
	"github.com/hansmi/dossier/pkg/geometry"
	"github.com/hansmi/dossier/proto/geometrypb"
	"github.com/hansmi/dossier/proto/sketchpb"
	"go.uber.org/multierr"
	"golang.org/x/exp/slices"
)

type edgePositionFunc func(Callbacks) (geometry.Length, error)

type dependencyDiscovery struct {
	nodes []NodeFeature
}

func (d *dependencyDiscovery) NodeFeaturePosition(name string, feature sketchpb.NodeFeature) (geometry.Point, error) {
	d.nodes = append(d.nodes, NodeFeature{
		name:    name,
		feature: feature,
	})
	return geometry.Point{}, nil
}

func (d *dependencyDiscovery) get() []NodeFeature {
	slices.SortFunc(d.nodes, func(a, b NodeFeature) int {
		return a.compare(b)
	})

	return d.nodes
}

type FlexRect struct {
	left   edgePositionFunc
	top    edgePositionFunc
	right  edgePositionFunc
	bottom edgePositionFunc

	deps []NodeFeature
}

func FromProto(pb *sketchpb.FlexRect) (*FlexRect, error) {
	var top, right, bottom, left genericEdge
	var topLeft, topRight, bottomLeft, bottomRight genericVertex
	var width, height *geometry.Length

	var err, errAll error

	for _, i := range []struct {
		name       string
		dst        *genericEdge
		src        *sketchpb.FlexRect_Edge
		pointCoord pointDimensionFunc
	}{
		{"top", &top, pb.Top, pointTop},
		{"right", &right, pb.Right, pointLeft},
		{"bottom", &bottom, pb.Bottom, pointTop},
		{"left", &left, pb.Left, pointLeft},
	} {
		if i.src != nil {
			if *i.dst, err = edgeFromProto(i.src, i.name, i.pointCoord); err != nil {
				multierr.AppendInto(&errAll, err)
			}
		}
	}

	for _, i := range []struct {
		name string
		dst  *genericVertex
		src  *sketchpb.FlexRect_Vertex
	}{
		{"top_left", &topLeft, pb.TopLeft},
		{"top_right", &topRight, pb.TopRight},
		{"bottom_left", &bottomLeft, pb.BottomLeft},
		{"bottom_right", &bottomRight, pb.BottomRight},
	} {
		if i.src != nil {
			if *i.dst, err = vertexFromProto(i.src, i.name); err != nil {
				multierr.AppendInto(&errAll, err)
			}
		}
	}

	for _, i := range []struct {
		name string
		dst  **geometry.Length
		src  *geometrypb.Length
	}{
		{"width", &width, pb.Width},
		{"height", &height, pb.Height},
	} {
		if i.src != nil {
			if length, err := geometry.LengthFromProto(i.src); err != nil {
				multierr.AppendInto(&errAll, fmt.Errorf("%s: %w", i.name, err))
			} else {
				*i.dst = ref.Ref(length)
			}
		}
	}

	if errAll != nil {
		return nil, errAll
	}

	edges := []genericEdge{top, right, bottom, left}
	vertices := []genericVertex{topLeft, topRight, bottomRight, bottomLeft}
	distances := []*geometry.Length{height, width}
	discovery := dependencyDiscovery{}

	r := &FlexRect{}

	for idx, i := range []struct {
		name    string
		dst     *edgePositionFunc
		edge    genericEdge
		extract pointDimensionFunc
	}{
		{"top", &r.top, top, pointTop},
		{"right", &r.right, right, pointLeft},
		{"bottom", &r.bottom, bottom, pointTop},
		{"left", &r.left, left, pointLeft},
	} {
		p := edgePicker{
			name: i.name,
		}

		if i.edge != nil {
			p.addEdge(i.edge, 0)
		}

		for _, j := range []int{idx, (idx + 1) % len(vertices)} {
			if v := vertices[j]; v != nil {
				p.addVertex(v, i.extract, 0)
			}
		}

		if distp := distances[idx%len(distances)]; distp != nil {
			dist := *distp

			if i.name == "left" || i.name == "top" {
				dist = -dist
			}

			if oppositeEdge := edges[(idx+2)%len(edges)]; oppositeEdge != nil {
				p.addEdge(oppositeEdge, dist)
			}

			// Opposite vertices with distance
			for _, j := range []int{
				(idx + 2) % len(vertices),
				(idx + 3) % len(vertices),
			} {
				if v := vertices[j]; v != nil {
					p.addVertex(v, i.extract, dist)
				}
			}
		}

		src, err := p.pick()
		if err != nil {
			return nil, err
		}

		*i.dst = src.Position

		if _, err := src.Position(&discovery); err != nil {
			return nil, err
		}
	}

	r.deps = discovery.get()

	return r, nil
}

// RequiredNodeFeatures returns the all node features referenced by dimensions.
// These must be resolvable before the FlexRect can be calculated.
func (r *FlexRect) RequiredNodeFeatures() []NodeFeature {
	return r.deps
}

// Resolve calculates the absolute position of the rectangle. Callbacks are
// invoked as necessary.
func (r *FlexRect) Resolve(cb Callbacks) (geometry.Rect, error) {
	var result geometry.Rect
	var err, errAll error

	for _, i := range []struct {
		name string
		dst  *geometry.Length
		src  edgePositionFunc
	}{
		{"top", &result.Top, r.top},
		{"right", &result.Right, r.right},
		{"bottom", &result.Bottom, r.bottom},
		{"left", &result.Left, r.left},
	} {
		if *i.dst, err = i.src(cb); err != nil {
			multierr.AppendInto(&errAll, fmt.Errorf("%s: %w", i.name, err))
		}
	}

	if errAll != nil {
		return geometry.Rect{}, errAll
	}

	result = result.Normalize()

	if err := result.Validate(); err != nil {
		return geometry.Rect{}, fmt.Errorf("flexrect: %w", err)
	}

	return result, nil
}
