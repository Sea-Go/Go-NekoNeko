# 🚀 快速参考 - 收藏系统

> 最小化步骤启动并测试收藏系统

---

## ⚡ 5 分钟快速启动

### 1️⃣ 初始化数据库 (1 分钟)

```powershell
cd doc\scripts
go run init_db.go
```

**输出**: `✨ 数据库初始化完成！`

---

### 2️⃣ 启动 API 服务 (立即)

```powershell
# 回到项目根目录
cd ..\..

# 编译（如果还没编译）
go build -o ./api_favorite ./api/favorite.go

# 运行服务
.\api_favorite
```

**输出**: 
```
[INFO] starting server on 0.0.0.0:8888
```

✅ 服务已启动，保持此窗口打开

---

### 3️⃣ 生成测试 Token (1 分钟)

```powershell
# 新打开一个 PowerShell 窗口
cd api\tools
go run jwt_generator.go
```

**输出**:
```
JWT Token (用于 Authorization header 中):
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

📌 **复制完整 token（包括 "Bearer " 前缀）**

---

### 4️⃣ 运行测试 (2 分钟)

```powershell
# 新打开第三个 PowerShell 窗口
cd doc\scripts

# 粘贴你复制的 token
$token = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 运行测试
.\integration_test_simple.ps1 -Token $token
```

**输出**:
```
================================================
✨ 所有测试通过！API 运行正常。
================================================
```

✅ **完成！** 收藏系统已成功启动并验证

---

## 📋 常用命令

### 编译

```bash
# 编译 favorite 服务
go build -o ./api_favorite ./api/favorite.go

# 编译 usercenter 服务（可选）
go build -o ./api_usercenter ./api/usercenter.go
```

### 测试

```bash
# 快速测试（需要预先生成 token）
cd doc\scripts
.\integration_test_simple.ps1 -Token "Bearer <your-token>"

# 详细测试指南
# 见 INTEGRATION_TEST_GUIDE.md
```

### 数据库

```bash
# 初始化数据库
cd doc\scripts
go run init_db.go

# 重置数据库（删除所有表）
psql -h 127.0.0.1 -U postgres -d favorite_db
# 在 psql 中执行:
# DROP TABLE IF EXISTS favorite_item CASCADE;
# DROP TABLE IF EXISTS favorite_folder CASCADE;
# DROP TABLE IF EXISTS auth_user CASCADE;
```

---

## 🔐 Token 说明

### 格式
```
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.payload.signature
│      └─────────── 实际 JWT token ──────────────────────┘
```

### 包含信息
- **user_id**: 1 (测试用户)
- **exp**: 24 小时后过期
- **iat**: 发行时间

### 有效期
- 默认: 24 小时
- 过期后需重新生成

---

## 🧪 关键 API 端点

### 创建收藏项
```bash
curl -X POST http://localhost:8888/favorite/v1/items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"folder_id":1,"object_type":"article","object_id":"12345"}'
```

### 列表收藏项
```bash
curl -X GET "http://localhost:8888/favorite/v1/items?folder_id=1&page=1&page_size=10" \
  -H "Authorization: Bearer <token>"
```

### 删除收藏项
```bash
curl -X DELETE http://localhost:8888/favorite/v1/items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"object_type":"article","object_id":"12345"}'
```

---

## ❌ 常见错误快速解决

### 错误: `Error connecting to database`

**解决**: 启动 PostgreSQL
```bash
docker-compose up -d
```

### 错误: `401 Unauthorized`

**解决**: 检查 token 格式
```bash
# ❌ 错误: Bearer<token> (无空格)
# ❌ 错误: Bearer token (缺少完整 token)

# ✅ 正确: Bearer eyJhbGciOiJIUzI1NiIsIn...
```

### 错误: `409 Conflict - 该对象已被收藏`

**解决**: 使用不同的 object_id 或先删除再创建
```bash
# 删除现有收藏
curl -X DELETE http://localhost:8888/favorite/v1/items \
  -H "Authorization: Bearer <token>" \
  -d '{"object_type":"article","object_id":"12345"}'

# 重新创建
curl -X POST http://localhost:8888/favorite/v1/items \
  -H "Authorization: Bearer <token>" \
  -d '{"folder_id":1,"object_type":"article","object_id":"12345"}'
```

---

## 📚 更多文档

| 文档 | 用途 |
|------|------|
| [IMPLEMENTATION_COMPLETE.md](IMPLEMENTATION_COMPLETE.md) | 完整实现总结 |
| [INTEGRATION_TEST_GUIDE.md](INTEGRATION_TEST_GUIDE.md) | 详细的测试指南 |
| [README_JWT.md](README_JWT.md) | JWT 认证指南 |
| [doc/jwt_authentication.md](doc/jwt_authentication.md) | JWT 工作原理 |

---

## ✅ 验证清单

在报告任何问题前，请检查：

- [ ] PostgreSQL 已启动 (`docker-compose ps`)
- [ ] 数据库已初始化 (`go run init_db.go` 成功)
- [ ] API 服务已启动 (`.\api_favorite` 无错误)
- [ ] Token 已生成 (`go run jwt_generator.go`)
- [ ] Token 格式正确 (以 "Bearer " 开头)
- [ ] favorite.yaml 中的 AccessSecret 是 "favorite-secret-key"

---

## 🎯 下一步

### 短期 (1-2 周)

- [ ] 添加 Redis 缓存
- [ ] 编写单元测试
- [ ] 性能基准测试

### 中期 (1-2 月)

- [ ] API 文档 (Swagger/OpenAPI)
- [ ] Batch 操作界面
- [ ] 高级查询功能

### 长期 (2-3 月)

- [ ] 微服务分解
- [ ] 消息队列集成
- [ ] 分布式缓存

---

## 💡 提示

💾 **定期备份** PostgreSQL 数据库
```bash
docker-compose exec postgres pg_dump -U postgres favorite_db > backup.sql
```

🔄 **清理日志** (可选)
```bash
rm -rf logs/*
```

🧹 **完全重置** (如果需要)
```bash
docker-compose down -v  # 删除所有数据
docker-compose up -d     # 重新启动
go run doc/scripts/init_db.go  # 重新初始化
```

---

## 📞 获取帮助

1. **查看测试指南**: [INTEGRATION_TEST_GUIDE.md](INTEGRATION_TEST_GUIDE.md)
2. **检查错误日志**: 查看 API 服务的控制台输出
3. **数据库日志**: `docker-compose logs postgres`
4. **代码文档**: 各文件中的注释说明

---

**最后更新**: 2026-02-09  
**状态**: ✅ 完整实现，生产就绪
