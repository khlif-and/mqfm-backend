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

type ProgressHandler struct {
	service port.AudioProgressService
}

func NewProgressHandler(s port.AudioProgressService) *ProgressHandler {
	return &ProgressHandler{service: s}
}

func (h *ProgressHandler) Update(c *gin.Context) {
	var input request.UpdateProgressRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	progress, err := h.service.UpdateProgress(userID, input)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgProgressUpdateFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgProgressUpdateOK, toProgressResponse(progress))
}

func (h *ProgressHandler) Get(c *gin.Context) {
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

	progress, err := h.service.GetProgress(userID, audioID)
	if err != nil {
		resp.Error(c, http.StatusNotFound, constant.MsgProgressNotFound, nil)
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgProgressGetOK, toProgressResponse(progress))
}

func (h *ProgressHandler) GetAll(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	progresses, err := h.service.GetAllProgress(userID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgProgressListFail, err.Error())
		return
	}

	var result []response.ProgressResponse
	for _, p := range progresses {
		result = append(result, *toProgressResponse(&p))
	}

	resp.Success(c, http.StatusOK, constant.MsgProgressListOK, result)
}

func (h *ProgressHandler) GetCompleted(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	progresses, err := h.service.GetCompleted(userID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgProgressListFail, err.Error())
		return
	}

	var result []response.ProgressResponse
	for _, p := range progresses {
		result = append(result, *toProgressResponse(&p))
	}

	resp.Success(c, http.StatusOK, constant.MsgProgressListOK, result)
}
