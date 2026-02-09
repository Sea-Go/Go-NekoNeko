// Code scaffolded by goctl. DO NOT EDIT.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	PgDsn    string
	UserAuth struct {
		AccessSecret string
		AccessExpire int64
	}
	// Redis 缓存配置
	Redis struct {
		Enabled  bool
		Addr     string
		Password string
		TTL      int
	}
}
