# JWT 认证实现指南

## 概览

本项目采用 JWT (JSON Web Token) 进行用户身份认证。JWT token 在请求的 `Authorization` header 中传递。

## 工作原理

### 1. JWT Token 生成

JWT token 包含以下信息：
- **user_id**: 用户唯一标识符
- **exp**: token 过期时间（Unix timestamp）
- **iat**: token 颁发时间（Unix timestamp）

### 2. 认证流程

```
HTTP Request
    ↓
Authorization Header (Bearer <token>)
    ↓
JWT 验证（使用 AccessSecret）
    ↓
Extract user_id from claims
    ↓
执行业务逻辑
```

### 3. 配置

JWT 配置位于 `api/etc/favorite.yaml`：

```yaml
UserAuth:
  AccessSecret: "favorite-secret-key"  # 用于签名和验证 JWT
  AccessExpire: 7200                   # token 有效期（秒）
```

## 使用方法

### 生成测试 JWT Token

```bash
cd api/tools
go run jwt_generator.go
```

输出示例：
```
JWT Token (用于 Authorization header 中):
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 测试 API 端点

#### 1. 使用 curl 添加收藏

```bash
curl -X POST http://localhost:8888/favorite/v1/items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "folder_id": 1,
    "object_type": "article",
    "object_id": "12345"
  }'
```

#### 2. 列表收藏项

```bash
curl -X GET "http://localhost:8888/favorite/v1/items?folder_id=1&page=1&page_size=10" \
  -H "Authorization: Bearer <your-jwt-token>"
```

#### 3. 删除收藏

```bash
curl -X DELETE http://localhost:8888/favorite/v1/items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "object_type": "article",
    "object_id": "12345"
  }'
```

## 错误处理

### 401 Unauthorized

**场景**: Authorization header 缺失或 token 无效

**响应**:
```json
{
  "code": 401,
  "message": "invalid or missing authorization token"
}
```

**常见原因**:
- 没有提供 Authorization header
- Token 已过期
- Token 签名无效
- Secret key 不匹配

### 200 OK

**场景**: 请求成功

**响应示例** (创建收藏):
```json
{
  "id": 1,
  "user_id": 1,
  "folder_id": 1,
  "object_type": "article",
  "object_id": "12345",
  "sort_order": 0,
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

## 代码实现

### JWT 提取工具 (`api/internal/utils/jwt.go`)

```go
// GetUserIDFromRequest 从 HTTP 请求中提取 userID
func GetUserIDFromRequest(r *http.Request, secret string) (int64, error) {
  // 1. 获取 Authorization header
  // 2. 解析 JWT token
  // 3. 验证签名（使用 secret）
  // 4. 提取 user_id claim
  // 5. 返回 userID 或错误
}
```

### 中间件集成

所有收藏 API 端点自动配置了 JWT 中间件：

```go
// api/internal/handler/routes.go
routes := []http.Route{
	{
		Method: http.MethodPost,
		Path:   "/favorite/v1/items",
		Handler: CreateFavoriteHandler(serverCtx),
		Middlewares: []http.Middleware{
			rest.WithJwt(serverCtx.Config.UserAuth.AccessSecret),
		},
	},
	// ... 更多路由
}
```

### Handler 中的 JWT 提取

```go
// api/internal/handler/favorite/createfavoritehandler.go
userID, err := utils.GetUserIDFromRequest(r, serverCtx.Config.UserAuth.AccessSecret)
if err != nil {
  utils.WriteErrorResponse(w, r, http.StatusUnauthorized, 
    "invalid or missing authorization token")
  return
}
```

## 测试场景

### 场景 1: 有效 token，创建收藏

**步骤**:
1. 生成有效的 JWT token（user_id=1）
2. 发送 POST 请求到 `/favorite/v1/items`
3. 在 Authorization header 中包含 token

**预期结果**: 201/200，返回创建的收藏项信息

### 场景 2: 无效 token

**步骤**:
1. 使用错误的 secret 签名 token
2. 发送请求

**预期结果**: 401，提示 "invalid or missing authorization token"

### 场景 3: 过期 token

**步骤**:
1. 生成 exp 为过去时间的 token
2. 发送请求

**预期结果**: 401，提示认证失败

### 场景 4: 缺少 Authorization header

**步骤**:
1. 不提供 Authorization header
2. 发送请求

**预期结果**: 401，提示缺少认证

## 安全建议

1. **Secret 管理**:
   - 不要在代码中硬编码 secret
   - 使用环境变量或配置管理系统存储 secret
   - 定期轮换 secret

2. **Token 有效期**:
   - 不要设置过长的 AccessExpire
   - 默认设置为 7200 秒（2 小时）
   - 对于高安全性需求，考虑更短的有效期

3. **HTTPS**:
   - 生产环境必须使用 HTTPS
   - 防止 token 在传输中被截获

4. **刷新机制**:
   - 考虑实现 refresh token 机制
   - 允许客户端刷新过期的 token

## 后续优化

- [ ] 实现 refresh token 机制（长期有效）
- [ ] 添加 token 黑名单（logout 时使用）
- [ ] 实现基于角色的访问控制 (RBAC)
- [ ] 添加率限制 (rate limiting)
- [ ] 记录认证失败日志
