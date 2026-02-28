package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"mqfm-backend/internal/config"
	adminController "mqfm-backend/internal/controllers/auth/admin"
	userController "mqfm-backend/internal/controllers/auth/user"
	catAdminController "mqfm-backend/internal/controllers/category/admin"
	historyUserController "mqfm-backend/internal/controllers/history/user"
	likeUserController "mqfm-backend/internal/controllers/likes/user"
	lsController "mqfm-backend/internal/controllers/livestream"
	playlistUserController "mqfm-backend/internal/controllers/playlist/user"
	audioAdminController "mqfm-backend/internal/controllers/podcast/audio/admin"
	lsModel "mqfm-backend/internal/models/livestream"
	userAuthRepo "mqfm-backend/internal/repositories/auth/user"
	catAdminRepo "mqfm-backend/internal/repositories/category/admin"
	historyRepo "mqfm-backend/internal/repositories/history/user"
	likeRepo "mqfm-backend/internal/repositories/likes/user"
	audioAdminRepo "mqfm-backend/internal/repositories/podcast/audio/admin"
	"mqfm-backend/internal/routes"
	adminAuthService "mqfm-backend/internal/services/auth/admin"
	userAuthService "mqfm-backend/internal/services/auth/user"
	catAdminService "mqfm-backend/internal/services/category/admin"
	historyUserService "mqfm-backend/internal/services/history/user"
	likeUserService "mqfm-backend/internal/services/likes/user"
	lsService "mqfm-backend/internal/services/livestream"
	playlistUserService "mqfm-backend/internal/services/playlist/user"
	audioAdminService "mqfm-backend/internal/services/podcast/audio/admin"
	"mqfm-backend/internal/utils"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	youtubeAPIKey := os.Getenv("YOUTUBE_API_KEY")
	if youtubeAPIKey == "" {
		log.Fatal("YOUTUBE_API_KEY is missing in .env")
	}

	config.ConnectDatabase()
	db := config.DB

	db.AutoMigrate(&lsModel.LiveStream{})

	r := gin.Default()
	r.Static("/uploads", "./uploads")

	adminRepo := adminAuthService.NewAdminAuthService(db)
	adminCtrl := adminController.NewAdminAuthController(adminRepo)

	userRepository := userAuthRepo.NewUserAuthRepository(db)
	userService := userAuthService.NewUserAuthService(userRepository)
	userCtrl := userController.NewUserAuthController(userService)

	categoryRepository := catAdminRepo.NewCategoryRepository(db)
	categoryService := catAdminService.NewAdminCategoryService(categoryRepository)
	catCtrl := catAdminController.NewAdminCategoryController(categoryService)

	audioRepository := audioAdminRepo.NewAudioRepository(db)
	audioService := audioAdminService.NewAdminAudioService(audioRepository)

	historyRepository := historyRepo.NewHistoryRepository(db)
	historyService := historyUserService.NewUserHistoryService(historyRepository)

	audioCtrl := audioAdminController.NewAdminAudioController(audioService, historyService)

	playlistRepo := playlistUserService.NewUserPlaylistService(db)
	playlistCtrl := playlistUserController.NewUserPlaylistController(playlistRepo)

	likeRepository := likeRepo.NewLikeRepository(db)
	likeService := likeUserService.NewUserLikeService(likeRepository)
	likeCtrl := likeUserController.NewUserLikeController(likeService)

	historyCtrl := historyUserController.NewUserHistoryController(historyService)
	mqfmChannelID := "UCwa0rj5KY6bWoVzJtgoiaDw"
	lsRepo := lsService.NewLiveStreamService(db, youtubeAPIKey)
	lsCtrl := lsController.NewLiveStreamController(lsRepo)

	go func() {
		utils.Log.Info("🚀 [Scheduler] Background Task Started: Checking YouTube Live Status...")
		for {
			if err := lsRepo.UpdateLiveStatus(mqfmChannelID); err != nil {
				utils.Log.Error("⚠️ [Scheduler] Error updating status", zap.Error(err))
			}
			time.Sleep(10 * time.Minute)
		}
	}()

	routes.SetupRoutes(r, adminCtrl, userCtrl, catCtrl, audioCtrl, playlistCtrl, likeCtrl, historyCtrl, lsCtrl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	utils.Log.Info("✅ Server running", zap.String("port", port))
	r.Run(":" + port)
}