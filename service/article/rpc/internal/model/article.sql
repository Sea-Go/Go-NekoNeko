-- Article Table Definition
CREATE TABLE IF NOT EXISTS "article" (
    "id" varchar(32) NOT NULL,
    "title" varchar(255) NOT NULL,
    "brief" varchar(512),
    "content" text,
    "cover_image_url" varchar(255),
    "manual_type_tag" varchar(64),
    "secondary_tags" jsonb,
    "author_id" varchar(32),
    "status" smallint DEFAULT 0,
    "view_count" integer DEFAULT 0,
    "like_count" integer DEFAULT 0,
    "comment_count" integer DEFAULT 0,
    "share_count" integer DEFAULT 0,
    "ext_info" jsonb,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    PRIMARY KEY ("id")
);

-- Indexes
CREATE INDEX IF NOT EXISTS "idx_article_manual_type_tag" ON "article" ("manual_type_tag");
CREATE INDEX IF NOT EXISTS "idx_article_author_id" ON "article" ("author_id");
CREATE INDEX IF NOT EXISTS "idx_article_deleted_at" ON "article" ("deleted_at");

-- Comments (Optional but recommended for documentation)
COMMENT ON TABLE "article" IS '文章主表';
COMMENT ON COLUMN "article"."id" IS '文章ID (Snowflake)';
COMMENT ON COLUMN "article"."content" IS 'Markdown正文';
COMMENT ON COLUMN "article"."secondary_tags" IS '二级标签 (JSON Array)';
COMMENT ON COLUMN "article"."status" IS '状态: 0-未指定, 1-草稿, 2-已发布, 3-审核中, 4-拒绝';
COMMENT ON COLUMN "article"."ext_info" IS '扩展字段 (JSON Map)';