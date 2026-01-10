package postgres

import (
	"fmt"
	"sea-try-go/service/favorite/rpc/internal/config"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type FavoriteRepo struct {
	Db *gorm.DB
}

func NewFavoriteService(c config.Config) *FavoriteRepo {
	db, err := InitDB(c)
	if err != nil {
		logx.Errorf("init db error:%v", err)
		panic(err)
	}
	return &FavoriteRepo{
		Db: db,
	}
}
func InitDB(c config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.Dbname,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //禁用复数表名
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// 获取底层 *sql.DB 以设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 自动迁移模型
	err = db.AutoMigrate(
		&Favorite{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate models: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	return db, nil
}
func (favorite *FavoriteRepo) Delete(userId, ArticleId uint64) error {
	err := favorite.Db.Where("user_id = ? AND article_id = ?", userId, ArticleId).Delete(&Favorite{}).Error
	return err
}
func (favorite *FavoriteRepo) Insert(userId, ArticleId uint64) error {
	fa := &Favorite{
		UserID:    userId,
		ArticleID: ArticleId,
	}
	err := favorite.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(fa).Error
	return err
}

func (favorite *FavoriteRepo) GetArtocleIdListByUserId(userId uint64) (*[]Favorite, error) {
	var res *[]Favorite
	err := favorite.Db.Where("user_id = ?", userId).Find(res).Error
	return res, err
}
