package genorm

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrNullValue      = errors.New("null value")
)
