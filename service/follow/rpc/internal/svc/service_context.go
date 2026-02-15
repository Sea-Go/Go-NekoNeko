package svc

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sea-try-go/service/follow/rpc/internal/config"
	"sea-try-go/service/follow/rpc/internal/model" // 引入你的 model 包
)

type ServiceContext struct {
	Config      config.Config
	Neo4jDriver neo4j.DriverWithContext
	FollowModel model.FollowModel // 新增这一行：注入 Model 层
}

func NewServiceContext(c config.Config) *ServiceContext {
	driver, err := neo4j.NewDriverWithContext(
		c.Neo4j.Uri,
		neo4j.BasicAuth(c.Neo4j.Username, c.Neo4j.Password, ""),
	)
	if err != nil {
		panic("Failed to connect to Neo4j: " + err.Error())
	}

	err = driver.VerifyConnectivity(context.Background())
	if err != nil {
		panic("Neo4j connectivity verification failed: " + err.Error())
	}

	return &ServiceContext{
		Config:      c,
		Neo4jDriver: driver,
		FollowModel: model.NewFollowModel(driver), // 初始化 Model
	}
}
