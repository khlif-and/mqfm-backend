package router

import (
	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/infrastructure/middleware"
)

func registerPublicRoutes(api *gin.RouterGroup, h *Handlers) {
	categories := api.Group("/categories")
	{
		categories.GET("/", h.AdminCategory.FindAll)
		categories.GET("/search", h.AdminCategory.Search)
		categories.GET("/:id", h.AdminCategory.FindByID)
	}

	audios := api.Group("/audios")
	{
		audios.GET("/", h.AdminAudio.FindAll)
		audios.GET("/search", h.AdminAudio.Search)
		audios.GET("/:id", middleware.OptionalJWTAuth(), h.AdminAudio.FindByID)
	}

	youtube := api.Group("/youtube")
	{
		youtube.GET("/live-status", h.Livestream.GetStatus)
	}
}
