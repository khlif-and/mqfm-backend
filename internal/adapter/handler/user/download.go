package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

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

	resp.Success(c, http.StatusCreated, constant.MsgDownloadOK, response.DownloadResponse{
		ID:        dl.ID,
		UserID:    dl.UserID,
		AudioID:   dl.AudioID,
		FileSize:  dl.FileSize,
		CreatedAt: dl.CreatedAt,
	})
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

	var result []response.DownloadResponse
	for _, d := range downloads {
		result = append(result, response.DownloadResponse{
			ID:        d.ID,
			UserID:    d.UserID,
			AudioID:   d.AudioID,
			FileSize:  d.FileSize,
			CreatedAt: d.CreatedAt,
		})
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
