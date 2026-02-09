# 完整集成测试指南

## 📋 目录

1. [环境准备](#环境准备)
2. [数据库初始化](#数据库初始化)
3. [启动服务](#启动服务)
4. [运行测试](#运行测试)
5. [测试场景说明](#测试场景说明)
6. [常见问题](#常见问题)

---

## 环境准备

### 前置条件

- Go 1.18+ （用于编译和运行服务）
- PostgreSQL 12+ （可选，可用 Docker）
- curl 或 Postman （用于 API 测试）

### 安装 PostgreSQL（Docker）

```bash
cd d:\UGit\Sea-TryGo-feature-collect-system

# 启动 PostgreSQL 容器
docker-compose up -d

# 验证容器运行
docker-compose ps
```

---

## 数据库初始化

### 方法 1: 使用 Go 脚本（推荐）

```bash
cd doc/scripts

# 运行初始化脚本
go run init_db.go

# 如果需要指定数据库连接
$env:PG_DSN = "postgres://postgres:yourpassword@host:5432/favorite_db?sslmode=disable"
go run init_db.go
```

**输出示例**:
```
开始初始化数据库...
连接字符串: postgres://postgres:123456@127.0.0.1:5432/favorite_db?sslmode=disable

✅ 数据库连接成功

步骤 1: 创建 auth_user 表...
✅ auth_user 表已创建

步骤 2: 创建 favorite_folder 表...
✅ favorite_folder 表已创建

步骤 3: 创建 favorite_item 表...
✅ favorite_item 表已创建

步骤 4: 验证表结构...
✅ auth_user 表已验证
✅ favorite_folder 表已验证
✅ favorite_item 表已验证

✨ 数据库初始化完成！
```

### 方法 2: 使用 PowerShell 脚本

```powershell
cd doc\scripts
.\init_db.ps1 -Host "127.0.0.1" -Port 5432 -Username "postgres" -Password "123456" -Database "favorite_db"
```

### 方法 3: 手动使用 psql

```bash
# 连接数据库
psql -h 127.0.0.1 -U postgres -d favorite_db

# 运行SQL脚本（在psql中）
\i 'doc/sql/favorite_item.sql'
\i 'doc/sql/favorite_folder.sql'
```

---

## 启动服务

### 编译

```bash
cd d:\UGit\Sea-TryGo-feature-collect-system

# 编译 favorite API 服务
go build -o ./api_favorite ./api/favorite.go

# 编译 usercenter API 服务（可选）
go build -o ./api_usercenter ./api/usercenter.go
```

### 运行服务

```bash
# 方法 1: 直接运行编译后的二进制
.\api_favorite

# 输出应该类似:
# [INFO] starting server on 0.0.0.0:8888
```

### 验证服务运行

```bash
# 在新的 PowerShell 窗口中测试
curl -v http://localhost:8888/health 2>&1 | Select-Object -First 10
```

---

## 运行测试

### 步骤 1: 生成 JWT Token

```bash
cd api\tools
go run jwt_generator.go

# 输出:
# JWT Token (用于 Authorization header 中):
# Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDQxMTExNTUsImlhdCI6MTcwNDAyNDc1NSwidXNlcl9pZCI6MX0.abc123...
```

**复制完整 token（包括 "Bearer " 前缀**）

### 步骤 2: 运行集成测试

#### 使用 PowerShell 脚本

```powershell
cd doc\scripts

# 简单版本（使用已有 token）
.\integration_test_simple.ps1 `
    -BaseUrl "http://localhost:8888" `
    -Token "Bearer <your-jwt-token-here>"

# 详细版本（自动生成 token）
.\integration_test.ps1 `
    -BaseUrl "http://localhost:8888" `
    -JwtSecret "favorite-secret-key"
```

#### 使用 curl 手动测试

```bash
# 设置变量
$TOKEN = "Bearer <your-jwt-token-here>"
$BASE_URL = "http://localhost:8888"

# 1. 创建收藏夹
curl -X POST "$BASE_URL/favorite/v1/folders" `
  -H "Content-Type: application/json" `
  -H "Authorization: $TOKEN" `
  -d '{"name":"My Favorites","is_public":false}' | ConvertFrom-Json

# 2. 创建收藏项
curl -X POST "$BASE_URL/favorite/v1/items" `
  -H "Content-Type: application/json" `
  -H "Authorization: $TOKEN" `
  -d '{"folder_id":1,"object_type":"article","object_id":"12345"}' | ConvertFrom-Json

# 3. 列表收藏项
curl -X GET "$BASE_URL/favorite/v1/items?folder_id=1&page=1&page_size=10" `
  -H "Authorization: $TOKEN" | ConvertFrom-Json

# 4. 删除收藏项
curl -X DELETE "$BASE_URL/favorite/v1/items" `
  -H "Content-Type: application/json" `
  -H "Authorization: $TOKEN" `
  -d '{"object_type":"article","object_id":"12345"}' | ConvertFrom-Json
```

---

## 测试场景说明

### 场景 1️⃣: 创建收藏项（成功）

**请求**:
```bash
POST /favorite/v1/items
Content-Type: application/json
Authorization: Bearer <token>

{
  "folder_id": 1,
  "object_type": "article",
  "object_id": "12345"
}
```

**预期结果**:
- ✅ 状态码: 200
- ✅ 响应包含: id, user_id, folder_id, object_type, object_id, created_at, updated_at

**错误响应** (如果收藏夹不属于用户):
- ❌ 状态码: 403
- ❌ 错误信息: "收藏夹不属于当前用户"

### 场景 2️⃣: 重复收藏（冲突）

**请求**: 同上（两次相同请求）

**第二次预期结果**:
- ❌ 状态码: 409
- ❌ 错误信息: "该对象已被收藏"

**业务规则**: 同一用户的同一对象不能被收藏两次

### 场景 3️⃣: 列表收藏项

**请求**:
```bash
GET /favorite/v1/items?folder_id=1&page=1&page_size=10
Authorization: Bearer <token>
```

**预期结果**:
- ✅ 状态码: 200
- ✅ 响应包含: items (数组), total (总数), page, page_size

**分页说明**:
- page: 从 1 开始
- page_size: 1-100（超出范围会返回错误）
- total: 总记录数

### 场景 4️⃣: 删除收藏项

**请求**:
```bash
DELETE /favorite/v1/items
Content-Type: application/json
Authorization: Bearer <token>

{
  "object_type": "article",
  "object_id": "12345"
}
```

**预期结果**:
- ✅ 状态码: 200
- ✅ 响应: { "success": true, "message": "删除成功" }

**删除行为**: 软删除（deleted_at 设置为当前时间，不是真正删除）

### 场景 5️⃣: 无效 Token（权限错误）

**请求**: 不提供或无效的 Authorization header

**预期结果**:
- ❌ 状态码: 401
- ❌ 错误信息: "invalid or missing authorization token"

**原因包括**:
- 缺少 Authorization header
- Token 已过期
- Token 签名无效（secret 不匹配）

### 场景 6️⃣: 资源不存在（404）

**请求**:
```bash
GET /favorite/v1/items?folder_id=99999
```

**预期结果**:
- ❌ 状态码: 404
- ❌ 错误信息: "收藏夹不存在"

---

## HTTP 状态码映射

| 状态码 | 含义 | 场景 |
|--------|------|------|
| **200** | ✅ OK | 请求成功 |
| **400** | ❌ Bad Request | 参数验证失败 |
| **401** | ❌ Unauthorized | 缺少或无效的 JWT token |
| **403** | ❌ Forbidden | 权限不足（如收藏夹不属于用户） |
| **404** | ❌ Not Found | 资源不存在 |
| **409** | ❌ Conflict | 资源重复（如重复收藏） |
| **500** | ❌ Internal Server Error | 服务器错误 |

---

## 常见问题

### Q1: 无法连接到数据库

**症状**:
```
Error: connect ECONNREFUSED 127.0.0.1:5432
```

**解决方案**:
1. 确保 PostgreSQL 已启动
   ```bash
   docker-compose ps
   ```
2. 检查连接字符串中的 Host、Port、Username、Password
3. 检查数据库是否已创建
   ```bash
   psql -h 127.0.0.1 -U postgres -l | grep favorite_db
   ```

### Q2: 401 Unauthorized 错误

**症状**:
```json
{
  "code": 401,
  "message": "invalid or missing authorization token"
}
```

**解决方案**:
1. 确认 Authorization header 的格式: `Bearer <token>`（注意空格）
2. 生成新的 JWT token
   ```bash
   cd api/tools
   go run jwt_generator.go
   ```
3. 检查 token 是否过期（有效期默认 24 小时）

### Q3: 403 Forbidden 错误

**症状**:
```json
{
  "code": 403,
  "message": "收藏夹不属于当前用户"
}
```

**解决方案**:
- 确认收藏夹 ID 属于当前用户（token 中的 user_id）
- 检查收藏夹是否已被删除

### Q4: 409 Conflict 错误

**症状**:
```json
{
  "code": 409,
  "message": "该对象已被收藏"
}
```

**原因**: 同一用户已经收藏了这个对象

**解决方案**:
- 先删除现有收藏，再重新收藏
- 或使用不同的 object_id

### Q5: 如何重置数据库？

```bash
# 删除所有表
psql -h 127.0.0.1 -U postgres -d favorite_db -c "DROP TABLE IF EXISTS favorite_item CASCADE;"
psql -h 127.0.0.1 -U postgres -d favorite_db -c "DROP TABLE IF EXISTS favorite_folder CASCADE;"
psql -h 127.0.0.1 -U postgres -d favorite_db -c "DROP TABLE IF EXISTS auth_user CASCADE;"

# 重新初始化
cd doc/scripts
go run init_db.go
```

---

## 测试总结

### ✅ 成功标准

所有以下场景都应该返回正确的状态码和错误消息：

- [x] 创建收藏项（200）
- [x] 重复收藏（409）
- [x] 列表收藏项（200，包含分页）
- [x] 删除收藏项（200）
- [x] 无效 token（401）
- [x] 权限错误（403）
- [x] 资源不存在（404）

### 📊 性能基准

- 创建收藏: < 50ms
- 列表查询: < 100ms（无缓存）
- 删除操作: < 50ms
- JWT 验证: < 10ms

---

## 后续优化建议

- [ ] 添加 Redis 缓存
- [ ] 实现批量操作接口
- [ ] 添加排序功能
- [ ] 实现收藏统计 API
- [ ] 添加取消关注选项
