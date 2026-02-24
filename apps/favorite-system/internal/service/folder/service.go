package folder

import (
	"context"
	"strings"

	"favorite-system/internal/repo/db"
	folderrepo "favorite-system/internal/repo/folder"
)

type Service struct {
	repo *folderrepo.Repo
}

func New(repo *folderrepo.Repo) *Service {
	return &Service{
		repo: repo,
	}
}

//
// ========= Create =========
//

type CreateInput struct {
	UserID   int64  `json:"user_id"`
	Name     string `json:"name"`
	IsPublic bool   `json:"is_public"`
}

func (s *Service) Create(ctx context.Context, in CreateInput) (db.Folder, error) {
	in.Name = strings.TrimSpace(in.Name)

	if in.UserID <= 0 {
		return db.Folder{}, ErrInvalidUserID
	}

	if in.Name == "" {
		return db.Folder{}, ErrInvalidFolderName
	}

	return s.repo.Create(ctx, folderrepo.CreateInput{
		UserID:   in.UserID,
		Name:     in.Name,
		IsPublic: in.IsPublic,
	})
}

//
// ========= ListByUser =========
//

func (s *Service) ListByUser(ctx context.Context, userID int64) ([]db.Folder, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}

	return s.repo.ListByUser(ctx, userID)
}

//
// ========= GetByID =========
//

func (s *Service) GetByID(ctx context.Context, id int64) (db.Folder, error) {
	if id <= 0 {
		return db.Folder{}, ErrInvalidFolderID
	}

	return s.repo.GetByID(ctx, id)
}

//
// ========= SoftDelete =========
//

func (s *Service) SoftDelete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidFolderID
	}

	return s.repo.SoftDelete(ctx, id)
}
