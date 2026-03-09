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
				playlists.PUT("/:id", h.UserPlaylist.Update)
				playlists.DELETE("/:id", h.UserPlaylist.Delete)
				playlists.POST("/add-audio", h.UserPlaylist.AddAudio)
				playlists.POST("/remove-audio", h.UserPlaylist.RemoveAudio)
				playlists.POST("/:id/share", h.UserPlaylist.Share)
			}

			likes := protected.Group("/likes")
			{
				likes.POST("/", h.UserLike.Like)
				likes.DELETE("/", h.UserLike.Unlike)
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
				recommend.GET("/location-based", h.UserRecommend.GetLocationBased)
				recommend.GET("/time-based", h.UserRecommend.GetTimeBased)
			}

			bookmarks := protected.Group("/bookmarks")
			{
				bookmarks.POST("/", h.UserBookmark.Create)
				bookmarks.GET("/", h.UserBookmark.GetByUser)
				bookmarks.GET("/audio/:audio_id", h.UserBookmark.GetByAudio)
				bookmarks.DELETE("/:id", h.UserBookmark.Delete)
			}

			notifications := protected.Group("/notifications")
			{
				notifications.GET("/", h.UserNotification.GetNotifications)
				notifications.PUT("/:id/read", h.UserNotification.MarkAsRead)
				notifications.PUT("/read-all", h.UserNotification.MarkAllAsRead)
				notifications.GET("/unread-count", h.UserNotification.UnreadCount)
				notifications.GET("/settings", h.UserNotification.GetSetting)
				notifications.PUT("/settings", h.UserNotification.UpdateSetting)
			}

			progress := protected.Group("/progress")
			{
				progress.POST("/", h.UserProgress.Update)
				progress.GET("/", h.UserProgress.GetAll)
				progress.GET("/completed", h.UserProgress.GetCompleted)
				progress.GET("/:audio_id", h.UserProgress.Get)
			}

			downloads := protected.Group("/downloads")
			{
				downloads.POST("/", h.UserDownload.Record)
				downloads.GET("/", h.UserDownload.GetAll)
				downloads.DELETE("/:id", h.UserDownload.Delete)
				downloads.GET("/storage", h.UserDownload.StorageUsage)
				downloads.GET("/smart", h.UserDownload.SmartDownload)
			}

			stats := protected.Group("/stats")
			{
				stats.POST("/", h.UserStats.RecordStat)
				stats.GET("/recap", h.UserStats.GetRecap)
			}

			clips := protected.Group("/clips")
			{
				clips.POST("/", h.UserClip.Create)
				clips.GET("/", h.UserClip.GetByUser)
				clips.DELETE("/:id", h.UserClip.Delete)
			}

			events := protected.Group("/events")
			{
				events.GET("/upcoming", h.UserEvent.GetUpcoming)
				events.GET("/:id", h.UserEvent.GetByID)
				events.POST("/:id/rsvp", h.UserEvent.RSVP)
				events.DELETE("/:id/rsvp", h.UserEvent.CancelRSVP)
				events.GET("/my-rsvps", h.UserEvent.GetMyRSVPs)
			}

			preferences := protected.Group("/preferences")
			{
				preferences.GET("/", h.UserPreference.Get)
				preferences.PUT("/", h.UserPreference.Update)
			}

			votes := protected.Group("/votes")
			{
				votes.POST("/", h.UserVote.Vote)
				votes.DELETE("/:audio_id", h.UserVote.Unvote)
				votes.GET("/status/:audio_id", h.UserVote.Status)
				votes.GET("/ranking/weekly", h.UserVote.WeeklyRanking)
				votes.GET("/ranking/monthly", h.UserVote.MonthlyRanking)
			}

			resume := protected.Group("/resume")
			{
				resume.POST("/", h.UserResume.Update)
				resume.GET("/", h.UserResume.Get)
			}

			share := protected.Group("/share")
			{
				share.GET("/audio/:audio_id", h.UserShare.ShareAudio)
			}

			favArtists := protected.Group("/favorite-artists")
			{
				favArtists.POST("/", h.UserFavArtist.Add)
				favArtists.DELETE("/:artist", h.UserFavArtist.Remove)
				favArtists.GET("/", h.UserFavArtist.GetAll)
			}

			location := protected.Group("/location")
			{
				location.PUT("/", h.UserLocation.Update)
				location.GET("/", h.UserLocation.Get)
			}

			collab := protected.Group("/collab")
			{
				collab.POST("/collaborators", h.UserCollab.AddCollaborator)
				collab.DELETE("/playlists/:playlist_id/collaborators/:user_id", h.UserCollab.RemoveCollaborator)
				collab.GET("/playlists/:playlist_id/collaborators", h.UserCollab.GetCollaborators)
				collab.POST("/contribute", h.UserCollab.ContributeAudio)
				collab.POST("/join", h.UserCollab.JoinByToken)
			}
		}
	}
}
