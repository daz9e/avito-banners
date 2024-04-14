package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	Database *gorm.DB
	Redis    *redis.Client
	Ctx      = context.Background()
)
