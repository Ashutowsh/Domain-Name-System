package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppName       string
	AppEnv        string
	AppPort       int
	Database      string
	RedisURL      string
	REDIS_URL     string
	RedisPassword string
	RateLimit     int
}

func LoadConfig() AppConfig {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found : ", err)
	}

	appPort, err := strconv.Atoi(getEnv("APP_PORT", "2053"))
	if err != nil {
		fmt.Println("Invalid, app port number : ", err)
	}

	rateLimit, err := strconv.Atoi(getEnv("RATE_LIMIT", "5"))
	if err != nil {
		fmt.Println("Invalid, rate limit : ", err)
	}

	config := AppConfig{
		AppName:       getEnv("APP_NAME", "DNS-Server"),
		AppEnv:        getEnv("APP_ENV", "development"),
		AppPort:       appPort,
		Database:      getEnv("DATABASE_URL", ""),
		REDIS_URL:     getEnv("REDIS_URL", ""),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RateLimit:     rateLimit,
	}

	return config
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
