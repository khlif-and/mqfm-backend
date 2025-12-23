package main

import (
	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/config"
	adminController "mqfm-backend/internal/controllers/auth/admin"
	"mqfm-backend/internal/routes"
	adminService "mqfm-backend/internal/services/auth/admin"
	"mqfm-backend/internal/utils"

)

func main() {
	config.ConnectDatabase()

	aService := adminService.NewAdminAuthService(config.DB)
	aController := adminController.NewAdminAuthController(aService)

	r := gin.Default()

	routes.SetupRoutes(r, aController)

	utils.Log.Info("Server mqfm-backend berjalan di port 8080")
	if err := r.Run(":8080"); err != nil {
		utils.Log.Fatal("Gagal menjalankan server: " + err.Error())
	}
}