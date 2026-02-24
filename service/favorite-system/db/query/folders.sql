-- name: CreateFolder :one
INSERT INTO folders (user_id, name, is_public)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetFolderByID :one
SELECT *
FROM folders
WHERE id = $1
  AND deleted_at IS NULL;

-- name: ListFoldersByUser :many
SELECT *
FROM folders
WHERE user_id = $1
  AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: SoftDeleteFolder :exec
UPDATE folders
SET deleted_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;