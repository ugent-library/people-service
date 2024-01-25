package models

import (
	"errors"
)

var ErrNotFound = errors.New("not found")
var ErrFatal = errors.New("fatal error")
var ErrSkipped = errors.New("was skipped")
var ErrMissingArgument = errors.New("missing argument")
var ErrInvalidReference = errors.New("invalid reference")
var ErrInvalidURN = errors.New("invalid urn")
