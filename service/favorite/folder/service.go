package folder

import "context"

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

// CreateFolder 创建收藏夹
func (s *Service) CreateFolder(
	ctx context.Context,
	userID int64,
	name string,
	isPublic bool,
) (*Folder, error) {

	if name == "" {
		return nil, ErrFolderNameEmpty
	}
	if len(name) > 100 {
		return nil, ErrFolderNameTooLong
	}

	exist, err := s.repo.ExistsByUserAndName(ctx, userID, name)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, ErrFolderNameExists
	}

	folder := &Folder{
		UserID:   userID,
		Name:     name,
		IsPublic: isPublic,
	}

	if err := s.repo.Create(ctx, folder); err != nil {
		return nil, err
	}

	return folder, nil
}

// ListByUser 列出用户的收藏夹
func (s *Service) ListByUser(ctx context.Context, userID int64) ([]*Folder, error) {
	return s.repo.ListByUser(ctx, userID)
}
