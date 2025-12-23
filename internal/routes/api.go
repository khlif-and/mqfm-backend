package routes

import (
	"github.com/gin-gonic/gin"

	adminController "mqfm-backend/internal/controllers/auth/admin"
	userController "mqfm-backend/internal/controllers/auth/user"
	categoryAdminController "mqfm-backend/internal/controllers/category/admin"
	"mqfm-backend/internal/middleware"

)

func SetupRoutes(
	r *gin.Engine,
	aController *adminController.AdminAuthController,
	uController *userController.UserAuthController,
	catAdminController *categoryAdminController.AdminCategoryController,
) {
	api := r.Group("/api")
	{
		// --- CATEGORY Routes (Universal / Public) ---
		// Endpoint ini bisa diakses siapa saja tanpa perlu login (User/Admin)
		categories := api.Group("/categories")
		{
			categories.GET("/", catAdminController.FindAll)
			categories.GET("/:id", catAdminController.FindByID)
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

				// Category Management (Admin Only - Write Operations)
				// Create, Update, Delete tetap butuh login admin
				adminCategories := protectedAdmin.Group("/categories")
				{
					adminCategories.POST("/", catAdminController.Create)
					adminCategories.PUT("/:id", catAdminController.Update)
					adminCategories.DELETE("/:id", catAdminController.Delete)
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