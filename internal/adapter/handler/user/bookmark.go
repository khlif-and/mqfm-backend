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

type BookmarkHandler struct {
	service port.BookmarkService
}

func NewBookmarkHandler(s port.BookmarkService) *BookmarkHandler {
	return &BookmarkHandler{service: s}
}

func (h *BookmarkHandler) Create(c *gin.Context) {
	var input request.CreateBookmarkRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	bookmark, err := h.service.Create(userID, input)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgBookmarkCreateFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgBookmarkCreateOK, response.BookmarkResponse{
		ID:              bookmark.ID,
		UserID:          bookmark.UserID,
		AudioID:         bookmark.AudioID,
		PositionSeconds: bookmark.PositionSeconds,
		Label:           bookmark.Label,
		CreatedAt:       bookmark.CreatedAt,
	})
}

func (h *BookmarkHandler) GetByUser(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	bookmarks, err := h.service.GetByUser(userID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgBookmarkListFail, err.Error())
		return
	}

	var result []response.BookmarkResponse
	for _, b := range bookmarks {
		result = append(result, response.BookmarkResponse{
			ID:              b.ID,
			UserID:          b.UserID,
			AudioID:         b.AudioID,
			PositionSeconds: b.PositionSeconds,
			Label:           b.Label,
			CreatedAt:       b.CreatedAt,
		})
	}

	resp.Success(c, http.StatusOK, constant.MsgBookmarkListOK, result)
}

func (h *BookmarkHandler) GetByAudio(c *gin.Context) {
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

	bookmarks, err := h.service.GetByUserAndAudio(userID, audioID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgBookmarkListFail, err.Error())
		return
	}

	var result []response.BookmarkResponse
	for _, b := range bookmarks {
		result = append(result, response.BookmarkResponse{
			ID:              b.ID,
			UserID:          b.UserID,
			AudioID:         b.AudioID,
			PositionSeconds: b.PositionSeconds,
			Label:           b.Label,
			CreatedAt:       b.CreatedAt,
		})
	}

	resp.Success(c, http.StatusOK, constant.MsgBookmarkListOK, result)
}

func (h *BookmarkHandler) Delete(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	if err := h.service.Delete(id, userID); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgBookmarkDeleteFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgBookmarkDeleteOK, nil)
}
