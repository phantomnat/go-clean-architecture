package article

import "errors"

var (
	ErrInternalServer = errors.New("internal server error")
	ErrNotFound       = errors.New("your requested item is not found")
	ErrAlreadyExist   = errors.New("your item already exist")
)
