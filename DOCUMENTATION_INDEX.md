# 📚 文档索引

收藏系统完整文档导航

---

## 🚀 快速开始

### 首次使用？从这里开始

1. **[QUICK_START.md](QUICK_START.md)** ⭐ **5 分钟快速上手**
   - 最小化步骤启动服务
   - 关键命令速查
   - 常见错误解决

2. **[IMPLEMENTATION_COMPLETE.md](IMPLEMENTATION_COMPLETE.md)** - 实现总结
   - 项目完成状态
   - 架构说明
   - 功能清单

---

## 🔐 JWT 认证

### 需要理解 JWT 如何工作？

1. **[README_JWT.md](README_JWT.md)** ⭐ **JWT 快速参考**
   - 使用指南
   - 常见场景
   - 故障排除

2. **[doc/jwt_authentication.md](doc/jwt_authentication.md)** - 详细文档
   - 工作原理详解
   - 配置说明
   - 安全建议

3. **[doc/jwt_implementation_summary.md](doc/jwt_implementation_summary.md)** - 实现细节
   - 代码实现
   - 集成方式
   - 已知限制

---

## 🧪 测试与验证

### 想要运行测试？

1. **[INTEGRATION_TEST_GUIDE.md](INTEGRATION_TEST_GUIDE.md)** ⭐ **完整测试指南**
   - 测试环境准备
   - 数据库初始化方法
   - 所有测试场景说明
   - 常见问题解答

2. **[doc/scripts/integration_test_simple.ps1](doc/scripts/integration_test_simple.ps1)** - 简化版测试脚本
   - 快速验证 API
   - 自动化测试
   - 详细的输出报告

3. **[doc/scripts/init_db.go](doc/scripts/init_db.go)** - 数据库初始化脚本
   - Go 版本初始化工具
   - 自动建表和索引

4. **[doc/scripts/init_db.ps1](doc/scripts/init_db.ps1)** - PowerShell 初始化脚本
   - PowerShell 版本初始化工具

---

## 📊 项目信息

### 项目结构和设计

| 文档 | 内容 | 位置 |
|------|------|------|
| **架构** | 完整的系统架构图和分层设计 | [IMPLEMENTATION_COMPLETE.md](IMPLEMENTATION_COMPLETE.md#-完整架构图) |
| **文件结构** | 项目文件组织说明 | [IMPLEMENTATION_COMPLETE.md](IMPLEMENTATION_COMPLETE.md#-项目文件结构) |
| **API 文档** | HTTP 端点说明 | [INTEGRATION_TEST_GUIDE.md](INTEGRATION_TEST_GUIDE.md#-http-状态码映射) 和 [IMPLEMENTATION_COMPLETE.md](IMPLEMENTATION_COMPLETE.md#-api-文档) |
| **数据库架构** | 表设计和约束 | [doc/sql/](doc/sql/) 中的 SQL 脚本 |

---

## 🔧 工具和脚本

### 所有可用的脚本工具

| 脚本 | 功能 | 位置 | 用途 |
|------|------|------|------|
| **jwt_generator.go** | 生成测试 JWT token | [api/tools/jwt_generator.go](api/tools/jwt_generator.go) | 测试用 token 生成 |
| **init_db.go** | 初始化数据库 (Go) | [doc/scripts/init_db.go](doc/scripts/init_db.go) | 创建表结构 |
| **init_db.ps1** | 初始化数据库 (PowerShell) | [doc/scripts/init_db.ps1](doc/scripts/init_db.ps1) | Windows 脚本版本 |
| **integration_test_simple.ps1** | 简化版测试脚本 | [doc/scripts/integration_test_simple.ps1](doc/scripts/integration_test_simple.ps1) | 快速验证 API |
| **integration_test.ps1** | 完整测试脚本 | [doc/scripts/integration_test.ps1](doc/scripts/integration_test.ps1) | 自动化测试 |

---

## 📝 核心源代码文档

### Handler 层 (HTTP 端点)

| 文件 | 说明 |
|------|------|
| [api/internal/handler/favorite/createfavoritehandler.go](api/internal/handler/favorite/createfavoritehandler.go) | POST /items - 创建收藏 |
| [api/internal/handler/favorite/deletefavoritehandler.go](api/internal/handler/favorite/deletefavoritehandler.go) | DELETE /items - 删除收藏 |
| [api/internal/handler/favorite/listfavoritehandler.go](api/internal/handler/favorite/listfavoritehandler.go) | GET /items - 列表查询 |

### Logic 层 (业务逻辑)

| 文件 | 说明 |
|------|------|
| [api/internal/logic/favorite/createfavoritelogic.go](api/internal/logic/favorite/createfavoritelogic.go) | 创建逻辑处理 |
| [api/internal/logic/favorite/deletefavoritelogic.go](api/internal/logic/favorite/deletefavoritelogic.go) | 删除逻辑处理 |
| [api/internal/logic/favorite/listfavoritelogic.go](api/internal/logic/favorite/listfavoritelogic.go) | 列表逻辑处理 |

### Service 层 (服务实现)

| 文件 | 说明 |
|------|------|
| [service/favorite/favorite_item/service.go](service/favorite/favorite_item/service.go) | 核心业务逻辑 |
| [service/favorite/favorite_item/repo.go](service/favorite/favorite_item/repo.go) | 数据库操作 |
| [service/favorite/favorite_item/model.go](service/favorite/favorite_item/model.go) | 数据模型定义 |
| [service/favorite/favorite_item/error.go](service/favorite/favorite_item/error.go) | 错误定义 |

### 工具类

| 文件 | 说明 |
|------|------|
| [api/internal/utils/jwt.go](api/internal/utils/jwt.go) | JWT 验证工具 |
| [api/internal/utils/error_mapper.go](api/internal/utils/error_mapper.go) | 错误映射工具 |

---

## 🎯 按任务查找文档

### 我想...

#### 快速启动项目
→ [QUICK_START.md](QUICK_START.md)

#### 理解 JWT 认证
→ [README_JWT.md](README_JWT.md) 然后 [doc/jwt_authentication.md](doc/jwt_authentication.md)

#### 运行测试
→ [INTEGRATION_TEST_GUIDE.md](INTEGRATION_TEST_GUIDE.md)

#### 初始化数据库
→ [INTEGRATION_TEST_GUIDE.md#数据库初始化](INTEGRATION_TEST_GUIDE.md#数据库初始化)

#### 理解项目架构
→ [IMPLEMENTATION_COMPLETE.md#-完整架构图](IMPLEMENTATION_COMPLETE.md#-完整架构图)

#### 查看 API 端点
→ [INTEGRATION_TEST_GUIDE.md#测试场景说明](INTEGRATION_TEST_GUIDE.md#测试场景说明) 或 [IMPLEMENTATION_COMPLETE.md#-api-文档](IMPLEMENTATION_COMPLETE.md#-api-文档)

#### 解决遇到的问题
→ [QUICK_START.md#-常见错误快速解决](QUICK_START.md#-常见错误快速解决) 或 [INTEGRATION_TEST_GUIDE.md#常见问题](INTEGRATION_TEST_GUIDE.md#常见问题)

#### 扩展或修改代码
→ [IMPLEMENTATION_COMPLETE.md](IMPLEMENTATION_COMPLETE.md) 了解架构，然后查看相应的源代码文件

---

## 📋 SQL 脚本

### 数据库表定义

| 表 | 用途 | 文件 |
|---|------|------|
| `auth_user` | 用户信息 | [doc/sql/favorite_item.sql](doc/sql/favorite_item.sql) (创建时) |
| `favorite_folder` | 收藏夹 | [doc/sql/favorite_folder.sql](doc/sql/favorite_folder.sql) |
| `favorite_item` | 收藏项目 | [doc/sql/favorite_item.sql](doc/sql/favorite_item.sql) |

---

## 🔗 文档关系图

```
QUICK_START.md (快速开始)
    ├── IMPLEMENTATION_COMPLETE.md (项目总结)
    │   ├── 架构图 → INTEGRATION_TEST_GUIDE.md
    │   └── API 文档 → 各 Handler 源代码
    │
    ├── README_JWT.md (JWT 使用)
    │   └── 详细文档 → doc/jwt_authentication.md
    │
    └── INTEGRATION_TEST_GUIDE.md (测试指南)
        ├── 初始化 → doc/scripts/init_db.go
        ├── 测试 → doc/scripts/integration_test_simple.ps1
        └── 常见问题 → FAQ 答案
```

---

## 📞 快速参考

### 关键路径

| 需求 | 立即查看 |
|------|---------|
| 5分钟启动 | [QUICK_START.md](QUICK_START.md) |
| JWT 问题 | [README_JWT.md](README_JWT.md) |
| 测试失败 | [INTEGRATION_TEST_GUIDE.md#常见问题](INTEGRATION_TEST_GUIDE.md#常见问题) |
| 数据库错误 | [INTEGRATION_TEST_GUIDE.md#常见问题](INTEGRATION_TEST_GUIDE.md#常见问题) |
| 代码实现 | [IMPLEMENTATION_COMPLETE.md](IMPLEMENTATION_COMPLETE.md) |

---

## ✅ 文档检查清单

所有必必要文档都已准备：

- [x] **QUICK_START.md** - 快速上手指南
- [x] **IMPLEMENTATION_COMPLETE.md** - 完整实现总结
- [x] **INTEGRATION_TEST_GUIDE.md** - 测试指南
- [x] **README_JWT.md** - JWT 快速参考
- [x] **doc/jwt_authentication.md** - JWT 详细文档
- [x] **doc/jwt_implementation_summary.md** - 实现细节
- [x] **doc/scripts/init_db.go** - 数据库初始化工具
- [x] **doc/scripts/init_db.ps1** - PowerShell 初始化工具
- [x] **doc/scripts/integration_test_simple.ps1** - 测试脚本
- [x] **本文件 (DOCUMENTATION_INDEX.md)** - 文档索引

---

**最后更新**: 2026-02-09  
**版本**: 1.0.0  
**状态**: ✅ 完整，生产就绪
