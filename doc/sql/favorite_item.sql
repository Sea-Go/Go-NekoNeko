-- 收藏项表（PostgreSQL，存放用户收藏的具体内容）
CREATE TABLE favorite_item (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  folder_id BIGINT NOT NULL,
  -- 收藏的对象信息
  object_type VARCHAR(32) NOT NULL,  -- 例如："article", "video", "post" 等
  object_id BIGINT NOT NULL,         -- 对应对象的 ID
  -- 排序字段
  sort_order INT DEFAULT 0,
  -- 时间戳
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP DEFAULT NULL,
  -- 约束：用户对同一对象只能在一个收藏夹中存在一份
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES auth_user(id) ON DELETE CASCADE,
  CONSTRAINT fk_folder FOREIGN KEY (folder_id) REFERENCES favorite_folder(id) ON DELETE CASCADE,
  -- 唯一约束：(user_id, object_type, object_id) 组合唯一
  CONSTRAINT uk_user_object UNIQUE (user_id, object_type, object_id) WHERE deleted_at IS NULL
);

-- 索引优化
CREATE INDEX idx_folder_id ON favorite_item(folder_id);
CREATE INDEX idx_user_id ON favorite_item(user_id);
CREATE INDEX idx_object ON favorite_item(object_type, object_id);
CREATE INDEX idx_created_at ON favorite_item(created_at);
CREATE INDEX idx_deleted_at ON favorite_item(deleted_at);
