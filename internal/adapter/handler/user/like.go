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

type LikeHandler struct {
	service port.LikeService
}

func NewLikeHandler(s port.LikeService) *LikeHandler {
	return &LikeHandler{service: s}
}

func (h *LikeHandler) Like(c *gin.Context) {
	var input request.LikeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	like, err := h.service.LikeAudio(userID, input)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgLikeFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgLikeOK, response.LikeResponse{
		ID:        like.ID,
		UserID:    like.UserID,
		AudioID:   like.AudioID,
		CreatedAt: like.CreatedAt,
	})
}

func (h *LikeHandler) Unlike(c *gin.Context) {
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

	if err := h.service.UnlikeAudio(userID, audioID); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgUnlikeFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgUnlikeOK, nil)
}

func (h *LikeHandler) GetLikes(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	likes, err := h.service.GetLikedAudios(userID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgLikeListFail, err.Error())
		return
	}

	var result []response.LikeResponse
	for _, l := range likes {
		result = append(result, response.LikeResponse{
			ID:        l.ID,
			UserID:    l.UserID,
			AudioID:   l.AudioID,
			CreatedAt: l.CreatedAt,
		})
	}

	resp.Success(c, http.StatusOK, constant.MsgLikeListOK, result)
}
