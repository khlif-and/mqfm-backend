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
	likeUserController "mqfm-backend/internal/controllers/likes/user"
	lsController "mqfm-backend/internal/controllers/livestream"
	playlistUserController "mqfm-backend/internal/controllers/playlist/user"
	audioAdminController "mqfm-backend/internal/controllers/podcast/audio/admin"
	lsModel "mqfm-backend/internal/models/livestream"
	"mqfm-backend/internal/routes"
	adminAuthService "mqfm-backend/internal/services/auth/admin"
	userAuthRepo "mqfm-backend/internal/repositories/auth/user"
	userAuthService "mqfm-backend/internal/services/auth/user"
	catAdminService "mqfm-backend/internal/services/category/admin"
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

	catRepo := catAdminService.NewAdminCategoryService(db)
	catCtrl := catAdminController.NewAdminCategoryController(catRepo)

	audioRepo := audioAdminService.NewAdminAudioService(db)
	audioCtrl := audioAdminController.NewAdminAudioController(audioRepo, catRepo)

	playlistRepo := playlistUserService.NewUserPlaylistService(db)
	playlistCtrl := playlistUserController.NewUserPlaylistController(playlistRepo)

	likeRepo := likeUserService.NewUserLikeService(db)
	likeCtrl := likeUserController.NewUserLikeController(likeRepo)

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

	routes.SetupRoutes(r, adminCtrl, userCtrl, catCtrl, audioCtrl, playlistCtrl, likeCtrl, lsCtrl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	utils.Log.Info("✅ Server running", zap.String("port", port))
	r.Run(":" + port)
}