package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/dto/request"
	resp "mqfm-backend/internal/shared/response"
	"mqfm-backend/internal/shared/security"
)

type LocationHandler struct {
	service port.UserLocationService
}

func NewLocationHandler(s port.UserLocationService) *LocationHandler {
	return &LocationHandler{service: s}
}

func (h *LocationHandler) Update(c *gin.Context) {
	var input request.UpdateLocationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	loc, err := h.service.Update(userID, input)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgLocationUpdateFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgLocationUpdateOK, loc)
}

func (h *LocationHandler) Get(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	loc, err := h.service.Get(userID)
	if err != nil {
		resp.Error(c, http.StatusNotFound, constant.MsgLocationNotFound, nil)
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgLocationGetOK, loc)
}
