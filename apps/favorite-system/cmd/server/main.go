package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"favorite-system/internal/api"
	"favorite-system/internal/app"
	"favorite-system/internal/config"
	"favorite-system/internal/pkg/logger"
	db "favorite-system/internal/repo/db"
	"favorite-system/internal/repo/folder"
	"favorite-system/internal/repo/pg"

	"go.uber.org/zap"
)

func main() {
	// 1) load config
	cfg := config.Load()

	// 2) init logger
	logger.Init(cfg.AppEnv)
	defer logger.Sync()

	// 3) handle shutdown signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 4) init postgres pool
	pgdb, err := pg.New(ctx, cfg.DBDsn)
	if err != nil {
		logger.Log.Fatal("db init failed", zap.Error(err))
	}
	defer pgdb.Close()

	// 5) sqlc queries
	queries := db.New(pgdb.Pool)

	// 6) build app container
	a := &app.App{
		Cfg:     cfg,
		DB:      pgdb,
		Queries: queries,
	}

	// 7) repos
	a.FolderRepo = folder.New(a.Queries)

	// 8) router
	r := api.NewRouter(a)

	// 9) run server
	srvErr := make(chan error, 1)
	go func() {
		logger.Log.Info("server listening", zap.String("port", cfg.Port))
		srvErr <- r.Run(":" + cfg.Port)
	}()

	// 10) wait for exit
	select {
	case err := <-srvErr:
		logger.Log.Fatal("server failed", zap.Error(err))
	case <-ctx.Done():
		logger.Log.Info("shutdown signal received")
		time.Sleep(200 * time.Millisecond)
	}
}
