package errordef

import (
	"errors"
)

// Dummy
var ErrDummy = &DummyError{}

// NotFound
var ErrNotFound = &NotFoundError{}

// Uninitialized
var ErrUninitialized = errors.New("uninitialized")

// InvalidParameter
var ErrInvalidParameter = errors.New("invalid parameter")
