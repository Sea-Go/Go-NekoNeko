package folder

import "errors"

var (
	ErrInvalidUserID = errors.New("invalid user_id")
	ErrInvalidName   = errors.New("invalid name")
)
