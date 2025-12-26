package user

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

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
	// Mengambil User ID dari context (pastikan middleware Auth berjalan)
	userID := utils.GetUserID(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	playlists, err := ctrl.service.GetByUserID(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch playlists", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Playlists retrieved", playlists)
}

func (ctrl *UserPlaylistController) AddAudio(c *gin.Context) {
	var input struct {
		AudioID         uint   `form:"audio_id" binding:"required"`
		PlaylistID      *uint  `form:"playlist_id"` // Pointer agar bisa cek null/0
		NewPlaylistName string `form:"new_playlist_name"`
	}

	if err := c.ShouldBind(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	userID := utils.GetUserID(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Skenario 1: Menambahkan ke playlist yang SUDAH ADA
	if input.PlaylistID != nil && *input.PlaylistID != 0 {
		if err := ctrl.service.AddAudioToPlaylist(userID, *input.PlaylistID, input.AudioID); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Failed to add audio to playlist", err.Error())
			return
		}
		utils.SuccessResponse(c, http.StatusOK, "Audio added to existing playlist", nil)
		return
	}

	// Skenario 2: Membuat Playlist BARU
	if input.NewPlaylistName == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Playlist name is required for new playlist", nil)
		return
	}

	file, _ := c.FormFile("image_file")
	var imagePath string

	if file != nil {
		pwd, _ := os.Getwd()
		uploadDir := filepath.Join(pwd, "uploads", "playlists")
		// Pastikan folder exists
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create directory", err.Error())
			return
		}

		// Format nama file: userid_timestamp_filename
		filename := fmt.Sprintf("%d_%d_%s", userID, time.Now().Unix(), file.Filename)
		fullPath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, fullPath); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload image", err.Error())
			return
		}
		imagePath = "uploads/playlists/" + filename
	}

	newPlaylist := playlistModel.Playlist{
		UserID:   userID,
		Name:     input.NewPlaylistName,
		ImageURL: imagePath,
	}

	// Gunakan service transaction untuk buat playlist + tambah audio
	if err := ctrl.service.CreateAndAddAudio(&newPlaylist, input.AudioID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create playlist and add audio", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "New playlist created and audio added", newPlaylist)
}

func (ctrl *UserPlaylistController) GetDetail(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
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