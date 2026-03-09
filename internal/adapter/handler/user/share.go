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

type ShareHandler struct {
	service     port.ShareService
	playlistSvc port.PlaylistService
}

func NewShareHandler(s port.ShareService, playlistSvc port.PlaylistService) *ShareHandler {
	return &ShareHandler{service: s, playlistSvc: playlistSvc}
}

func (h *ShareHandler) ShareAudio(c *gin.Context) {
	audioID, ok := helper.ParamToUint(c, "audio_id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidAudioID, nil)
		return
	}

	link := h.service.GenerateAudioShareLink(audioID)
	resp.Success(c, http.StatusOK, constant.MsgShareOK, response.ShareLinkResponse{
		ShareLink: link,
		Type:      "audio",
	})
}

func (h *ShareHandler) SharePlaylist(c *gin.Context) {
	playlistID, ok := helper.ParamToUint(c, "playlist_id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	token, err := h.playlistSvc.SharePlaylist(playlistID)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgShareFail, err.Error())
		return
	}

	link := h.service.GeneratePlaylistShareLink(token)
	resp.Success(c, http.StatusOK, constant.MsgShareOK, response.ShareLinkResponse{
		ShareLink: link,
		Type:      "playlist",
	})
}

func (h *ShareHandler) GetSharedPlaylist(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, nil)
		return
	}

	playlist, err := h.playlistSvc.GetByShareToken(token)
	if err != nil {
		resp.Error(c, http.StatusNotFound, constant.MsgPlaylistNotFound, nil)
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgPlaylistGetOK, playlist)
}
