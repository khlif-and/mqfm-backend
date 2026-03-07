package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/dto/response"
	"mqfm-backend/internal/shared/helper"
	resp "mqfm-backend/internal/shared/response"
	"mqfm-backend/internal/shared/security"
)

type HistoryHandler struct {
	service port.HistoryService
}

func NewHistoryHandler(s port.HistoryService) *HistoryHandler {
	return &HistoryHandler{service: s}
}

func (h *HistoryHandler) GetHistory(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	histories, err := h.service.GetHistory(userID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgHistoryListFail, err.Error())
		return
	}

	var result []response.HistoryResponse
	for _, h := range histories {
		result = append(result, response.HistoryResponse{
			ID:        h.ID,
			UserID:    h.UserID,
			AudioID:   h.AudioID,
			PlayCount: h.PlayCount,
			PlayedAt:  h.PlayedAt,
			CreatedAt: h.CreatedAt,
		})
	}

	resp.Success(c, http.StatusOK, constant.MsgHistoryListOK, result)
}

func (h *HistoryHandler) DeleteHistory(c *gin.Context) {
	audioID, ok := helper.ParamToUint(c, "audio_id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidAudioID, nil)
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	if err := h.service.DeleteHistory(userID, audioID); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgHistoryDeleteFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgHistoryDeleteOK, nil)
}

func (h *HistoryHandler) ClearHistory(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	if err := h.service.ClearHistory(userID); err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgHistoryClearFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgHistoryClearOK, nil)
}
