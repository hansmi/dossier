package sketch

import (
	"errors"

	"github.com/hansmi/dossier/internal/sketcherror"
)

var (
	ErrIncompleteConfig = sketcherror.ErrIncompleteConfig
	ErrBadConfig        = sketcherror.ErrBadConfig

	ErrNodeFeatureUnavailable = errors.New("node feature unavailable")
	ErrNodePositionUnknown    = errors.New("node position unknown")
)
