package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func Init(env string) {
	if env == "local" {
		Log, _ = zap.NewDevelopment()
	} else {
		Log, _ = zap.NewProduction()
	}
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
