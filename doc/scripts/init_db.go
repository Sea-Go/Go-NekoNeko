//go:build tools
// +build tools

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// 数据库初始化脚本
// 使用方法: go run init_db.go
func main() {
	// 数据库连接字符串
	// 需要修改为你的实际连接信息
	dsn := os.Getenv("PG_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:123456@127.0.0.1:5432/favorite_db?sslmode=disable"
	}

	fmt.Println("开始初始化数据库...")
	fmt.Printf("连接字符串: %s\n\n", dsn)

	// 连接数据库
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer pool.Close()

	// 测试连接
	ctx := context.Background()
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("数据库连接测试失败: %v", err)
	}
	fmt.Println("✅ 数据库连接成功\n")

	// 创建 auth_user 表（如果不存在）
	fmt.Println("步骤 1: 创建 auth_user 表...")
	if err := createAuthUserTable(ctx, pool); err != nil {
		log.Fatalf("创建 auth_user 表失败: %v", err)
	}
	fmt.Println("✅ auth_user 表已创建\n")

	// 创建 favorite_folder 表
	fmt.Println("步骤 2: 创建 favorite_folder 表...")
	if err := createFavoriteFolderTable(ctx, pool); err != nil {
		log.Fatalf("创建 favorite_folder 表失败: %v", err)
	}
	fmt.Println("✅ favorite_folder 表已创建\n")

	// 创建 favorite_item 表
	fmt.Println("步骤 3: 创建 favorite_item 表...")
	if err := createFavoriteItemTable(ctx, pool); err != nil {
		log.Fatalf("创建 favorite_item 表失败: %v", err)
	}
	fmt.Println("✅ favorite_item 表已创建\n")

	// 验证表
	fmt.Println("步骤 4: 验证表结构...")
	if err := verifyTables(ctx, pool); err != nil {
		log.Fatalf("验证表失败: %v", err)
	}

	fmt.Println("\n✨ 数据库初始化完成！")
	fmt.Println("\n已创建的表:")
	fmt.Println("  - auth_user (用户表)")
	fmt.Println("  - favorite_folder (收藏夹表)")
	fmt.Println("  - favorite_item (收藏项表)")
}

// createAuthUserTable 创建用户表（如果不存在）
func createAuthUserTable(ctx context.Context, pool *pgxpool.Pool) error {
	sql := `
	CREATE TABLE IF NOT EXISTS auth_user (
		id BIGSERIAL PRIMARY KEY,
		username VARCHAR(128) NOT NULL UNIQUE,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		deleted_at TIMESTAMP DEFAULT NULL
	);
	
	CREATE INDEX IF NOT EXISTS idx_username ON auth_user(username);
	CREATE INDEX IF NOT EXISTS idx_email ON auth_user(email);
	CREATE INDEX IF NOT EXISTS idx_deleted_at ON auth_user(deleted_at);
	`

	if _, err := pool.Exec(ctx, sql); err != nil {
		return err
	}

	return nil
}

// createFavoriteFolderTable 创建收藏夹表
func createFavoriteFolderTable(ctx context.Context, pool *pgxpool.Pool) error {
	sql := `
	CREATE TABLE IF NOT EXISTS favorite_folder (
		id BIGSERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL,
		name VARCHAR(64) NOT NULL,
		is_public BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		deleted_at TIMESTAMP DEFAULT NULL,
		
		CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES auth_user(id) ON DELETE CASCADE
	);
	
	CREATE UNIQUE INDEX IF NOT EXISTS uk_user_name
	ON favorite_folder (user_id, name)
	WHERE deleted_at IS NULL;
	
	CREATE INDEX IF NOT EXISTS idx_user_id
	ON favorite_folder (user_id);
	
	CREATE INDEX IF NOT EXISTS idx_deleted_at
	ON favorite_folder (deleted_at);
	`

	if _, err := pool.Exec(ctx, sql); err != nil {
		return err
	}

	return nil
}

// createFavoriteItemTable 创建收藏项表
func createFavoriteItemTable(ctx context.Context, pool *pgxpool.Pool) error {
	sql := `
	CREATE TABLE IF NOT EXISTS favorite_item (
		id BIGSERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL,
		folder_id BIGINT NOT NULL,
		object_type VARCHAR(32) NOT NULL,
		object_id BIGINT NOT NULL,
		sort_order INT DEFAULT 0,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		deleted_at TIMESTAMP DEFAULT NULL,
		
		CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES auth_user(id) ON DELETE CASCADE,
		CONSTRAINT fk_folder FOREIGN KEY (folder_id) REFERENCES favorite_folder(id) ON DELETE CASCADE,
		CONSTRAINT uk_user_object UNIQUE (user_id, object_type, object_id) WHERE deleted_at IS NULL
	);
	
	CREATE INDEX IF NOT EXISTS idx_folder_id ON favorite_item(folder_id);
	CREATE INDEX IF NOT EXISTS idx_user_id ON favorite_item(user_id);
	CREATE INDEX IF NOT EXISTS idx_object ON favorite_item(object_type, object_id);
	CREATE INDEX IF NOT EXISTS idx_created_at ON favorite_item(created_at);
	CREATE INDEX IF NOT EXISTS idx_deleted_at ON favorite_item(deleted_at);
	`

	if _, err := pool.Exec(ctx, sql); err != nil {
		return err
	}

	return nil
}

// verifyTables 验证表是否创建成功
func verifyTables(ctx context.Context, pool *pgxpool.Pool) error {
	tables := []string{"auth_user", "favorite_folder", "favorite_item"}

	for _, table := range tables {
		var exists bool
		err := pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM information_schema.tables 
				WHERE table_name = $1
			)
		`, table).Scan(&exists)

		if err != nil {
			return fmt.Errorf("查询表 %s 失败: %w", table, err)
		}

		if !exists {
			return fmt.Errorf("表 %s 不存在", table)
		}

		fmt.Printf("✅ %s 表已验证\n", table)
	}

	return nil
}
