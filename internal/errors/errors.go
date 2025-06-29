package errors

import "errors"

var (
	ErrOrderNotFound   = errors.New("order not found")
	ErrInvalidOrderUID = errors.New("invalid order UID")
	ErrOrderExists     = errors.New("order already exists")
)
