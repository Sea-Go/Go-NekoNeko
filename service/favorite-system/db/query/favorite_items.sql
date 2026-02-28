-- name: AddFavoriteItem :one
INSERT INTO favorite_items (folder_id, user_id, object_type, object_id, title)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFavoriteItem :one
SELECT * FROM favorite_items
WHERE user_id = $1 AND object_type = $2 AND object_id = $3 AND deleted_at IS NULL;

-- name: ListFavoriteItems :many
SELECT * FROM favorite_items
WHERE folder_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListAllFavoriteItems :many
SELECT * FROM favorite_items
WHERE folder_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: CountFavoriteItems :one
SELECT COUNT(*) FROM favorite_items
WHERE folder_id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteFavoriteItem :exec
UPDATE favorite_items
SET deleted_at = NOW()
WHERE folder_id = $1 AND user_id = $2 AND object_type = $3 AND object_id = $4 AND deleted_at IS NULL;

-- name: SoftDeleteFavoriteItemsByFolder :exec
UPDATE favorite_items
SET deleted_at = NOW()
WHERE folder_id = $1 AND deleted_at IS NULL;
