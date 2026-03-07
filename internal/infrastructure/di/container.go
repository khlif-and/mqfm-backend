package di

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"mqfm-backend/internal/adapter/cache"
	adminHandler "mqfm-backend/internal/adapter/handler/admin"
	publicHandler "mqfm-backend/internal/adapter/handler/public"
	userHandler "mqfm-backend/internal/adapter/handler/user"
	mysqlRepo "mqfm-backend/internal/adapter/repository/mysql"
	"mqfm-backend/internal/app/service"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/infrastructure/router"
)

type Container struct {
	Handlers           *router.Handlers
	LivestreamService  port.LivestreamService
	Cache              port.CacheRepository
}

func NewContainer(db *gorm.DB, redisClient *redis.Client, youtubeAPIKey string) *Container {
	adminRepo := mysqlRepo.NewAdminRepository(db)
	userRepo := mysqlRepo.NewUserRepository(db)
	categoryRepo := mysqlRepo.NewCategoryRepository(db)
	audioRepo := mysqlRepo.NewAudioRepository(db)
	playlistRepo := mysqlRepo.NewPlaylistRepository(db)
	likeRepo := mysqlRepo.NewLikeRepository(db)
	historyRepo := mysqlRepo.NewHistoryRepository(db)
	livestreamRepo := mysqlRepo.NewLivestreamRepository(db)

	adminAuthSvc := service.NewAdminAuthService(adminRepo)
	userAuthSvc := service.NewUserAuthService(userRepo)
	categorySvc := service.NewCategoryService(categoryRepo)
	audioSvc := service.NewAudioService(audioRepo)
	playlistSvc := service.NewPlaylistService(playlistRepo)
	likeSvc := service.NewLikeService(likeRepo)
	historySvc := service.NewHistoryService(historyRepo)
	livestreamSvc := service.NewLivestreamService(livestreamRepo, youtubeAPIKey)

	var cacheRepo port.CacheRepository
	if redisClient != nil {
		cacheRepo = cache.NewRedisCache(redisClient)
	}

	handlers := &router.Handlers{
		AdminAuth:     adminHandler.NewAuthHandler(adminAuthSvc),
		AdminCategory: adminHandler.NewCategoryHandler(categorySvc),
		AdminAudio:    adminHandler.NewAudioHandler(audioSvc, historySvc),
		UserAuth:      userHandler.NewAuthHandler(userAuthSvc),
		UserPlaylist:  userHandler.NewPlaylistHandler(playlistSvc),
		UserLike:      userHandler.NewLikeHandler(likeSvc),
		UserHistory:   userHandler.NewHistoryHandler(historySvc),
		Livestream:    publicHandler.NewLivestreamHandler(livestreamSvc),
	}

	return &Container{
		Handlers:          handlers,
		LivestreamService: livestreamSvc,
		Cache:             cacheRepo,
	}
}
