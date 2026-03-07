package router

import (
	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/infrastructure/middleware"
)

func registerAdminRoutes(api *gin.RouterGroup, h *Handlers) {
	admin := api.Group("/admin")
	{
		admin.POST("/auth/register", h.AdminAuth.Register)
		admin.POST("/auth/login", h.AdminAuth.Login)

		protected := admin.Group("/")
		protected.Use(middleware.JWTAuth())
		{
			protected.GET("/auth/me", h.AdminAuth.Me)
			protected.PUT("/auth/update/:id", h.AdminAuth.Update)
			protected.POST("/auth/logout", h.AdminAuth.Logout)

			categories := protected.Group("/categories")
			{
				categories.POST("/", h.AdminCategory.Create)
				categories.PUT("/:id", h.AdminCategory.Update)
				categories.DELETE("/:id", h.AdminCategory.Delete)
			}

			audios := protected.Group("/audios")
			{
				audios.POST("/", h.AdminAudio.Create)
				audios.PUT("/:id", h.AdminAudio.Update)
				audios.DELETE("/:id", h.AdminAudio.Delete)
			}
		}
	}
}
