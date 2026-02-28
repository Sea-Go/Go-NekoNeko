# 收藏系统设计文档

## 1. 需求概述
实现一个通用的收藏系统，支持用户创建收藏夹、收藏任意类型的对象（如文章、视频等）。

### 核心功能
1.  **收藏夹管理 (CRUD)**:
    *   创建、删除、更新收藏夹。
    *   查询用户的收藏夹列表。
    *   支持公開/私有设置 (`is_public`)。
2.  **收藏内容管理**:
    *   添加内容到收藏夹。
    *   从收藏夹移除内容。
    *   查询收藏夹内的内容列表。
3.  **约束**:
    *   同一用户对同一对象不能重複收藏（全局唯一约束 `user_id, object_type, object_id`）。
4.  **性能优化**:
    *   使用 Redis 缓存活跃用户的收藏列表。
    *   缓存自动过期（TTL），不活跃用户数据自动从 Redis 清除。

## 2. 接口设计 (API)

### 2.1 收藏夹 (Folder)
*   `POST /api/v1/favorite/folder/create`: 创建收藏夹
*   `POST /api/v1/favorite/folder/update`: 更新收藏夹 (名称, 公开性)
*   `POST /api/v1/favorite/folder/delete`: 删除收藏夹
*   `GET /api/v1/favorite/folder/list`: 获取用户收藏夹列表

### 2.2 收藏项 (Item)
*   `POST /api/v1/favorite/item/add`: 添加收藏
*   `POST /api/v1/favorite/item/remove`: 移除收藏
*   `GET /api/v1/favorite/item/list`: 获取收藏夹内容列表

## 3. 数据库设计

### 3.1 收藏夹表 (folders)
| 字段 | 类型 | 说明 |
| --- | --- | --- |
| id | bigint | 主键 |
| user_id | bigint | 用户ID |
| name | varchar | 收藏夹名称 |
| is_public | boolean | 是否公开 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |
| deleted_at | timestamp | 软删除时间 |

### 3.2 收藏项表 (favorite_items)
| 字段 | 类型 | 说明 |
| --- | --- | --- |
| id | bigint | 主键 |
| folder_id | bigint | 归属收藏夹ID |
| user_id | bigint | 用户ID (冗余，用于快速校验唯一性) |
| object_type | varchar | 对象类型 (article, video...) |
| object_id | bigint | 对象ID |
| title | varchar | 标题 (快照) |
| created_at | timestamp | 创建时间 |

**索引**:
*   `unique_index(user_id, object_type, object_id)`: 防止重复收藏。

## 4. 缓存策略 (Redis)
*   **Key**: `favorite:user:<user_id>:folders` (Hash/List)
*   **Key**: `favorite:folder:<folder_id>:items` (ZSet/List)
*   **TTL**: 设置过期时间（如 1 小时），访问时自动续期。
*   **一致性**: 写操作（增删改）时直接删除/更新缓存。
