package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/internal/shared/helper"
	resp "mqfm-backend/internal/shared/response"
	"mqfm-backend/internal/shared/security"
)

type PlaylistCollabHandler struct {
	service port.PlaylistCollabService
}

func NewPlaylistCollabHandler(s port.PlaylistCollabService) *PlaylistCollabHandler {
	return &PlaylistCollabHandler{service: s}
}

func (h *PlaylistCollabHandler) AddCollaborator(c *gin.Context) {
	var input request.AddCollaboratorRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	if err := h.service.AddCollaborator(userID, input.PlaylistID, input.CollaboratorID); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgCollabAddFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgCollabAddOK, nil)
}

func (h *PlaylistCollabHandler) RemoveCollaborator(c *gin.Context) {
	playlistID, ok := helper.ParamToUint(c, "playlist_id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}
	collaboratorID, ok := helper.ParamToUint(c, "user_id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	if err := h.service.RemoveCollaborator(userID, playlistID, collaboratorID); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgCollabRemoveFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgCollabRemoveOK, nil)
}

func (h *PlaylistCollabHandler) GetCollaborators(c *gin.Context) {
	playlistID, ok := helper.ParamToUint(c, "playlist_id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	collabs, err := h.service.GetCollaborators(playlistID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgCollabListFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgCollabListOK, collabs)
}

func (h *PlaylistCollabHandler) ContributeAudio(c *gin.Context) {
	var input request.ContributeAudioRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	if err := h.service.ContributeAudio(userID, input.PlaylistID, input.AudioID); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgCollabContributeFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgCollabContributeOK, nil)
}

func (h *PlaylistCollabHandler) JoinByToken(c *gin.Context) {
	var input request.JoinPlaylistRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	if err := h.service.JoinByShareToken(userID, input.ShareToken); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgCollabJoinFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgCollabJoinOK, nil)
}
