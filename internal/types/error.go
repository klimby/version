package types

import "errors"

var (
	ErrInvalidArguments = errors.New("invalid arguments")
	ErrNotInitialized   = errors.New("not initialized")
)
