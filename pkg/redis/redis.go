package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"neotype-backend/pkg/config"
	"os"
)

const (
	envProdRedisURL = "REDIS_URL" // used in production only
	envPassword     = "REDIS_PASSWORD"
	redisDB         = 0 // default
)

func NewConnection() *redis.Client {
	redisURL, exists := os.LookupEnv(envProdRedisURL)
	if exists { // Production mode
		parsed, err := redis.ParseURL(redisURL)
		if err != nil {
			panic("failed to parse redis URL")
		}
		return redis.NewClient(parsed)
	}

	// Local development mode
	addr, err := config.Get("redis", "addr")
	port, err := config.Get("redis", "port")
	if err != nil {
		panic("failed to get config for redis")
	}

	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", addr, port),
		Password: os.Getenv(envPassword),
		DB:       redisDB,
	})
}
