package config_conn

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl string
	RedisUrl    string
	HMACKey     string
}

func GetConfig() Config {
	_ = godotenv.Load()
	cfg := Config{
		DatabaseUrl: os.Getenv("DATABASE_URL"),
		RedisUrl:    os.Getenv("REDIS_URL"),
		HMACKey:     os.Getenv("HMAC_KEY"),
	}
	return cfg
}
