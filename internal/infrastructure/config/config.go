package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	DBUser          string
	DBPassword      string
	DBHost          string
	DBPort          string
	DBName          string
	JWTSecret       string
	YouTubeAPIKey   string
	GoogleClientID  string
	RedisAddr       string
	RedisPassword   string
	RedisDB         int
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:           getEnv("PORT", "8080"),
		DBUser:         getEnv("DB_USER", ""),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBHost:         getEnv("DB_HOST", "127.0.0.1"),
		DBPort:         getEnv("DB_PORT", "3306"),
		DBName:         getEnv("DB_NAME", ""),
		JWTSecret:      getEnv("JWT_SECRET", "mqfm_secret_key_default"),
		YouTubeAPIKey:  getEnv("YOUTUBE_API_KEY", ""),
		GoogleClientID: getEnv("GOOGLE_CLIENT_ID", ""),
		RedisAddr:      getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
