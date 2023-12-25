package sketch

import (
	"github.com/hansmi/dossier/pkg/geometry"
	"github.com/hansmi/dossier/proto/reportpb"
)

// DocumentReport is the result of a document analysis.
type DocumentReport struct {
	tags  []string
	pages []*PageReport
}

func (r *DocumentReport) Tags() []string {
	return r.tags
}

func (r *DocumentReport) Pages() []*PageReport {
	return r.pages
}

func (r *DocumentReport) AsProto(unit geometry.LengthUnit) *reportpb.Document {
	pb := &reportpb.Document{
		Tags: r.tags,
	}

	for _, p := range r.pages {
		pb.Pages = append(pb.Pages, p.AsProto(unit))
	}

	return pb
}
