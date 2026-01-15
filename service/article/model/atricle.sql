-- service/article/model/article.sql
CREATE TABLE "articles" (
                            "id" bigserial PRIMARY KEY,
                            "title" varchar(255) NOT NULL DEFAULT '',
                            "content" text NOT NULL DEFAULT '',
                            "author" varchar(255) NOT NULL DEFAULT '',
                            "status" smallint NOT NULL DEFAULT 0, -- 0: Processing, 1: Published, 2: Failed
                            "storage_url" varchar(1024) NOT NULL DEFAULT '',
                            "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);