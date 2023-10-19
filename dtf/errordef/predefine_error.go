package errordef

import (
	"errors"
)

// Dummy
var ErrDummy = &DummyError{}

// NotFound
var ErrNotFound = &NotFoundError{}

var ErrUninitialized = errors.New("uninitialized")
var ErrInvalidParameter = errors.New("invalid parameter")
var ErrOperationFailed = errors.New("operation failed")
var ErrInternalError = errors.New("internal error")
var ErrAccessDenied = errors.New("access denied")
