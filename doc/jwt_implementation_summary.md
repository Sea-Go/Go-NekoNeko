# JWT 认证实现完成总结

## ✅ 实现完成

### 1. JWT 认证工具 (`api/internal/utils/jwt.go`)

**功能**:
- 从 HTTP Authorization header 中提取 JWT token
- 验证 JWT 签名（使用配置中的 AccessSecret）
- 从 JWT claims 中提取 user_id
- 处理多种错误场景（缺失 header、无效 token、签名错误等）

**关键方法**:
```go
GetUserIDFromRequest(r *http.Request, secret string) (int64, error)
WriteErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, errMsg string)
```

**支持的 Claims 字段**:
- `user_id`: 用户 ID (支持 float64 和 string 类型)
- `exp`: 过期时间 (Unix timestamp)
- `iat`: 颁发时间 (Unix timestamp)

---

### 2. Handler 层更新

所有三个 Favorite handler 已更新为从 JWT 中提取 user_id：

**文件**:
- ✅ `api/internal/handler/favorite/createfavoritehandler.go`
- ✅ `api/internal/handler/favorite/deletefavoritehandler.go`
- ✅ `api/internal/handler/favorite/listfavoritehandler.go`

**变更内容**:
```go
// 之前（硬编码）:
userID := int64(1)

// 现在（从 JWT 中提取）:
userID, err := utils.GetUserIDFromRequest(r, serverCtx.Config.UserAuth.AccessSecret)
if err != nil {
    utils.WriteErrorResponse(w, r, http.StatusUnauthorized, "invalid or missing authorization token")
    return
}
```

---

### 3. JWT 配置

**文件**: `api/etc/favorite.yaml`

```yaml
UserAuth:
  AccessSecret: "favorite-secret-key"  # JWT 签名密钥
  AccessExpire: 7200                   # Token 有效期（秒）
```

---

### 4. 测试工具

#### 4.1 JWT Token 生成器 (`api/tools/jwt_generator.go`)

**功能**: 生成用于测试的有效 JWT token

**使用方法**:
```bash
cd api/tools
go run jwt_generator.go
```

**输出示例**:
```
JWT Token (用于 Authorization header 中):
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDQxMTExNTUsImlhdCI6MTcwNDAyNDc1NSwidXNlcl9pZCI6MX0.abc123...
```

#### 4.2 API 测试脚本 (Windows PowerShell)

**文件**: `api/tools/test_api.ps1`

**功能**: 自动化测试所有 API 端点（包括 JWT 认证）

**使用方法**:
```powershell
# 1. 生成 JWT token
cd api\tools
go run jwt_generator.go

# 2. 复制输出的完整 token (包括 "Bearer " 前缀)

# 3. 运行测试脚本
$token = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
.\test_api.ps1 -BaseUrl 'http://localhost:8888' -JwtToken $token
```

**测试场景**:
- ✅ 创建收藏项（POST /favorite/v1/items）
- ✅ 列表收藏项（GET /favorite/v1/items）
- ✅ 删除收藏项（DELETE /favorite/v1/items）
- ✅ 无效 token 验证（应返回 401）

#### 4.3 API 测试脚本 (Linux/Mac Bash)

**文件**: `api/tools/test_api.sh`

**功能**: 使用 curl 进行 API 测试

---

### 5. 文档

**文件**: `doc/jwt_authentication.md`

**内容**:
- JWT 工作原理和流程
- 配置说明
- API 测试示例
- 错误处理说明
- 安全建议
- 代码实现细节

---

## 📊 工作流程架构

```
┌─────────────────────────────────────┐
│     HTTP Request                    │
│   Authorization: Bearer <token>     │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│ Handler 层                          │
│ (CreateFavoriteHandler 等)          │
└────────────┬────────────────────────┘
             │
             ▼ GetUserIDFromRequest()
┌─────────────────────────────────────┐
│ JWT 工具 (api/internal/utils/)      │
│ 1. 提取 Authorization header        │
│ 2. 解析 JWT token                   │
│ 3. 验证签名（用 AccessSecret）     │
│ 4. 提取 user_id claim               │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│ Handler 重新实现                    │
│ (获取到 userID，继续业务逻辑)       │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│ Logic 层                            │
│ (CreateFavoriteLogic 等)            │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│ Service 层                          │
│ (FavoriteItemService)               │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│ Repository 层                       │
│ (Database 操作)                     │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│     PostgreSQL                      │
└─────────────────────────────────────┘
```

---

## 🔐 安全特性

1. **JWT 签名验证**: 使用 HS256 (HMAC) 算法
2. **Token 过期验证**: 自动检查 exp claim
3. **错误消息**: 统一的 401 错误响应，不泄露详细信息
4. **多类型支持**: user_id 可以是 float64 或 string 格式

---

## 📋 编译状态

✅ **编译成功**
```
Exit Code: 0 (成功)
```

所有文件均已编译验证，没有语法或类型错误。

---

## 🧪 测试设置说明

### 前置条件

1. **启动 PostgreSQL 数据库** (如果还没启动)
   ```powershell
   docker-compose up -d
   ```

2. **初始化数据库**
   ```sql
   -- 执行 doc/sql/ 中的所有 SQL 脚本
   ```

### 运行服务

```bash
# 交叉编译（可选）
go build -o ./api_favorite ./api/favorite.go

# 运行服务
./api_favorite
```

### 测试 API

```powershell
# 1. 生成 JWT token
cd api\tools
$token = go run jwt_generator.go

# 2. 复制标记为 "Bearer ..." 的完整 token

# 3. 运行测试脚本
.\test_api.ps1 -JwtToken "<paste-token-here>"
```

---

## 📝 关键代码片段

### JWT 提取（在 Handler 中）

```go
// 从请求中提取 userID
userID, err := utils.GetUserIDFromRequest(r, serverCtx.Config.UserAuth.AccessSecret)
if err != nil {
    utils.WriteErrorResponse(w, r, http.StatusUnauthorized, 
        "invalid or missing authorization token")
    return
}

// 继续业务逻辑，使用 userID
logic := favorite.NewCreateFavoriteLogic(serverCtx.FavoriteItemService)
item, err := logic.Execute(r.Context(), req, userID)
```

### JWT Token 生成（用于测试）

```go
claims := jwt.MapClaims{
    "user_id": int64(1),
    "exp":     time.Now().Add(time.Hour * 24).Unix(),
    "iat":     time.Now().Unix(),
}

token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
tokenString, _ := token.SignedString([]byte(secret))
```

---

## ⚠️ 已知限制和改进项

### 当前限制

1. **单一 Secret**: 所有 token 使用同一个 secret（无轮换机制）
2. **无 Refresh Token**: 生成的 token 过期后需重新认证
3. **无 Token 黑名单**: Logout 后旧 token 仍可用（直到过期）
4. **无速率限制**: 没有实现针对认证失败的限制

### 建议的后续改进

- [ ] 实现 Refresh Token 机制
- [ ] 添加 Token 黑名单（用于 Logout）
- [ ] 实现 Secret 轮换策略
- [ ] 添加速率限制 (rate limiting)
- [ ] 支持多种认证方法（OAuth2、SAML 等）
- [ ] 添加角色和权限管理 (RBAC)
- [ ] 实现 Audit 日志记录

---

## 📞 故障排除

### 问题 1: "invalid or missing authorization token"

**原因**:
- 没有提供 Authorization header
- Token 已过期
- Token 签名无效（secret 不匹配）

**解决方案**:
1. 确认是否提供了 Authorization header
2. 重新生成新的 JWT token
3. 验证 secret key 是否与 favorite.yaml 中的一致

### 问题 2: Token 生成失败

**原因**:
- golang-jwt 库未安装

**解决方案**:
```bash
go get github.com/golang-jwt/jwt/v4
```

### 问题 3: 测试脚本运行错误

**原因**:
- PowerShell 执行策略限制
- 脚本路径不正确

**解决方案**:
```powershell
# 允许运行脚本
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# 从正确的目录运行
cd api\tools
.\test_api.ps1 -JwtToken "<token>"
```

---

## 总结

✅ **JWT 认证实现完成**

所有 Favorite API 端点现在都：
1. 验证 Authorization header 中的 JWT token
2. 验证 token 的签名和过期时间
3. 从 token 中安全地提取 user_id
4. 使用 user_id 进行后续的业务逻辑处理

代码已编译验证，测试工具已提供。下一步可以进行集成测试和性能测试。
