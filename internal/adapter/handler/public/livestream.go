package public

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	resp "mqfm-backend/internal/shared/response"
)

type LivestreamHandler struct {
	service port.LivestreamService
}

func NewLivestreamHandler(s port.LivestreamService) *LivestreamHandler {
	return &LivestreamHandler{service: s}
}

func (h *LivestreamHandler) GetStatus(c *gin.Context) {
	status, err := h.service.GetStatus()
	if err != nil {
		resp.Success(c, http.StatusOK, constant.MsgLivestreamNoData, gin.H{"is_live": false})
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgLivestreamOK, status)
}
