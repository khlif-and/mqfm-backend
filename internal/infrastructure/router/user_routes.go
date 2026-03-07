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
		user.POST("/auth/send-otp", h.UserOTP.SendOTP)
		user.POST("/auth/verify-otp", h.UserOTP.VerifyOTP)

		protected := user.Group("/")
		protected.Use(middleware.JWTAuth())
		{
			protected.GET("/auth/me", h.UserAuth.Me)
			protected.PUT("/auth/update/:id", h.UserAuth.Update)
			protected.POST("/auth/logout", h.UserAuth.Logout)
			protected.POST("/auth/google/link", h.UserAuth.LinkGoogle)
			protected.DELETE("/auth/google/unlink", h.UserAuth.UnlinkGoogle)

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
				history.DELETE("/clear", h.UserHistory.ClearHistory)
				history.DELETE("/:audio_id", h.UserHistory.DeleteHistory)
			}

			recommend := protected.Group("/recommendations")
			{
				recommend.GET("/popular", h.UserRecommend.GetPopular)
				recommend.GET("/most-listened", h.UserRecommend.GetMostListened)
				recommend.GET("/by-artist", h.UserRecommend.GetByArtist)
				recommend.GET("/similar/:audio_id", h.UserRecommend.GetSimilar)
				recommend.GET("/quick-pick", h.UserRecommend.GetQuickPick)
				recommend.GET("/onboarding", h.UserRecommend.GetOnboarding)
				recommend.GET("/personalized", h.UserRecommend.GetPersonalized)
			}
		}
	}
}
