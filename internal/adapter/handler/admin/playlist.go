package admin

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
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/internal/shared/dto/response"
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

func (h *PlaylistHandler) Create(c *gin.Context) {
	var input request.CreatePlaylistRequest
	if err := c.ShouldBind(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	adminID := security.GetUserID(c)
	if adminID == 0 {
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

		filename := fmt.Sprintf("admin_%d_%d_%s", adminID, time.Now().Unix(), filepath.Base(file.Filename))
		fullPath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, fullPath); err != nil {
			logger.Error("failed to save playlist image")
			resp.Error(c, http.StatusInternalServerError, constant.MsgFileUploadFail, err.Error())
			return
		}
		imagePath = "uploads/playlists/" + filename
	}

	playlist := entity.Playlist{
		UserID:      adminID,
		CreatorRole: "admin",
		Name:        input.Name,
		ImageURL:    imagePath,
	}

	if err := h.service.Create(&playlist, file); err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgPlaylistCreateFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgPlaylistCreateOK, h.toResponse(&playlist))
}

func (h *PlaylistHandler) GetAll(c *gin.Context) {
	playlists, err := h.service.Search("")
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgPlaylistListFail, err.Error())
		return
	}

	var result []response.PlaylistResponse
	for i := range playlists {
		result = append(result, h.toResponse(&playlists[i]))
	}

	resp.Success(c, http.StatusOK, constant.MsgPlaylistListOK, result)
}

func (h *PlaylistHandler) GetDetail(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	playlist, err := h.service.GetByID(id)
	if err != nil {
		resp.Error(c, http.StatusNotFound, constant.MsgPlaylistNotFound, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgPlaylistGetOK, h.toResponse(playlist))
}

func (h *PlaylistHandler) Update(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	var input request.UpdatePlaylistRequest
	if err := c.ShouldBind(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	adminID := security.GetUserID(c)
	updates := make(map[string]interface{})
	if input.Name != "" {
		updates["name"] = input.Name
	}

	file, _ := c.FormFile("image_file")
	if file != nil {
		pwd, _ := os.Getwd()
		uploadDir := filepath.Join(pwd, "uploads", "playlists")
		_ = os.MkdirAll(uploadDir, 0755)
		filename := fmt.Sprintf("admin_%d_%d_%s", adminID, time.Now().Unix(), filepath.Base(file.Filename))
		fullPath := filepath.Join(uploadDir, filename)
		if err := c.SaveUploadedFile(file, fullPath); err != nil {
			resp.Error(c, http.StatusInternalServerError, constant.MsgFileUploadFail, err.Error())
			return
		}
		updates["image_url"] = "uploads/playlists/" + filename
	}

	if len(updates) == 0 {
		resp.Error(c, http.StatusBadRequest, constant.MsgNoUpdatesProvided, nil)
		return
	}

	updated, err := h.service.Update(id, updates)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgPlaylistUpdateFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgPlaylistUpdateOK, h.toResponse(updated))
}

func (h *PlaylistHandler) Delete(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	if err := h.service.Delete(id, 0); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgPlaylistDeleteFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgPlaylistDeleteOK, nil)
}

func (h *PlaylistHandler) AddAudio(c *gin.Context) {
	var input request.PlaylistAudioRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	if err := h.service.AddAudioToPlaylist(input.PlaylistID, input.AudioID); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgPlaylistAddAudioFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgPlaylistAddAudioOK, nil)
}

func (h *PlaylistHandler) RemoveAudio(c *gin.Context) {
	var input request.PlaylistAudioRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	if err := h.service.RemoveAudioFromPlaylist(input.PlaylistID, input.AudioID); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgPlaylistRemoveAudioFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgPlaylistRemoveAudioOK, nil)
}

func (h *PlaylistHandler) toResponse(p *entity.Playlist) response.PlaylistResponse {
	r := response.PlaylistResponse{
		ID:            p.ID,
		UserID:        p.UserID,
		CreatorRole:   p.CreatorRole,
		Name:          p.Name,
		ImageURL:      p.ImageURL,
		DominantColor: p.DominantColor,
		ShareToken:    p.ShareToken,
		IsPublic:      p.IsPublic,
		TimeSince:     p.TimeSinceCreated(),
		AudioCount:    len(p.Audios),
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}

	if p.User != nil {
		r.CreatorName = p.User.Username
	}

	for _, a := range p.Audios {
		r.Audios = append(r.Audios, response.AudioResponse{
			ID:            a.ID,
			Title:         a.Title,
			Artist:        a.Artist,
			FilePath:      a.FilePath,
			Duration:      a.Duration,
			DurationFmt:   response.FormatDuration(a.Duration),
			FileSize:      a.FileSize,
			Status:        a.Status,
			CategoryID:    a.CategoryID,
			Thumbnail:     a.Thumbnail,
			DominantColor: a.DominantColor,
			CreatedAt:     a.CreatedAt,
			UpdatedAt:     a.UpdatedAt,
		})
	}

	return r
}
