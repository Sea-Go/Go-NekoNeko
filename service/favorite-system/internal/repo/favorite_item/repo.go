package favorite_item

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

func (r *Repo) Add(ctx context.Context, arg db.AddFavoriteItemParams) (db.FavoriteItem, error) {
	return r.q.AddFavoriteItem(ctx, arg)
}

func (r *Repo) Remove(ctx context.Context, arg db.SoftDeleteFavoriteItemParams) error {
	return r.q.SoftDeleteFavoriteItem(ctx, arg)
}

func (r *Repo) ListByFolder(ctx context.Context, folderID int64) ([]db.FavoriteItem, error) {
	return r.q.ListAllFavoriteItems(ctx, folderID)
}

func (r *Repo) Get(ctx context.Context, arg db.GetFavoriteItemParams) (db.FavoriteItem, error) {
	return r.q.GetFavoriteItem(ctx, arg)
}

func (r *Repo) SoftDeleteByFolder(ctx context.Context, folderID int64) error {
	return r.q.SoftDeleteFavoriteItemsByFolder(ctx, folderID)
}
