package favorite_item

import "errors"

var (
	// ErrFolderNotFound 收藏夹不存在
	ErrFolderNotFound = errors.New("folder not found")

	// ErrFolderNotOwned 收藏夹不属于当前用户
	ErrFolderNotOwned = errors.New("folder not owned by current user")

	// ErrItemAlreadyExists 收藏项已存在（同一用户的同一对象）
	ErrItemAlreadyExists = errors.New("favorite item already exists")

	// ErrItemNotFound 收藏项不存在
	ErrItemNotFound = errors.New("favorite item not found")

	// ErrPermissionDenied 权限不足
	ErrPermissionDenied = errors.New("permission denied")

	// ErrDuplicateFavorite 重复收藏同一对象
	ErrDuplicateFavorite = errors.New("cannot favorite same object twice")
)
