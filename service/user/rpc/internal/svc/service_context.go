package svc

import (
	"sea-try-go/service/user/rpc/internal/config"
	"sea-try-go/service/user/rpc/internal/model"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Config    config.Config
	UserModel *model.UserModel
}

func NewServiceContext(cfg config.Config, db *gorm.DB) *ServiceContext {
	return &ServiceContext{
		Config:    cfg,
		UserModel: model.NewUserModel(db),
	}
}
