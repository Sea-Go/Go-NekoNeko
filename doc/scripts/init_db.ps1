# PostgreSQL 数据库初始化脚本（PowerShell）
# 注：需要安装 PostgreSQL 客户端工具（psql）

param(
    [string]$Host = "127.0.0.1",
    [int]$Port = 5432,
    [string]$Username = "postgres",
    [string]$Password = "123456",
    [string]$Database = "favorite_db"
)

$ErrorActionPreference = "Stop"

Write-Host "PostgreSQL 数据库初始化脚本" -ForegroundColor Cyan
Write-Host "================================================" -ForegroundColor Cyan
Write-Host ""

# 检查 psql 是否安装
Write-Host "检查 psql 工具..."
$psqlPath = (Get-Command psql -ErrorAction SilentlyContinue).Path

if (-not $psqlPath) {
    Write-Host "❌ 错误: psql 未找到，请确保已安装 PostgreSQL 客户端工具" -ForegroundColor Red
    Write-Host ""
    Write-Host "Windows 安装 PostgreSQL 客户端:" -ForegroundColor Yellow
    Write-Host "  1. 从 https://www.postgresql.org/download/windows/ 下载安装程序"
    Write-Host "  2. 仅选择 'Command Line Tools' 组件安装"
    Write-Host "  3. 将 PostgreSQL bin 目录添加到 PATH 环境变量"
    exit 1
}

Write-Host "✅ psql 已找到: $psqlPath" -ForegroundColor Green
Write-Host ""

# 设置环境变量
$env:PGPASSWORD = $Password

Write-Host "连接信息:"
Write-Host "  Host: $Host"
Write-Host "  Port: $Port"
Write-Host "  User: $Username"
Write-Host "  Database: $Database"
Write-Host ""

try {
    # 步骤 1: 测试连接
    Write-Host "步骤 1: 测试数据库连接..." -ForegroundColor Yellow
    $testQuery = "SELECT version();"
    
    $output = psql -h $Host -p $Port -U $Username -d $Database -c $testQuery 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ 数据库连接成功" -ForegroundColor Green
    } else {
        Write-Host "❌ 数据库连接失败" -ForegroundColor Red
        Write-Host "错误信息: $output"
        exit 1
    }
    
    Write-Host ""
    
    # 步骤 2: 创建 auth_user 表
    Write-Host "步骤 2: 创建 auth_user 表..." -ForegroundColor Yellow
    $sqlAuthUser = @"
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
"@
    
    $sqlAuthUser | psql -h $Host -p $Port -U $Username -d $Database
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ auth_user 表已创建" -ForegroundColor Green
    } else {
        Write-Host "❌ 创建 auth_user 表失败" -ForegroundColor Red
        exit 1
    }
    
    Write-Host ""
    
    # 步骤 3: 创建 favorite_folder 表
    Write-Host "步骤 3: 创建 favorite_folder 表..." -ForegroundColor Yellow
    $sqlFolder = @"
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
"@
    
    $sqlFolder | psql -h $Host -p $Port -U $Username -d $Database
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ favorite_folder 表已创建" -ForegroundColor Green
    } else {
        Write-Host "❌ 创建 favorite_folder 表失败" -ForegroundColor Red
        exit 1
    }
    
    Write-Host ""
    
    # 步骤 4: 创建 favorite_item 表
    Write-Host "步骤 4: 创建 favorite_item 表..." -ForegroundColor Yellow
    $sqlItem = @"
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
"@
    
    $sqlItem | psql -h $Host -p $Port -U $Username -d $Database
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ favorite_item 表已创建" -ForegroundColor Green
    } else {
        Write-Host "❌ 创建 favorite_item 表失败" -ForegroundColor Red
        exit 1
    }
    
    Write-Host ""
    
    # 步骤 5: 验证表
    Write-Host "步骤 5: 验证表结构..." -ForegroundColor Yellow
    
    $tables = @("auth_user", "favorite_folder", "favorite_item")
    foreach ($table in $tables) {
        $query = "SELECT COUNT(*) FROM information_schema.tables WHERE table_name='$table';"
        $result = psql -h $Host -p $Port -U $Username -d $Database -t -c $query
        
        if ([int]$result -gt 0) {
            Write-Host "✅ $table 表已验证" -ForegroundColor Green
        } else {
            Write-Host "❌ $table 表验证失败" -ForegroundColor Red
            exit 1
        }
    }
    
    Write-Host ""
    Write-Host "================================================" -ForegroundColor Green
    Write-Host "✨ 数据库初始化完成！" -ForegroundColor Green
    Write-Host "================================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "已创建的表:"
    Write-Host "  ✓ auth_user (用户表)"
    Write-Host "  ✓ favorite_folder (收藏夹表)"
    Write-Host "  ✓ favorite_item (收藏项表)"
    Write-Host ""
    
} catch {
    Write-Host "❌ 错误: $_" -ForegroundColor Red
    exit 1
} finally {
    $env:PGPASSWORD = ""
}
