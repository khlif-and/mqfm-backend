package routes

import (
	"github.com/gin-gonic/gin"

	adminController "mqfm-backend/internal/controllers/auth/admin"
	userController "mqfm-backend/internal/controllers/auth/user"
	categoryAdminController "mqfm-backend/internal/controllers/category/admin"
	likeUserController "mqfm-backend/internal/controllers/likes/user"
	lsController "mqfm-backend/internal/controllers/livestream"
	playlistUserController "mqfm-backend/internal/controllers/playlist/user"
	audioAdminController "mqfm-backend/internal/controllers/podcast/audio/admin"
	"mqfm-backend/internal/middleware"

)

func SetupRoutes(
	r *gin.Engine,
	aController *adminController.AdminAuthController,
	uController *userController.UserAuthController,
	catAdminController *categoryAdminController.AdminCategoryController,
	audioAdminController *audioAdminController.AdminAudioController,
	playlistController *playlistUserController.UserPlaylistController,
	likeController *likeUserController.UserLikeController,
	lsController *lsController.LiveStreamController,
) {
	api := r.Group("/api")
	{
		categories := api.Group("/categories")
		{
			categories.GET("/", catAdminController.FindAll)
			categories.GET("/search", catAdminController.Search)
			categories.GET("/:id", catAdminController.FindByID)
		}

		audios := api.Group("/audios")
		{
			audios.GET("/", audioAdminController.FindAll)
			audios.GET("/search", audioAdminController.Search)
			audios.GET("/:id", audioAdminController.FindByID)
		}

		youtube := api.Group("/youtube")
		{
			youtube.GET("/live-status", lsController.GetStatus)
		}

		adminAuth := api.Group("/admin")
		{
			adminAuth.POST("/auth/register", aController.Register)
			adminAuth.POST("/auth/login", aController.Login)

			protectedAdmin := adminAuth.Group("/")
			protectedAdmin.Use(middleware.JWTMiddleware())
			{
				protectedAdmin.GET("/auth/me", aController.Me)
				protectedAdmin.PUT("/auth/update/:id", aController.Update)
				protectedAdmin.POST("/auth/logout", aController.Logout)

				adminCategories := protectedAdmin.Group("/categories")
				{
					adminCategories.POST("/", catAdminController.Create)
					adminCategories.PUT("/:id", catAdminController.Update)
					adminCategories.DELETE("/:id", catAdminController.Delete)
				}

				adminAudios := protectedAdmin.Group("/audios")
				{
					adminAudios.POST("/", audioAdminController.Create)
					adminAudios.PUT("/:id", audioAdminController.Update)
					adminAudios.DELETE("/:id", audioAdminController.Delete)
				}
			}
		}

		userAuth := api.Group("/user")
		{
			userAuth.POST("/auth/register", uController.Register)
			userAuth.POST("/auth/login", uController.Login)

			protectedUser := userAuth.Group("/")
			protectedUser.Use(middleware.JWTMiddleware())
			{
				protectedUser.GET("/auth/me", uController.Me)
				protectedUser.PUT("/auth/update/:id", uController.Update)
				protectedUser.POST("/auth/logout", uController.Logout)

				playlists := protectedUser.Group("/playlists")
				{
					playlists.GET("/", playlistController.GetMyPlaylists)
					playlists.GET("/search", playlistController.Search)
					playlists.GET("/:id", playlistController.GetDetail)
					playlists.POST("/", playlistController.Create)
					playlists.POST("/add-audio", playlistController.AddAudio)
				}

				likes := protectedUser.Group("/likes")
				{
					likes.POST("/", likeController.Like)
					likes.DELETE("/:audio_id", likeController.Unlike)
					likes.GET("/", likeController.GetLikes)
				}
			}
		}
	}
}