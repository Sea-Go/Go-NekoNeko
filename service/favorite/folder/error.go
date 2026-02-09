package folder

import "errors"

var (
	ErrFolderNameEmpty   = errors.New("folder name is empty")
	ErrFolderNameTooLong = errors.New("folder name too long")
	ErrFolderNameExists  = errors.New("folder name already exists")
)
