package svc

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"sea-try-go/service/follow/rpc/internal/config"
)

type ServiceContext struct {
	Config      config.Config
	Neo4jDriver neo4j.DriverWithContext // 注入 Neo4j 驱动池
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 Neo4j 连接
	driver, err := neo4j.NewDriverWithContext(
		c.Neo4j.Uri,
		neo4j.BasicAuth(c.Neo4j.Username, c.Neo4j.Password, ""),
	)
	if err != nil {
		// 竞赛思维：如果连不上数据库，直接 panic 终止进程 (Fail Fast)，不要带病运行
		panic("Failed to connect to Neo4j: " + err.Error())
	}

	// 验证连通性
	err = driver.VerifyConnectivity(context.Background())
	if err != nil {
		panic("Neo4j connectivity verification failed: " + err.Error())
	}

	return &ServiceContext{
		Config:      c,
		Neo4jDriver: driver,
	}
}
