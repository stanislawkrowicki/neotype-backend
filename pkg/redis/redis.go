package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"neotype-backend/pkg/config"
	"os"
)

const (
	envPassword = "REDIS_PASSWORD"
	configAddr  = "addr"
	configPort  = "port"
	redisDB     = 0 // default
)

func NewConnection() *redis.Client {
	addr, err := config.Get("redis", "addr")
	port, err := config.Get("redis", "port")
	if err != nil {
		panic("failed to get config for redis")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", addr, port),
		Password: os.Getenv(envPassword),
		DB:       redisDB,
	})

	return rdb
}
