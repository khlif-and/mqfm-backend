package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	resp "mqfm-backend/internal/shared/response"
	"mqfm-backend/internal/shared/security"
)

type FavoriteArtistHandler struct {
	service port.FavoriteArtistService
}

func NewFavoriteArtistHandler(s port.FavoriteArtistService) *FavoriteArtistHandler {
	return &FavoriteArtistHandler{service: s}
}

func (h *FavoriteArtistHandler) Add(c *gin.Context) {
	var input struct {
		ArtistName string `json:"artist_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	if err := h.service.Add(userID, input.ArtistName); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgFavoriteArtistAddFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgFavoriteArtistAddOK, nil)
}

func (h *FavoriteArtistHandler) Remove(c *gin.Context) {
	artistName := c.Param("artist")
	if artistName == "" {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, nil)
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	if err := h.service.Remove(userID, artistName); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgFavoriteArtistRemoveFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgFavoriteArtistRemoveOK, nil)
}

func (h *FavoriteArtistHandler) GetAll(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	favs, err := h.service.GetByUser(userID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgFavoriteArtistListFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgFavoriteArtistListOK, favs)
}
