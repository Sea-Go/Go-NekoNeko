package model

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 初始化 GORM 连接
func InitDB(dataSource string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dataSource), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // DEBUG
	})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	err = db.AutoMigrate(&Like{}, &Comment{})
	if err != nil {
		log.Fatalf("数据表迁移失败: %v\n", err)
	}
	return db
}
