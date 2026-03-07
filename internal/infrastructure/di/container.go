package di

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"mqfm-backend/internal/adapter/cache"
	adminHandler "mqfm-backend/internal/adapter/handler/admin"
	userHandler "mqfm-backend/internal/adapter/handler/user"
	mysqlRepo "mqfm-backend/internal/adapter/repository/mysql"
	"mqfm-backend/internal/app/service"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/infrastructure/config"
	"mqfm-backend/internal/infrastructure/email"
	"mqfm-backend/internal/infrastructure/router"
)

type Container struct {
	Handlers              *router.Handlers
	RecommendationService port.RecommendationService
	Cache                 port.CacheRepository
}

func NewContainer(db *gorm.DB, redisClient *redis.Client, cfg *config.Config) *Container {
	adminRepo := mysqlRepo.NewAdminRepository(db)
	userRepo := mysqlRepo.NewUserRepository(db)
	categoryRepo := mysqlRepo.NewCategoryRepository(db)
	audioRepo := mysqlRepo.NewAudioRepository(db)
	playlistRepo := mysqlRepo.NewPlaylistRepository(db)
	likeRepo := mysqlRepo.NewLikeRepository(db)
	historyRepo := mysqlRepo.NewHistoryRepository(db)
	otpRepo := mysqlRepo.NewOTPRepository(db)
	audioScoreRepo := mysqlRepo.NewAudioScoreRepository(db)

	emailSender := email.NewSender(cfg)
	colorExtractorSvc := service.NewColorExtractorService()

	adminAuthSvc := service.NewAdminAuthService(adminRepo)
	categorySvc := service.NewCategoryService(categoryRepo)
	audioSvc := service.NewAudioService(audioRepo, colorExtractorSvc)
	playlistSvc := service.NewPlaylistService(playlistRepo)
	likeSvc := service.NewLikeService(likeRepo)
	historySvc := service.NewHistoryService(historyRepo)
	otpSvc := service.NewOTPService(otpRepo, userRepo, emailSender)
	userAuthSvc := service.NewUserAuthService(userRepo, otpRepo, otpSvc)
	recommendationSvc := service.NewRecommendationService(audioRepo, audioScoreRepo, historyRepo, likeRepo)

	var cacheRepo port.CacheRepository
	if redisClient != nil {
		cacheRepo = cache.NewRedisCache(redisClient)
	}

	handlers := &router.Handlers{
		AdminAuth:      adminHandler.NewAuthHandler(adminAuthSvc),
		AdminCategory:  adminHandler.NewCategoryHandler(categorySvc),
		AdminAudio:     adminHandler.NewAudioHandler(audioSvc, historySvc),
		UserAuth:       userHandler.NewAuthHandler(userAuthSvc),
		UserPlaylist:   userHandler.NewPlaylistHandler(playlistSvc),
		UserLike:       userHandler.NewLikeHandler(likeSvc),
		UserHistory:    userHandler.NewHistoryHandler(historySvc),
		UserOTP:        userHandler.NewOTPHandler(otpSvc, userAuthSvc),
		UserRecommend:  userHandler.NewRecommendationHandler(recommendationSvc),
	}

	return &Container{
		Handlers:              handlers,
		RecommendationService: recommendationSvc,
		Cache:                 cacheRepo,
	}
}
