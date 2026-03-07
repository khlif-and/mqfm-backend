package redis

import (
	"context"

	goredis "github.com/redis/go-redis/v9"

	"mqfm-backend/internal/infrastructure/config"
	"mqfm-backend/internal/shared/logger"
)

func NewRedisClient(cfg *config.Config) *goredis.Client {
	client := goredis.NewClient(&goredis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		logger.Warn("redis connection failed, caching disabled: " + err.Error())
		return nil
	}

	logger.Info("redis connected successfully")
	return client
}
