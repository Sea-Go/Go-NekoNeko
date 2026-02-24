package app

import (
	"favorite-system/internal/config"
	db "favorite-system/internal/repo/db"
	"favorite-system/internal/repo/folder"
	"favorite-system/internal/repo/pg"
)

type App struct {
	Cfg     *config.Config
	DB      *pg.DB
	Queries *db.Queries

	FolderRepo *folder.Repo
}
