package sketcherror

import "errors"

var ErrIncompleteConfig = errors.New("incomplete configuration")
var ErrBadConfig = errors.New("bad configuration")
