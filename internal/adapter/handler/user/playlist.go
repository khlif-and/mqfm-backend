package user

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/helper"
	"mqfm-backend/internal/shared/logger"
	resp "mqfm-backend/internal/shared/response"
	"mqfm-backend/internal/shared/security"
)

type PlaylistHandler struct {
	service port.PlaylistService
}

func NewPlaylistHandler(s port.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{service: s}
}

func (h *PlaylistHandler) GetMyPlaylists(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	playlists, err := h.service.GetByUserID(userID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgPlaylistListFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgPlaylistListOK, playlists)
}

func (h *PlaylistHandler) Search(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	query := c.Query("q")
	if query == "" {
		resp.Error(c, http.StatusBadRequest, constant.MsgSearchRequired, nil)
		return
	}

	playlists, err := h.service.Search(userID, query)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgPlaylistSearchFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgPlaylistSearchOK, playlists)
}

func (h *PlaylistHandler) Create(c *gin.Context) {
	var input struct {
		Name string `form:"name" binding:"required"`
	}
	if err := c.ShouldBind(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	file, _ := c.FormFile("image_file")
	var imagePath string

	if file != nil {
		pwd, _ := os.Getwd()
		uploadDir := filepath.Join(pwd, "uploads", "playlists")
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			logger.Error("failed to create upload dir")
			resp.Error(c, http.StatusInternalServerError, constant.MsgDirCreateFail, err.Error())
			return
		}

		filename := fmt.Sprintf("%d_%d_%s", userID, time.Now().Unix(), file.Filename)
		fullPath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, fullPath); err != nil {
			logger.Error("failed to save playlist image")
			resp.Error(c, http.StatusInternalServerError, constant.MsgFileUploadFail, err.Error())
			return
		}
		imagePath = "uploads/playlists/" + filename
	}

	playlist := entity.Playlist{
		UserID:   userID,
		Name:     input.Name,
		ImageURL: imagePath,
	}

	if err := h.service.Create(&playlist); err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgPlaylistCreateFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgPlaylistCreateOK, playlist)
}

func (h *PlaylistHandler) AddAudio(c *gin.Context) {
	var input struct {
		AudioID    uint `json:"audio_id" binding:"required"`
		PlaylistID uint `json:"playlist_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	if err := h.service.AddAudioToPlaylist(userID, input.PlaylistID, input.AudioID); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgPlaylistAddAudioFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgPlaylistAddAudioOK, nil)
}

func (h *PlaylistHandler) GetDetail(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	playlist, err := h.service.GetByID(id, userID)
	if err != nil {
		resp.Error(c, http.StatusNotFound, constant.MsgPlaylistNotFound, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgPlaylistGetOK, playlist)
}

func (h *PlaylistHandler) RemoveAudio(c *gin.Context) {
	var input struct {
		AudioID    uint `json:"audio_id" binding:"required"`
		PlaylistID uint `json:"playlist_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	if err := h.service.RemoveAudioFromPlaylist(userID, input.PlaylistID, input.AudioID); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgPlaylistRemoveAudioFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgPlaylistRemoveAudioOK, nil)
}
