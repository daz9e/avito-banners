package tools

import (
	"avito-banners/config"
	"github.com/go-redis/redis/v8"
)

func SetupRedis() {
	config.Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
