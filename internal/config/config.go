package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"os"
)

var (
	Database *gorm.DB
	Redis    *redis.Client
	Ctx      = context.Background()
	Cfg      *Config
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSSLMode  string
	AdminToken string
	UserToken  string
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	config := Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "user"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "bannerservice"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	return config, nil
}

func GetUserToken() string {
	return getEnv("USER_TOKEN", "user_token")
}

func GetAdminToken() string {
	return getEnv("ADMIN_TOKEN", "admin_token")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
