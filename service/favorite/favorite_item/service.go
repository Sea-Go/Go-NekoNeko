package favorite_item

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zeromicro/go-zero/core/logx"
)

// Service 收藏项业务逻辑
// 职责：
// 1. 将 HTTP 请求的数据转换为业务对象
// 2. 实施业务规则（权限、唯一性约束）
// 3. 调用 Repo 进行数据操作
// 4. 处理缓存失效
// 5. 返回业务错误而不是数据库错误
type Service struct {
	itemRepo   RepoInterface
	folderRepo interface{} // 这里假设 folder 也有 Repo，会在后续整合
	db         *pgxpool.Pool
	cache      CacheInterface
	cacheKey   *CacheKeyBuilder
}

// NewService 创建服务实例
func NewService(db *pgxpool.Pool, cache CacheInterface) *Service {
	if cache == nil {
		cache = &NopCache{} // 使用空缓存实现
	}
	return &Service{
		itemRepo: NewRepo(db),
		db:       db,
		cache:    cache,
		cacheKey: NewCacheKeyBuilder(),
	}
}

// CreateItem 创建收藏项（添加收藏）
// 业务规则：
// 1. 收藏夹必须存在且属于当前用户
// 2. 用户不能重复收藏同一对象
func (s *Service) CreateItem(ctx context.Context, userID int64, folderID int64, objectType string, objectID int64) (*ItemInfo, error) {
	// 业务规则 1: 检查重复收藏
	exists, _, _, err := s.itemRepo.CheckExists(ctx, userID, objectType, objectID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrItemAlreadyExists
	}

	// 业务规则 2: 检查收藏夹是否存在且属于当前用户
	// （这里会调用 folderRepo，但现在先简化）
	// TODO: 添加收藏夹权限检查

	// 执行数据操作
	item := &FavoriteItem{
		UserID:     userID,
		FolderID:   folderID,
		ObjectType: objectType,
		ObjectID:   objectID,
		SortOrder:  0,
	}
	if err := s.itemRepo.Add(ctx, item); err != nil {
		return nil, err
	}

	// 缓存失效：清除该用户该收藏夹的所有缓存
	// 使用模式匹配删除所有相关缓存
	pattern := s.cacheKey.FavoriteFolderPattern(userID, folderID)
	if err := s.cache.DeletePattern(ctx, pattern); err != nil {
		// 缓存失效失败不应影响业务逻辑，只记录日志
		logx.Errorf("cache DeletePattern failed: %v", err)
	} else {
		logx.Infof("cache invalidated for pattern=%s", pattern)
	}

	// 返回响应数据
	return &ItemInfo{
		ID:         item.ID,
		FolderID:   item.FolderID,
		ObjectType: item.ObjectType,
		ObjectID:   item.ObjectID,
		CreatedAt:  item.CreatedAt,
	}, nil
}

// DeleteItem 删除收藏项
// 业务规则：
// 1. 只能删除自己的收藏
// 2. 如果对象已经不存在，仍然返回成功（幂等性）
func (s *Service) DeleteItem(ctx context.Context, userID int64, objectType string, objectID int64) error {
	// 验证所有权（可选，也可以由 DELETE 语句的 WHERE 条件保证）
	exists, _, folderID, err := s.itemRepo.CheckExists(ctx, userID, objectType, objectID)
	if err != nil {
		return err
	}
	if !exists {
		// 保证幂等性：即使不存在也返回成功
		return nil
	}

	// 执行删除
	if err := s.itemRepo.Delete(ctx, userID, objectType, objectID); err != nil {
		return err
	}

	// 缓存失效：清除该用户该收藏夹的所有缓存
	if folderID > 0 {
		pattern := s.cacheKey.FavoriteFolderPattern(userID, folderID)
		if err := s.cache.DeletePattern(ctx, pattern); err != nil {
			// 缓存失效失败不应影响业务逻辑，只记录日志
			logx.Errorf("cache DeletePattern failed: %v", err)
		} else {
			logx.Infof("cache invalidated for pattern=%s", pattern)
		}
	}

	return nil
}

// ListItems 获取收藏夹中的项目
// 业务规则：
// 1. 只能查看自己的收藏夹
// 2. 分页返回
// 3. 优先从缓存读取，缓存未命中则从数据库读取并写入缓存
func (s *Service) ListItems(ctx context.Context, userID int64, folderID int64, page, pageSize int) ([]*ItemInfo, int64, error) {
	// 业务规则 1: 检查收藏夹所有权
	// TODO: 添加收藏夹所有权检查

	// 分页计算
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100 // 防止一次查询太多数据
	}
	offset := (page - 1) * pageSize

	// 生成缓存键
	cacheKey := s.cacheKey.FavoriteListKey(userID, folderID, page)

	// 步骤 1: 尝试从缓存读取
	var cachedResult struct {
		Items []*ItemInfo `json:"items"`
		Total int64       `json:"total"`
	}

	if err := s.cache.Get(ctx, cacheKey, &cachedResult); err == nil && cachedResult.Items != nil {
		// 缓存命中，直接返回
		logx.Infof("cache hit key=%s", cacheKey)
		return cachedResult.Items, cachedResult.Total, nil
	}
	logx.Infof("cache miss key=%s", cacheKey)

	// 步骤 2: 缓存未命中，从数据库读取
	items, err := s.itemRepo.ListByFolder(ctx, folderID, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 获取总数
	total, err := s.itemRepo.CountByFolder(ctx, folderID)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var infos []*ItemInfo
	for _, item := range items {
		infos = append(infos, &ItemInfo{
			ID:         item.ID,
			FolderID:   item.FolderID,
			ObjectType: item.ObjectType,
			ObjectID:   item.ObjectID,
			CreatedAt:  item.CreatedAt,
		})
	}

	// 步骤 3: 写入缓存
	// 为了避免缓存错误影响业务逻辑，缓存写入错误不应该返回给客户端
	toCache := struct {
		Items []*ItemInfo `json:"items"`
		Total int64       `json:"total"`
	}{
		Items: infos,
		Total: total,
	}

	if err := s.cache.Set(ctx, cacheKey, toCache); err != nil {
		// 缓存写入失败不影响返回，只记录日志
		logx.Errorf("cache set failed key=%s err=%v", cacheKey, err)
	} else {
		logx.Infof("cache set key=%s", cacheKey)
	}

	return infos, total, nil
}

// GetItemCount 获取收藏夹中的项目数
func (s *Service) GetItemCount(ctx context.Context, folderID int64) (int64, error) {
	return s.itemRepo.CountByFolder(ctx, folderID)
}
