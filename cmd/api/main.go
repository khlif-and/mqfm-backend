package main

import (
	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/config"
	adminController "mqfm-backend/internal/controllers/auth/admin"
	userController "mqfm-backend/internal/controllers/auth/user"
	categoryAdminController "mqfm-backend/internal/controllers/category/admin"
	"mqfm-backend/internal/routes"
	adminService "mqfm-backend/internal/services/auth/admin"
	userService "mqfm-backend/internal/services/auth/user"
	categoryAdminService "mqfm-backend/internal/services/category/admin"
	"mqfm-backend/internal/utils"

)

func main() {
	config.ConnectDatabase()

	// 1. Init Admin Auth
	aService := adminService.NewAdminAuthService(config.DB)
	aController := adminController.NewAdminAuthController(aService)

	// 2. Init User Auth
	uService := userService.NewUserAuthService(config.DB)
	uController := userController.NewUserAuthController(uService)

	// 3. Init Admin Category
	catAdminService := categoryAdminService.NewAdminCategoryService(config.DB)
	catAdminController := categoryAdminController.NewAdminCategoryController(catAdminService)

	r := gin.Default()

	// Setup Routes dengan semua controller
	routes.SetupRoutes(r, aController, uController, catAdminController)

	utils.Log.Info("Server mqfm-backend berjalan di port 8080")
	if err := r.Run(":8080"); err != nil {
		utils.Log.Fatal("Gagal menjalankan server: " + err.Error())
	}
}