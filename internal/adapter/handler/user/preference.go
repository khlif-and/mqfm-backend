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

type PreferenceHandler struct {
	service port.UserPreferenceService
}

func NewPreferenceHandler(s port.UserPreferenceService) *PreferenceHandler {
	return &PreferenceHandler{service: s}
}

func (h *PreferenceHandler) Get(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	pref, err := h.service.GetOrCreate(userID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgInternalError, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgPreferenceGetOK, response.PreferenceResponse{
		PlaybackSpeed:     pref.PlaybackSpeed,
		SleepTimerMinutes: pref.SleepTimerMinutes,
		AutoDownloadWifi:  pref.AutoDownloadWifi,
	})
}

func (h *PreferenceHandler) Update(c *gin.Context) {
	var input request.UpdatePreferenceRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	pref, err := h.service.Update(userID, input)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgPreferenceUpdateFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgPreferenceUpdateOK, response.PreferenceResponse{
		PlaybackSpeed:     pref.PlaybackSpeed,
		SleepTimerMinutes: pref.SleepTimerMinutes,
		AutoDownloadWifi:  pref.AutoDownloadWifi,
	})
}
