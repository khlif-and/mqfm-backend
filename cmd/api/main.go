package main

import (
	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/config"
	adminController "mqfm-backend/internal/controllers/auth/admin"
	userController "mqfm-backend/internal/controllers/auth/user"
	categoryAdminController "mqfm-backend/internal/controllers/category/admin"
	audioAdminController "mqfm-backend/internal/controllers/podcast/audio/admin"
	playlistUserController "mqfm-backend/internal/controllers/playlist/user"
	likeUserController "mqfm-backend/internal/controllers/likes/user" 
	"mqfm-backend/internal/routes"
	adminService "mqfm-backend/internal/services/auth/admin"
	userService "mqfm-backend/internal/services/auth/user"
	categoryAdminService "mqfm-backend/internal/services/category/admin"
	audioAdminService "mqfm-backend/internal/services/podcast/audio/admin"
	playlistUserService "mqfm-backend/internal/services/playlist/user"
	likeUserService "mqfm-backend/internal/services/likes/user" 
	"mqfm-backend/internal/utils"

)

func main() {
	config.ConnectDatabase()

	aService := adminService.NewAdminAuthService(config.DB)
	aController := adminController.NewAdminAuthController(aService)

	uService := userService.NewUserAuthService(config.DB)
	uController := userController.NewUserAuthController(uService)

	catAdminService := categoryAdminService.NewAdminCategoryService(config.DB)
	catAdminController := categoryAdminController.NewAdminCategoryController(catAdminService)

	audAdminService := audioAdminService.NewAdminAudioService(config.DB)
	audAdminController := audioAdminController.NewAdminAudioController(audAdminService, catAdminService)

	plUserService := playlistUserService.NewUserPlaylistService(config.DB)
	plUserController := playlistUserController.NewUserPlaylistController(plUserService)

	lUserService := likeUserService.NewUserLikeService(config.DB)
	lUserController := likeUserController.NewUserLikeController(lUserService)

	r := gin.Default()

	r.Static("/uploads", "./uploads")

	routes.SetupRoutes(r, aController, uController, catAdminController, audAdminController, plUserController, lUserController)

	utils.Log.Info("Server mqfm-backend berjalan di port 8080")
	if err := r.Run(":8080"); err != nil {
		utils.Log.Fatal("Gagal menjalankan server: " + err.Error())
	}
}