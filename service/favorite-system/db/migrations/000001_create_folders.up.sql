CREATE TABLE IF NOT EXISTS folders (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,

    is_public BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- 用户 + 名称 唯一（未删除情况下）
CREATE UNIQUE INDEX idx_folders_user_name_unique
ON folders(user_id, name)
WHERE deleted_at IS NULL;

-- 查询优化
CREATE INDEX idx_folders_user_id
ON folders(user_id)
WHERE deleted_at IS NULL;