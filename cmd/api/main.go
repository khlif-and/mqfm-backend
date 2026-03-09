package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"mqfm-backend/internal/infrastructure/config"
	"mqfm-backend/internal/infrastructure/database"
	"mqfm-backend/internal/infrastructure/di"
	infraRedis "mqfm-backend/internal/infrastructure/redis"
	"mqfm-backend/internal/infrastructure/router"
	"mqfm-backend/internal/infrastructure/scheduler"
	"mqfm-backend/internal/shared/logger"
)

func main() {
	logger.Init()
	defer logger.Sync()

	cfg := config.Load()

	db := database.NewMySQL(cfg)
	redisClient := infraRedis.NewRedisClient(cfg)

	container := di.NewContainer(db, redisClient, cfg)

	cron := scheduler.NewScoreRecalculator(container.RecommendationService, 3)
	cron.Start()

	votingScheduler := scheduler.NewVotingScheduler(container.AudioVoteService)
	votingScheduler.Start()

	notifScheduler := scheduler.NewNotificationScheduler(container.NotificationService)
	notifScheduler.Start()

	r := gin.New()
	r.Use(gin.Recovery())

	router.Setup(r, container.Handlers)

	logger.Info("server started", zap.String("port", cfg.Port))
	if err := r.Run(":" + cfg.Port); err != nil {
		logger.Fatal("server failed to start", zap.Error(err))
	}
}