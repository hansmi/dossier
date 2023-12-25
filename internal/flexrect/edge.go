package flexrect

import (
	"fmt"
	"strings"

	"github.com/hansmi/dossier/internal/sketcherror"
	"github.com/hansmi/dossier/pkg/geometry"
	"github.com/hansmi/dossier/proto/sketchpb"
)

type genericEdge interface {
	fmt.Stringer

	Position(Callbacks) (geometry.Length, error)
}

type absoluteEdge struct {
	name     string
	distance geometry.Length
}

var _ genericEdge = (*absoluteEdge)(nil)

func (e *absoluteEdge) String() string {
	if e.distance == 0 {
		return e.name
	}

	return fmt.Sprintf("%s (distance %s)", e.name, e.distance.String())
}

func (e *absoluteEdge) Position(_ Callbacks) (geometry.Length, error) {
	return e.distance, nil
}

type relativeEdge struct {
	name    string
	feature NodeFeature
	offset  geometry.Length

	extract pointDimensionFunc
}

var _ genericEdge = (*relativeEdge)(nil)

func (e *relativeEdge) String() string {
	var buf strings.Builder

	fmt.Fprintf(&buf, "%s (feature %q", e.name, e.feature.String())

	if e.offset != 0 {
		fmt.Fprintf(&buf, ", offset %s", e.offset.String())
	}

	buf.WriteString(")")

	return buf.String()
}

func (e *relativeEdge) Position(cb Callbacks) (geometry.Length, error) {
	pos, err := e.feature.get(cb)
	if err != nil {
		return 0, err
	}

	return e.extract(pos) + e.offset, nil
}

func edgeFromProto(pb *sketchpb.FlexRect_Edge, name string, extract pointDimensionFunc) (genericEdge, error) {
	var err error

	switch m := pb.GetMethod().(type) {
	case *sketchpb.FlexRect_Edge_Abs:
		e := &absoluteEdge{
			name: name,
		}

		if e.distance, err = geometry.LengthFromProto(m.Abs); err != nil {
			return nil, err
		}

		return e, nil

	case *sketchpb.FlexRect_Edge_Rel:
		rel := &relativeEdge{
			name:    name,
			extract: extract,
		}

		if rel.feature, err = newNodeFeature(m.Rel); err != nil {
			return nil, err
		}

		if m.Rel.Offset != nil {
			if rel.offset, err = geometry.LengthFromProto(m.Rel.GetOffset()); err != nil {
				return nil, err
			}
		}

		return rel, nil
	}

	return nil, fmt.Errorf("%w: edge %q requires absolute or relative position", sketcherror.ErrIncompleteConfig, name)
}

type shiftedEdge struct {
	edge   genericEdge
	offset geometry.Length
}

var _ genericEdge = (*shiftedEdge)(nil)

func (e *shiftedEdge) String() string {
	if e.offset == 0 {
		return e.edge.String()
	}

	return fmt.Sprintf("%s from %s", e.offset.String(), e.edge.String())
}

func (e *shiftedEdge) Position(cb Callbacks) (geometry.Length, error) {
	pos, err := e.edge.Position(cb)
	if err != nil {
		return 0, err
	}

	return pos + e.offset, nil
}

type edgeFromVertex struct {
	vertex  genericVertex
	extract pointDimensionFunc
}

var _ genericEdge = (*edgeFromVertex)(nil)

func (e *edgeFromVertex) String() string {
	return e.vertex.String()
}

func (e *edgeFromVertex) Position(cb Callbacks) (geometry.Length, error) {
	pos, err := e.vertex.Position(cb)
	if err != nil {
		return 0, err
	}

	return e.extract(pos), nil
}
