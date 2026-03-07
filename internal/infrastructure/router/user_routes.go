package router

import (
	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/infrastructure/middleware"
)

func registerUserRoutes(api *gin.RouterGroup, h *Handlers) {
	user := api.Group("/user")
	{
		user.POST("/auth/register", h.UserAuth.Register)
		user.POST("/auth/login", h.UserAuth.Login)
		user.POST("/auth/google", h.UserAuth.GoogleLogin)

		protected := user.Group("/")
		protected.Use(middleware.JWTAuth())
		{
			protected.GET("/auth/me", h.UserAuth.Me)
			protected.PUT("/auth/update/:id", h.UserAuth.Update)
			protected.POST("/auth/logout", h.UserAuth.Logout)

			playlists := protected.Group("/playlists")
			{
				playlists.GET("/", h.UserPlaylist.GetMyPlaylists)
				playlists.GET("/search", h.UserPlaylist.Search)
				playlists.GET("/:id", h.UserPlaylist.GetDetail)
				playlists.POST("/", h.UserPlaylist.Create)
				playlists.POST("/add-audio", h.UserPlaylist.AddAudio)
			}

			likes := protected.Group("/likes")
			{
				likes.POST("/", h.UserLike.Like)
				likes.DELETE("/:audio_id", h.UserLike.Unlike)
				likes.GET("/", h.UserLike.GetLikes)
			}

			history := protected.Group("/history")
			{
				history.GET("/", h.UserHistory.GetHistory)
				history.DELETE("/:audio_id", h.UserHistory.DeleteHistory)
				history.DELETE("/clear", h.UserHistory.ClearHistory)
			}
		}
	}
}
