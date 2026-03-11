package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/internal/shared/dto/response"
	"mqfm-backend/internal/shared/helper"
	resp "mqfm-backend/internal/shared/response"
	"mqfm-backend/internal/shared/security"
)

type DownloadHandler struct {
	service port.DownloadService
}

func NewDownloadHandler(s port.DownloadService) *DownloadHandler {
	return &DownloadHandler{service: s}
}

func (h *DownloadHandler) Record(c *gin.Context) {
	var input request.DownloadRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	dl, err := h.service.RecordDownload(userID, input)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgDownloadFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgDownloadOK, toDownloadResponse(dl))
}

func (h *DownloadHandler) GetAll(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	downloads, err := h.service.GetDownloads(userID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgDownloadListFail, err.Error())
		return
	}

	result := make([]response.DownloadResponse, 0, len(downloads))
	for i := range downloads {
		result = append(result, toDownloadResponse(&downloads[i]))
	}

	resp.Success(c, http.StatusOK, constant.MsgDownloadListOK, result)
}

func (h *DownloadHandler) Delete(c *gin.Context) {
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

	if err := h.service.DeleteDownload(id, userID); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgDownloadDeleteFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgDownloadDeleteOK, nil)
}

func (h *DownloadHandler) StorageUsage(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	size, err := h.service.GetStorageUsage(userID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgInternalError, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgDownloadStorageOK, response.StorageUsageResponse{
		TotalBytes: size,
		TotalMB:    size / (1024 * 1024),
	})
}

func (h *DownloadHandler) SmartDownload(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	audios, err := h.service.GetNewFromFavorites(userID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgDownloadSmartFail, err.Error())
		return
	}

	var result []response.AudioResponse
	for _, a := range audios {
		result = append(result, toAudioResponseVal(a))
	}

	resp.Success(c, http.StatusOK, constant.MsgDownloadSmartOK, result)
}

func toDownloadResponse(d *entity.Download) response.DownloadResponse {
	r := response.DownloadResponse{
		ID:         d.ID,
		UserID:     d.UserID,
		AudioID:    d.AudioID,
		PlaylistID: d.PlaylistID,
		FileSize:   d.FileSize,
		ExpiresAt:  d.ExpiresAt,
		CreatedAt:  d.CreatedAt,
	}
	if !d.ExpiresAt.IsZero() {
		remaining := int(time.Until(d.ExpiresAt).Hours() / 24)
		if remaining < 0 {
			remaining = 0
		}
		r.DaysRemaining = remaining
	}
	if d.Audio != nil {
		r.Title = d.Audio.Title
		r.Artist = d.Audio.Artist
		r.Thumbnail = d.Audio.Thumbnail
		r.DominantColor = d.Audio.DominantColor
		r.Duration = d.Audio.Duration
		r.DurationFmt = response.FormatDuration(d.Audio.Duration)
		ar := toAudioResponseVal(*d.Audio)
		r.Audio = &ar
	}
	return r
}
