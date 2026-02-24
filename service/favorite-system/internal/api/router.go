package api

import (
	"context"
	"net/http"
	"time"

	"favorite-system/internal/api/handler"
	"favorite-system/internal/app"
	"favorite-system/internal/pkg/httpx"

	"github.com/gin-gonic/gin"
)

func NewRouter(a *app.App) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.Use(gin.Recovery())

	// healthz
	r.GET("/healthz", func(c *gin.Context) {
		httpx.OK(c, gin.H{"ok": true, "ver": "v1-folders-routes"})
	})

	// readyz (db ping)
	r.GET("/readyz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		if err := a.DB.Pool.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"ready": false})
			return
		}
		httpx.OK(c, gin.H{"ready": true})
	})

	v1 := r.Group("/api/v1")
	{
		// ping
		v1.GET("/ping", func(c *gin.Context) {
			httpx.OK(c, gin.H{"pong": true})
		})

		// folders
		fh := handler.NewFolderHandler(a)
		v1.POST("/folders", fh.Create)
		v1.GET("/folders", fh.ListByUser)
		v1.GET("/folders/:id", fh.GetByID)
		v1.DELETE("/folders/:id", fh.SoftDelete)
	}

	return r
}
