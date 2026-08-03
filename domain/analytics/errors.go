package analytics

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrKeyNotFound  = errors.New("key not found")
)
