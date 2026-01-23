package svc

import (
	"sea-try-go/service/hotspot/rpc/internal/config"
	"sea-try-go/service/hotspot/rpc/model"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	// LikesModel    model.LikesModel
	// CommentsModel model.CommentsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// DB := model.InitDB(c.DB.DataSource)
	return &ServiceContext{
		Config: c,
		DB:     model.InitDB(c.DB.DataSource),
	}
}
