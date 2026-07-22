package analytics

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")

	ErrKeyNotFound = errors.New("key not found")
	ErrValueIsNil  = errors.New("value is nil")
)

var ErrReadOnly = errors.New("metric is read-only")
