CREATE TABLE IF NOT EXISTS favorite_items (
    id BIGSERIAL PRIMARY KEY,
    folder_id BIGINT NOT NULL REFERENCES folders(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL,
    object_type VARCHAR(50) NOT NULL,
    object_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

CREATE UNIQUE INDEX idx_favorite_items_user_object_unique
ON favorite_items(user_id, object_type, object_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_favorite_items_folder_id
ON favorite_items(folder_id)
WHERE deleted_at IS NULL;
