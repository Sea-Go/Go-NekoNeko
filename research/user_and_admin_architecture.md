# 用户与管理员服务架构调研

## 概述

本系统包含两个核心用户管理服务：`user`（普通用户）和 `admin`（管理员）。两者均遵循 Go-Zero 微服务架构，采用 API Gateway + RPC Service 的分层模式。API 层负责 HTTP 协议处理、JWT 认证和请求/响应转换；RPC 层负责核心业务逻辑和数据库交互。两个服务共享部分通用模块（如密码加密、雪花ID生成、错误码、日志），但拥有各自独立的数据库表（`users` 和 `admins`）和 gRPC 接口定义。管理员服务 (`admin`) 具备管理普通用户 (`user`) 的能力，例如查询、封禁、删除等。

### 核心组件
- **User API/RPC**: 位于 `service/user/user/`，处理普通用户的注册、登录、信息维护。
- **Admin API/RPC**: 位于 `service/user/admin/`，处理管理员的创建、登录，并提供对普通用户的管理功能。
- **通用模块**: 位于 `service/common/` 和 `service/user/common/`，提供 JWT、密码加密、错误码、日志、雪花ID等共享功能。
- **数据存储**: 使用 PostgreSQL，通过 GORM 进行 ORM 操作。`user` 服务操作 `users` 表，`admin` 服务操作 `admins` 和 `users` 表。

## User 服务 (普通用户)

### API 层 (`service/user/user/api`)
- **入口**: `usercenter.go` 初始化服务，加载 `etc/usercenter.yaml` 配置。
- **配置**: `config/config.go` 定义了 REST 配置、JWT 配置 (`UserAuth`) 和指向 `user.rpc` 的 RPC 客户端配置 (`UserRpc`)。
- **路由**: `handler/routes.go` 定义了两组路由：
  - 无需认证: `/user/register`, `/user/login`。
  - 需 JWT 认证: `/user/get`, `/user/update`, `/user/delete`, `/user/logout`。
- **Handler**: 位于 `handler/user/`，负责解析 HTTP 请求，调用 Logic 层，并将结果包装成统一的 `response.Response` 格式返回。
- **Logic**: 位于 `internal/logic/user/`，是业务逻辑的核心。它从 JWT 的上下文中提取 `userId`，并调用 RPC 服务。例如，在 `getuserlogic.go` 中，它会验证 `ctx` 中的 `userId` 并将其作为查询 UID 调用 RPC。
- **类型定义**: `types/types.go` 定义了 API 层使用的请求/响应结构体。
- **服务上下文**: `svc/servicecontext.go` 将配置和 RPC 客户端 (`userservice.UserService`) 组装成 `ServiceContext`，供 Handler 和 Logic 使用。

### RPC 层 (`service/user/user/rpc`)
- **入口**: `user.go` 初始化 RPC 服务，加载 `etc/user.yaml` 配置，并初始化数据库连接。
- **配置**: `internal/config/config.go` 定义了 RPC 服务器配置和 PostgreSQL 数据库连接信息。
- **服务上下文**: `internal/svc/servicecontext.go` 初始化 GORM 数据库连接，并创建 `UserModel` 实例。
- **模型 (Model)**: 
  - `internal/model/data_model.go` 定义了 `User` 结构体，对应 `users` 表。
  - `internal/model/user.go` 提供了 `UserModel`，封装了对 `users` 表的 CRUD 操作，如 `FindOneByUserName`, `Insert`, `UpdateUserById` 等。
- **逻辑 (Logic)**: 位于 `internal/logic/`，实现了具体的业务逻辑。例如：
  - `registerlogic.go`: 处理注册，检查用户名唯一性，生成雪花ID，加密密码，存入数据库。
  - `loginlogic.go`: 处理登录，验证用户名密码，并根据用户状态 (`Status`) 返回不同结果（0:成功, 1:用户/密码错, 2:用户被封禁）。
  - `updateuserlogic.go`: 处理用户信息更新，支持更新用户名、密码、邮箱和额外信息，并检查用户名冲突。
- **gRPC 定义**: `proto/user.proto` 定义了 `UserService` 及其所有 RPC 方法和消息类型。
- **自动生成代码**: `pb/` 目录包含由 `protoc` 生成的 gRPC 服务端和客户端代码。`userservice/` 目录包含由 `goctl` 生成的 RPC 客户端封装。

## Admin 服务 (管理员)

### API 层 (`service/user/admin/api`)
- **入口**: `admincenter.go` 初始化服务，加载 `etc/admincenter.yaml` 配置。
- **配置**: `config/config.go` 定义了 REST 配置、JWT 配置 (`AdminAuth`) 和指向 `admin.rpc` 的 RPC 客户端配置 (`AdminRpc`)。
- **路由**: `handler/routes.go` 定义了两组路由：
  - 无需认证: `/admin/create`, `/admin/login`。
  - 需 JWT 认证: 所有其他管理接口，如 `/admin/getuser`, `/admin/banuser` 等。
- **Handler/Logic/Types/SVC**: 结构与 User API 层类似，位于对应的 `internal/` 子目录中。Logic 层会调用 `admin.rpc` 来执行操作。

### RPC 层 (`service/user/admin/rpc`)
- **入口**: `admin.go` 初始化 RPC 服务，加载 `etc/admin.yaml` 配置，并初始化数据库连接。
- **配置**: `internal/config/config.go` 定义了 RPC 服务器配置、PostgreSQL DSN (`DataSource`) 和系统默认密码 (`System.DefaultPassword`)。
- **服务上下文**: `internal/svc/servicecontext.go` 初始化 GORM 数据库连接，并创建 `AdminModel` 实例。
- **模型 (Model)**:
  - `internal/model/data_model.go` 同时定义了 `Admin` 和 `User` 两个结构体，分别对应 `admins` 和 `users` 表。
  - `internal/model/admin.go` 提供了 `AdminModel`，它同时封装了对 `admins` 表和 `users` 表的操作。这使得管理员服务可以直接管理普通用户数据。
- **逻辑 (Logic)**: 位于 `internal/logic/`，实现了管理员专属逻辑：
  - `createadminlogic.go`: 创建新管理员。
  - `banuserlogic.go` / `unbanuserlogic.go`: 通过更新 `users` 表的 `status` 字段来封禁/解封用户。
  - `deleteuserlogic.go`: 直接从 `users` 表中删除用户记录。
  - `resetuserpasswordlogic.go`: 重置用户密码为系统默认密码。
  - `updateuserlogic.go`: 更新普通用户信息。
- **gRPC 定义**: `proto/admin.proto` 定义了 `AdminService`，它不仅包含管理员自身的操作（`GetSelf`, `CreateAdmin`），还包含了大量对普通用户的管理操作（`GetUser`, `BanUser`, `DeleteUser` 等）。
- **自动生成代码**: `pb/` 和 `adminservice/` 目录包含生成的 gRPC 代码。

## 共享模块

### 密码与 JWT (`service/user/common/`)
- **密码加密**: `cryptx/password.go` 使用 `bcrypt` 提供密码哈希和校验功能。
- **JWT**: `jwt/jwt.go` 封装了 JWT token 的生成逻辑。

### 错误处理 (`service/user/common/errmsg`)
- **错误码**: `errmsg.go` 定义了一套全局的错误码（如 `ErrorUserExist=1001`）和对应的错误信息。API 层的 Logic 在调用 RPC 后，会根据 gRPC 状态码（`codes.AlreadyExists`, `codes.NotFound` 等）映射到这套业务错误码，并通过 `errmsg.GetErrMsg` 获取描述返回给前端。

### 日志 (`service/common/logger`)
- **结构化日志**: `logger.go` 提供了 `LogBusinessErr`, `LogInfo` 等方法，用于记录带有丰富上下文（如 `service_name`, `trace_id`, `event_id`, 调用链路）的结构化日志。它依赖于 `sea-try-go/service/common/snowflake` 生成唯一的 `event_id`。

### 雪花ID (`service/common/snowflake`)
- **分布式ID**: `snowflake.go` 基于 `bwmarrin/snowflake` 库实现，使用服务名和进程ID（PID）生成唯一的节点ID，以保证在单机多实例环境下ID的唯一性。

## 数据流示例：管理员封禁用户

1.  **HTTP 请求**: `POST /admincenter/v1/admin/banuser` 携带 JWT token 和 `{"uid": 123}`。
2.  **Admin API**:
    - `handler/admin/banuserhandler.go` 解析请求。
    - `logic/admin/banuserlogic.go` 从 JWT 中获取管理员自己的 `userId`（此处未使用），并将 `uid=123` 作为 `BanUserReq` 通过 `AdminRpc` 发送给 Admin RPC 服务。
3.  **Admin RPC**:
    - `server/adminserviceserver.go` 接收请求并调用 `logic/banuserlogic.go`。
    - `logic/banuserlogic.go` 调用 `AdminModel.UpdateUserStatusByUid(123, 1)`。
    - `model/admin.go` 执行 SQL `UPDATE users SET status = 1 WHERE uid = 123`。
4.  **响应**: 操作成功后，结果逐层返回，最终 Admin API 返回 JSON 响应 `{"code": 200, "msg": "OK", "data": {"success": true}}`。
