package submsg

import (
	"errors"
)

var (
	ErrUnknownOneOfField = errors.New("unknown oneof field")
)
