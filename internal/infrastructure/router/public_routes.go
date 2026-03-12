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

	shared := api.Group("/shared")
	{
		shared.GET("/clip/:token", h.UserClip.GetByShareToken)
		shared.GET("/playlist/:token", h.UserShare.GetSharedPlaylist)
	}

	api.GET("/series", h.AdminSeries.FindAll)
	api.GET("/series/search", h.AdminSeries.Search)
	api.GET("/series/:id", h.AdminSeries.FindByID)
	api.GET("/events/upcoming", middleware.OptionalJWTAuth(), h.UserEvent.GetUpcoming)
	api.GET("/events/:id", middleware.OptionalJWTAuth(), h.UserEvent.GetByID)
	api.GET("/votes/ranking/weekly", h.UserVote.WeeklyRanking)
	api.GET("/votes/ranking/monthly", h.UserVote.MonthlyRanking)

	api.GET("/radios", h.AdminRadio.FindActive)
	api.GET("/radios/:id", h.AdminRadio.FindByID)
}
