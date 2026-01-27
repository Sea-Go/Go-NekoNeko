package model

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConf struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Mode     string
}

func InitDB(conf DBConf) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
		conf.Host,
		conf.User,
		conf.Password,
		conf.DBName,
		conf.Port,
		conf.Mode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("数据库连接失败")
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalln("数据表迁移失败")
	}
	return db
}
