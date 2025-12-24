package main

import (
	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/config"
	adminController "mqfm-backend/internal/controllers/auth/admin"
	userController "mqfm-backend/internal/controllers/auth/user"
	categoryAdminController "mqfm-backend/internal/controllers/category/admin"
	audioAdminController "mqfm-backend/internal/controllers/podcast/audio/admin"
	"mqfm-backend/internal/routes"
	adminService "mqfm-backend/internal/services/auth/admin"
	userService "mqfm-backend/internal/services/auth/user"
	categoryAdminService "mqfm-backend/internal/services/category/admin"
	audioAdminService "mqfm-backend/internal/services/podcast/audio/admin"
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
	
	// UPDATE DISINI: Masukkan catAdminService sebagai parameter kedua
	audAdminController := audioAdminController.NewAdminAudioController(audAdminService, catAdminService)

	r := gin.Default()

	routes.SetupRoutes(r, aController, uController, catAdminController, audAdminController)

	utils.Log.Info("Server mqfm-backend berjalan di port 8080")
	if err := r.Run(":8080"); err != nil {
		utils.Log.Fatal("Gagal menjalankan server: " + err.Error())
	}
}