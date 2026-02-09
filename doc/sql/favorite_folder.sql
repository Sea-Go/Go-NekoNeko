-- 收藏夹表（PostgreSQL，支持软删除）
CREATE TABLE favorite_folder (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  name VARCHAR(64) NOT NULL,
  is_public BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP DEFAULT NULL
);
-- 同一用户下收藏夹名唯一（只针对未删除的）
CREATE UNIQUE INDEX uk_user_name
ON favorite_folder (user_id, name)
WHERE deleted_at IS NULL;
CREATE INDEX idx_user_id
ON favorite_folder (user_id);
CREATE INDEX idx_deleted_at
ON favorite_folder (deleted_at);
