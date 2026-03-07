package router

import (
	"github.com/gin-gonic/gin"

	adminHandler "mqfm-backend/internal/adapter/handler/admin"
	userHandler "mqfm-backend/internal/adapter/handler/user"
	"mqfm-backend/internal/infrastructure/middleware"
)

type Handlers struct {
	AdminAuth      *adminHandler.AuthHandler
	AdminCategory  *adminHandler.CategoryHandler
	AdminAudio     *adminHandler.AudioHandler
	UserAuth       *userHandler.AuthHandler
	UserPlaylist   *userHandler.PlaylistHandler
	UserLike       *userHandler.LikeHandler
	UserHistory    *userHandler.HistoryHandler
	UserOTP        *userHandler.OTPHandler
	UserRecommend  *userHandler.RecommendationHandler
}

func Setup(r *gin.Engine, h *Handlers) {
	r.Use(middleware.Security())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.RateLimit(10, 20))

	r.Static("/uploads", "./uploads")

	api := r.Group("/api")

	registerPublicRoutes(api, h)
	registerAdminRoutes(api, h)
	registerUserRoutes(api, h)
}
