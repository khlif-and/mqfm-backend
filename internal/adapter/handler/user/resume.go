package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/internal/shared/dto/response"
	resp "mqfm-backend/internal/shared/response"
	"mqfm-backend/internal/shared/security"
)

type ResumeHandler struct {
	service port.SmartResumeService
}

func NewResumeHandler(s port.SmartResumeService) *ResumeHandler {
	return &ResumeHandler{service: s}
}

func (h *ResumeHandler) Update(c *gin.Context) {
	var input request.UpdateResumeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	resume, err := h.service.Update(userID, input)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgResumeUpdateFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgResumeUpdateOK, response.ResumeResponse{
		AudioID:         resume.AudioID,
		Audio:           toAudioResponsePtr(resume.Audio),
		PlaylistID:      resume.PlaylistID,
		PositionSeconds: resume.PositionSeconds,
		UpdatedAt:       resume.UpdatedAt,
	})
}

func (h *ResumeHandler) Get(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	resume, err := h.service.Get(userID)
	if err != nil {
		resp.Error(c, http.StatusNotFound, constant.MsgResumeNotFound, nil)
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgResumeGetOK, response.ResumeResponse{
		AudioID:         resume.AudioID,
		Audio:           toAudioResponsePtr(resume.Audio),
		PlaylistID:      resume.PlaylistID,
		PositionSeconds: resume.PositionSeconds,
		UpdatedAt:       resume.UpdatedAt,
	})
}
