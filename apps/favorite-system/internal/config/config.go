package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	AppEnv string
	Port   string
	DBDsn  string
}

func Load() *Config {
	_ = godotenv.Load()

	viper.AutomaticEnv()

	cfg := &Config{
		AppEnv: viper.GetString("APP_ENV"),
		Port:   viper.GetString("PORT"),
		DBDsn:  viper.GetString("DB_DSN"),
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	log.Println("config loaded")
	return cfg
}
