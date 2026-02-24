package folder

import (
	"context"

	db "favorite-system/internal/repo/db"
)

type Repo struct {
	q *db.Queries
}

func New(q *db.Queries) *Repo {
	return &Repo{q: q}
}

func (r *Repo) Create(ctx context.Context, arg db.CreateFolderParams) (db.Folder, error) {
	return r.q.CreateFolder(ctx, arg)
}

func (r *Repo) GetByID(ctx context.Context, id int64) (db.Folder, error) {
	return r.q.GetFolderByID(ctx, id)
}

func (r *Repo) ListByUser(ctx context.Context, userID int64) ([]db.Folder, error) {
	return r.q.ListFoldersByUser(ctx, userID)
}

func (r *Repo) SoftDelete(ctx context.Context, id int64) error {
	return r.q.SoftDeleteFolder(ctx, id)
}
