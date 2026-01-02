package user

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	playlistModel "mqfm-backend/internal/models/playlist/user"
	playlistService "mqfm-backend/internal/services/playlist/user"
	"mqfm-backend/internal/utils"
)

type UserPlaylistController struct {
	service *playlistService.UserPlaylistService
}

func NewUserPlaylistController(s *playlistService.UserPlaylistService) *UserPlaylistController {
	return &UserPlaylistController{service: s}
}

func (ctrl *UserPlaylistController) GetMyPlaylists(c *gin.Context) {
	userID := utils.GetUserID(c)
	if userID == 0 {
		utils.Log.Warn("[Controller] Unauthorized access attempt",
			zap.String("ip", c.ClientIP()),
			zap.String("endpoint", "GetMyPlaylists"),
		)
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	playlists, err := ctrl.service.GetByUserID(userID)
	if err != nil {
		// Log error sudah ada di service, tapi kita bisa log context HTTP-nya disini
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch playlists", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Playlists retrieved", playlists)
}

func (ctrl *UserPlaylistController) Search(c *gin.Context) {
	userID := utils.GetUserID(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	query := c.Query("q")
	if query == "" {
		utils.Log.Info("[Controller] Search query missing",
			zap.Uint("user_id", userID),
		)
		utils.ErrorResponse(c, http.StatusBadRequest, "Search query is required", nil)
		return
	}

	playlists, err := ctrl.service.Search(userID, query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to search playlists", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Playlists found", playlists)
}

func (ctrl *UserPlaylistController) Create(c *gin.Context) {
	var input struct {
		Name string `form:"name" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		utils.Log.Warn("[Controller] Invalid playlist input",
			zap.Error(err),
			zap.String("ip", c.ClientIP()),
		)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	userID := utils.GetUserID(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	file, _ := c.FormFile("image_file")
	var imagePath string

	if file != nil {
		pwd, _ := os.Getwd()
		uploadDir := filepath.Join(pwd, "uploads", "playlists")
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			utils.Log.Error("[Controller] Failed to create upload directory",
				zap.Error(err),
				zap.String("path", uploadDir),
			)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create directory", err.Error())
			return
		}

		filename := fmt.Sprintf("%d_%d_%s", userID, time.Now().Unix(), file.Filename)
		fullPath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, fullPath); err != nil {
			utils.Log.Error("[Controller] Failed to save uploaded file",
				zap.Error(err),
				zap.String("path", fullPath),
			)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload image", err.Error())
			return
		}
		imagePath = "uploads/playlists/" + filename
	}

	newPlaylist := playlistModel.Playlist{
		UserID:   userID,
		Name:     input.Name,
		ImageURL: imagePath,
	}

	if err := ctrl.service.Create(&newPlaylist); err != nil {
		// Error database sudah dilog di Service
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create playlist", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Playlist created successfully", newPlaylist)
}

func (ctrl *UserPlaylistController) AddAudio(c *gin.Context) {
	var input struct {
		AudioID    uint `json:"audio_id" binding:"required"`
		PlaylistID uint `json:"playlist_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warn("[Controller] Invalid AddAudio input",
			zap.Error(err),
			zap.String("ip", c.ClientIP()),
		)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	userID := utils.GetUserID(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	if err := ctrl.service.AddAudioToPlaylist(userID, input.PlaylistID, input.AudioID); err != nil {
		// Log specific logic errors from service as Info/Warn in controller context if needed
		// But mostly service logs capture the system error.
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to add audio", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Audio added successfully", nil)
}

func (ctrl *UserPlaylistController) GetDetail(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.Log.Info("[Controller] Invalid Playlist ID param",
			zap.String("id_param", idParam),
		)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	userID := utils.GetUserID(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	playlist, err := ctrl.service.GetByID(uint(id), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Playlist not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Playlist detail retrieved", playlist)
}