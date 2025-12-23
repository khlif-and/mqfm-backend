package routes

import (
	"github.com/gin-gonic/gin"

	adminController "mqfm-backend/internal/controllers/auth/admin"
	"mqfm-backend/internal/middleware"

)

func SetupRoutes(
	r *gin.Engine,
	aController *adminController.AdminAuthController,
) {
	api := r.Group("/api")
	{
		adminAuth := api.Group("/admin/auth")
		{
			adminAuth.POST("/register", aController.Register)
			adminAuth.POST("/login", aController.Login)

			protectedAdmin := adminAuth.Group("/")
			protectedAdmin.Use(middleware.JWTMiddleware())
			{
				protectedAdmin.PUT("/update/:id", aController.Update)
				protectedAdmin.POST("/logout", aController.Logout)
			}
		}
	}
}