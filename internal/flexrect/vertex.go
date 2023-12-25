package flexrect

import (
	"fmt"
	"strings"

	"github.com/hansmi/dossier/internal/sketcherror"
	"github.com/hansmi/dossier/pkg/geometry"
	"github.com/hansmi/dossier/proto/sketchpb"
)

type genericVertex interface {
	fmt.Stringer

	Position(Callbacks) (geometry.Point, error)
}

type absoluteVertex struct {
	name string
	pos  geometry.Point
}

var _ genericVertex = (*absoluteVertex)(nil)

func (v *absoluteVertex) String() string {
	return fmt.Sprintf("%s (position %s)", v.name, v.pos.String())
}

func (v *absoluteVertex) Position(_ Callbacks) (geometry.Point, error) {
	return v.pos, nil
}

type relativeVertex struct {
	name    string
	feature NodeFeature
	offset  geometry.Size
}

var _ genericVertex = (*relativeVertex)(nil)

func (v *relativeVertex) String() string {
	var buf strings.Builder

	fmt.Fprintf(&buf, "%s (feature %q", v.name, v.feature.String())

	if v.offset != (geometry.Size{}) {
		fmt.Fprintf(&buf, ", offset %s", v.offset.String())
	}

	buf.WriteString(")")

	return buf.String()
}

func (v *relativeVertex) Position(cb Callbacks) (geometry.Point, error) {
	pos, err := v.feature.get(cb)
	if err != nil {
		return geometry.Point{}, err
	}

	return pos.Shift(v.offset), nil
}

func vertexFromProto(pb *sketchpb.FlexRect_Vertex, name string) (genericVertex, error) {
	var err error

	switch m := pb.GetMethod().(type) {
	case *sketchpb.FlexRect_Vertex_Abs:
		v := &absoluteVertex{
			name: name,
		}

		if v.pos, err = geometry.PointFromProto(m.Abs); err != nil {
			return nil, err
		}

		return v, nil

	case *sketchpb.FlexRect_Vertex_Rel:
		rel := &relativeVertex{
			name: name,
		}

		if rel.feature, err = newNodeFeature(m.Rel); err != nil {
			return nil, err
		}

		if rel.offset, err = geometry.SizeFromProto(m.Rel.GetOffset()); err != nil {
			return nil, err
		}

		return rel, nil
	}

	return nil, fmt.Errorf("%w: vertex %q requires absolute or relative position", sketcherror.ErrIncompleteConfig, name)
}
