package routes

import (
	"github.com/gin-gonic/gin"

	adminController "mqfm-backend/internal/controllers/auth/admin"
	userController "mqfm-backend/internal/controllers/auth/user"
	categoryAdminController "mqfm-backend/internal/controllers/category/admin"
	audioAdminController "mqfm-backend/internal/controllers/podcast/audio/admin" // Import baru
	"mqfm-backend/internal/middleware"

)

func SetupRoutes(
	r *gin.Engine,
	aController *adminController.AdminAuthController,
	uController *userController.UserAuthController,
	catAdminController *categoryAdminController.AdminCategoryController,
	audioAdminController *audioAdminController.AdminAudioController, // Parameter baru
) {
	api := r.Group("/api")
	{
		categories := api.Group("/categories")
		{
			categories.GET("/", catAdminController.FindAll)
			categories.GET("/:id", catAdminController.FindByID)
		}


audios := api.Group("/audios")
		{
			audios.GET("/", audioAdminController.FindAll)
		
			audios.GET("/search", audioAdminController.Search) 
			
			audios.GET("/:id", audioAdminController.FindByID)
		}

		// --- ADMIN Routes ---
		adminAuth := api.Group("/admin")
		{
			// Public Auth
			adminAuth.POST("/auth/register", aController.Register)
			adminAuth.POST("/auth/login", aController.Login)

			// Protected Admin Routes
			protectedAdmin := adminAuth.Group("/")
			protectedAdmin.Use(middleware.JWTMiddleware())
			{
				// Auth Profile
				protectedAdmin.GET("/auth/me", aController.Me)
				protectedAdmin.PUT("/auth/update/:id", aController.Update)
				protectedAdmin.POST("/auth/logout", aController.Logout)

				// Category Management (Admin Only)
				adminCategories := protectedAdmin.Group("/categories")
				{
					adminCategories.POST("/", catAdminController.Create)
					adminCategories.PUT("/:id", catAdminController.Update)
					adminCategories.DELETE("/:id", catAdminController.Delete)
				}

				// Audio Management (Admin Only)
				adminAudios := protectedAdmin.Group("/audios")
				{
					adminAudios.POST("/", audioAdminController.Create)
					adminAudios.PUT("/:id", audioAdminController.Update)
					adminAudios.DELETE("/:id", audioAdminController.Delete)
				}
			}
		}

		// --- USER Routes ---
		userAuth := api.Group("/user/auth")
		{
			userAuth.POST("/register", uController.Register)
			userAuth.POST("/login", uController.Login)

			protectedUser := userAuth.Group("/")
			protectedUser.Use(middleware.JWTMiddleware())
			{
				protectedUser.GET("/me", uController.Me)
				protectedUser.PUT("/update/:id", uController.Update)
				protectedUser.POST("/logout", uController.Logout)
			}
		}
	}
}