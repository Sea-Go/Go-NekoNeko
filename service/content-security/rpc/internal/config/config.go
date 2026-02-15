package config

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	logx.LogConf
	AdDetection struct {
		ApiEndpoint string  `json:"apiEndpoint"`
		ApiKey      string  `json:"apiKey"`
		Threshold   float64 `json:"threshold"`
		Timeout     int     `json:"timeout"`
	}
	HtmlSanitization struct {
		AllowedTags []string `json:"allowedTags"`
	}
	Cache cache.CacheConf
}
