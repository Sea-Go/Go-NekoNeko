// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"log"
	"sea-try-go/service/user/api/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(postgres.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		log.Fatalln("数据库连接失败:", err)
	}
	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
