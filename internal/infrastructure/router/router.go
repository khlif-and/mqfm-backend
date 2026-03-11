package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	adminHandler "mqfm-backend/internal/adapter/handler/admin"
	userHandler "mqfm-backend/internal/adapter/handler/user"
	"mqfm-backend/internal/infrastructure/middleware"
)

type Handlers struct {
	AdminAuth        *adminHandler.AuthHandler
	AdminCategory    *adminHandler.CategoryHandler
	AdminAudio       *adminHandler.AudioHandler
	AdminPlaylist    *adminHandler.PlaylistHandler
	AdminEvent       *adminHandler.EventHandler
	AdminSeries      *adminHandler.SeriesHandler
	UserAuth         *userHandler.AuthHandler
	UserPlaylist     *userHandler.PlaylistHandler
	UserLike         *userHandler.LikeHandler
	UserHistory      *userHandler.HistoryHandler
	UserOTP          *userHandler.OTPHandler
	UserRecommend    *userHandler.RecommendationHandler
	UserBookmark     *userHandler.BookmarkHandler
	UserNotification *userHandler.NotificationHandler
	UserProgress     *userHandler.ProgressHandler
	UserDownload     *userHandler.DownloadHandler
	UserStats        *userHandler.StatsHandler
	UserClip         *userHandler.ClipHandler
	UserEvent        *userHandler.EventHandler
	UserPreference   *userHandler.PreferenceHandler
	UserVote         *userHandler.VoteHandler
	UserResume       *userHandler.ResumeHandler
	UserShare        *userHandler.ShareHandler
	UserFavArtist    *userHandler.FavoriteArtistHandler
	UserLocation     *userHandler.LocationHandler
	UserCollab       *userHandler.PlaylistCollabHandler
}

func Setup(r *gin.Engine, h *Handlers) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(middleware.Security())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.RateLimit(100, 200))

	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	registerPublicRoutes(api, h)
	registerAdminRoutes(api, h)
	registerUserRoutes(api, h)

	apiv1 := r.Group("/api/v1")
	registerPublicRoutes(apiv1, h)
	registerAdminRoutes(apiv1, h)
	registerUserRoutes(apiv1, h)
}
