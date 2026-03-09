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

			events := protected.Group("/events")
			{
				events.POST("/", h.AdminEvent.Create)
				events.GET("/", h.AdminEvent.FindAll)
				events.GET("/:id", h.AdminEvent.FindByID)
				events.PUT("/:id", h.AdminEvent.Update)
				events.DELETE("/:id", h.AdminEvent.Delete)
			}

			series := protected.Group("/series")
			{
				series.POST("/", h.AdminSeries.Create)
				series.GET("/", h.AdminSeries.FindAll)
				series.GET("/search", h.AdminSeries.Search)
				series.GET("/:id", h.AdminSeries.FindByID)
				series.PUT("/:id", h.AdminSeries.Update)
				series.DELETE("/:id", h.AdminSeries.Delete)
				series.POST("/items", h.AdminSeries.AddItem)
				series.DELETE("/:id/items/:audio_id", h.AdminSeries.RemoveItem)
			}
		}
	}
}
