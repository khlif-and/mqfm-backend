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
	"mqfm-backend/internal/shared/helper"
)

type Container struct {
	Handlers              *router.Handlers
	RecommendationService port.RecommendationService
	AudioVoteService      port.AudioVoteService
	NotificationService   port.NotificationService
	Cache                 port.CacheRepository
	CacheManager          port.CacheManager
	TokenStore            port.TokenStore
	RankingCache          *cache.RankingCache
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
	bookmarkRepo := mysqlRepo.NewBookmarkRepository(db)
	notificationRepo := mysqlRepo.NewNotificationRepository(db)
	progressRepo := mysqlRepo.NewAudioProgressRepository(db)
	downloadRepo := mysqlRepo.NewDownloadRepository(db)
	statRepo := mysqlRepo.NewListeningStatRepository(db)
	clipRepo := mysqlRepo.NewAudioClipRepository(db)
	eventRepo := mysqlRepo.NewEventRepository(db)
	prefRepo := mysqlRepo.NewUserPreferenceRepository(db)
	seriesRepo := mysqlRepo.NewAudioSeriesRepository(db)
	voteRepo := mysqlRepo.NewAudioVoteRepository(db)
	rankingRepo := mysqlRepo.NewAudioRankingRepository(db)
	favArtistRepo := mysqlRepo.NewFavoriteArtistRepository(db)
	resumeRepo := mysqlRepo.NewSmartResumeRepository(db)
	locationRepo := mysqlRepo.NewUserLocationRepository(db)
	collabRepo := mysqlRepo.NewPlaylistCollaboratorRepository(db)

	emailSender := email.NewSender(cfg)
	colorExtractorSvc := service.NewColorExtractorService()
	audioConverter := helper.NewAudioConverter()

	var cacheRepo port.CacheRepository
	var cacheMgr port.CacheManager
	var tokenStore port.TokenStore
	var rankingCache *cache.RankingCache

	if redisClient != nil {
		cacheRepo = cache.NewRedisCache(redisClient)
		cacheMgr = cache.NewCacheManager(cacheRepo, audioRepo)
		tokenStore = cache.NewRedisTokenStore(redisClient)
		rankingCache = cache.NewRankingCache(cacheRepo, audioScoreRepo)
	}

	adminAuthSvc := service.NewAdminAuthService(adminRepo, tokenStore)
	categorySvc := service.NewCategoryService(categoryRepo)
	audioSvc := service.NewAudioService(audioRepo, colorExtractorSvc)
	playlistSvc := service.NewPlaylistService(playlistRepo, colorExtractorSvc)
	likeSvc := service.NewLikeService(likeRepo)
	historySvc := service.NewHistoryService(historyRepo)
	otpSvc := service.NewOTPService(otpRepo, userRepo, emailSender)
	userAuthSvc := service.NewUserAuthService(userRepo, otpRepo, otpSvc, tokenStore)
	recommendationSvc := service.NewRecommendationService(audioRepo, audioScoreRepo, historyRepo, likeRepo, locationRepo)
	bookmarkSvc := service.NewBookmarkService(bookmarkRepo)
	notificationSvc := service.NewNotificationService(notificationRepo)
	progressSvc := service.NewAudioProgressService(progressRepo)
	downloadSvc := service.NewDownloadService(downloadRepo, audioRepo, favArtistRepo)
	statSvc := service.NewListeningStatService(statRepo)
	clipSvc := service.NewAudioClipService(clipRepo, audioRepo, audioConverter)
	eventSvc := service.NewEventService(eventRepo, notificationSvc)
	prefSvc := service.NewUserPreferenceService(prefRepo)
	seriesSvc := service.NewAudioSeriesService(seriesRepo)
	voteSvc := service.NewAudioVoteService(voteRepo, rankingRepo, audioRepo)
	resumeSvc := service.NewSmartResumeService(resumeRepo)
	shareSvc := service.NewShareService(cfg.BaseURL)
	favArtistSvc := service.NewFavoriteArtistService(favArtistRepo)
	locationSvc := service.NewUserLocationService(locationRepo)
	collabSvc := service.NewPlaylistCollabService(collabRepo, playlistRepo)

	handlers := &router.Handlers{
		AdminAuth:      adminHandler.NewAuthHandler(adminAuthSvc),
		AdminCategory:  adminHandler.NewCategoryHandler(categorySvc),
		AdminAudio:     adminHandler.NewAudioHandler(audioSvc, historySvc),
		AdminPlaylist:  adminHandler.NewPlaylistHandler(playlistSvc),
		AdminEvent:     adminHandler.NewEventHandler(eventSvc),
		AdminSeries:    adminHandler.NewSeriesHandler(seriesSvc),
		UserAuth:       userHandler.NewAuthHandler(userAuthSvc),
		UserPlaylist:   userHandler.NewPlaylistHandler(playlistSvc),
		UserLike:       userHandler.NewLikeHandler(likeSvc),
		UserHistory:    userHandler.NewHistoryHandler(historySvc),
		UserOTP:        userHandler.NewOTPHandler(otpSvc, userAuthSvc),
		UserRecommend:  userHandler.NewRecommendationHandler(recommendationSvc),
		UserBookmark:   userHandler.NewBookmarkHandler(bookmarkSvc),
		UserNotification: userHandler.NewNotificationHandler(notificationSvc),
		UserProgress:   userHandler.NewProgressHandler(progressSvc),
		UserDownload:   userHandler.NewDownloadHandler(downloadSvc),
		UserStats:      userHandler.NewStatsHandler(statSvc),
		UserClip:       userHandler.NewClipHandler(clipSvc, shareSvc),
		UserEvent:      userHandler.NewEventHandler(eventSvc),
		UserPreference: userHandler.NewPreferenceHandler(prefSvc),
		UserVote:       userHandler.NewVoteHandler(voteSvc),
		UserResume:     userHandler.NewResumeHandler(resumeSvc),
		UserShare:      userHandler.NewShareHandler(shareSvc, playlistSvc),
		UserFavArtist:  userHandler.NewFavoriteArtistHandler(favArtistSvc),
		UserLocation:   userHandler.NewLocationHandler(locationSvc),
		UserCollab:     userHandler.NewPlaylistCollabHandler(collabSvc),
	}

	return &Container{
		Handlers:              handlers,
		RecommendationService: recommendationSvc,
		AudioVoteService:      voteSvc,
		NotificationService:   notificationSvc,
		Cache:                 cacheRepo,
		CacheManager:          cacheMgr,
		TokenStore:            tokenStore,
		RankingCache:          rankingCache,
	}
}
