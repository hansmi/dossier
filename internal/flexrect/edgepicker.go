package flexrect

import (
	"fmt"

	"github.com/hansmi/dossier/internal/sketcherror"
	"github.com/hansmi/dossier/pkg/geometry"
)

type edgePicker struct {
	name    string
	sources []genericEdge
}

func (p *edgePicker) addEdge(edge genericEdge, offset geometry.Length) {
	if offset != 0 {
		edge = &shiftedEdge{
			edge:   edge,
			offset: offset,
		}
	}

	p.sources = append(p.sources, edge)
}

func (p *edgePicker) addVertex(v genericVertex, extract pointDimensionFunc, offset geometry.Length) {
	var edge genericEdge = &edgeFromVertex{
		vertex:  v,
		extract: extract,
	}

	if offset != 0 {
		edge = &shiftedEdge{
			edge:   edge,
			offset: offset,
		}
	}

	p.sources = append(p.sources, edge)
}

func (p *edgePicker) pick() (genericEdge, error) {
	if len(p.sources) < 1 {
		return nil, fmt.Errorf("%w: no position specified for %q", sketcherror.ErrIncompleteConfig, p.name)
	}

	if len(p.sources) > 1 {
		var desc []string

		for _, i := range p.sources {
			desc = append(desc, i.String())
		}

		return nil, fmt.Errorf("%w: conflicting position sources for %q: %q", sketcherror.ErrBadConfig, p.name, desc)
	}

	return p.sources[0], nil
}
