package router

import (
	"github.com/gin-gonic/gin"

	adminHandler "mqfm-backend/internal/adapter/handler/admin"
	publicHandler "mqfm-backend/internal/adapter/handler/public"
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
	Livestream     *publicHandler.LivestreamHandler
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
