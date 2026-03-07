package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"mqfm-backend/internal/infrastructure/config"
	"mqfm-backend/internal/infrastructure/database"
	"mqfm-backend/internal/infrastructure/di"
	infraRedis "mqfm-backend/internal/infrastructure/redis"
	"mqfm-backend/internal/infrastructure/router"
	"mqfm-backend/internal/shared/logger"
)

func main() {
	logger.Init()
	defer logger.Sync()

	cfg := config.Load()

	if cfg.YouTubeAPIKey == "" {
		logger.Fatal("YOUTUBE_API_KEY is missing")
	}

	db := database.NewMySQL(cfg)
	redisClient := infraRedis.NewRedisClient(cfg)

	container := di.NewContainer(db, redisClient, cfg.YouTubeAPIKey)

	r := gin.New()
	r.Use(gin.Recovery())

	router.Setup(r, container.Handlers)

	mqfmChannelID := "UCwa0rj5KY6bWoVzJtgoiaDw"
	go func() {
		logger.Info("scheduler started: checking youtube live status")
		for {
			if err := container.LivestreamService.UpdateLiveStatus(mqfmChannelID); err != nil {
				logger.Error("scheduler: youtube status update failed", zap.Error(err))
			}
			time.Sleep(10 * time.Minute)
		}
	}()

	logger.Info("server started", zap.String("port", cfg.Port))
	if err := r.Run(":" + cfg.Port); err != nil {
		logger.Fatal("server failed to start", zap.Error(err))
	}
}