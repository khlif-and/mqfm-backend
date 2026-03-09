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

type ClipHandler struct {
	service  port.AudioClipService
	shareSvc port.ShareService
}

func NewClipHandler(s port.AudioClipService, shareSvc port.ShareService) *ClipHandler {
	return &ClipHandler{service: s, shareSvc: shareSvc}
}

func (h *ClipHandler) Create(c *gin.Context) {
	var input request.CreateClipRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	clip, err := h.service.CreateClip(userID, input)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgClipCreateFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgClipCreateOK, response.ClipResponse{
		ID:         clip.ID,
		UserID:     clip.UserID,
		AudioID:    clip.AudioID,
		StartTime:  clip.StartTime,
		EndTime:    clip.EndTime,
		ClipPath:   clip.ClipPath,
		ShareToken: clip.ShareToken,
		ShareLink:  h.shareSvc.GenerateClipShareLink(clip.ShareToken),
		CreatedAt:  clip.CreatedAt,
	})
}

func (h *ClipHandler) GetByUser(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	clips, err := h.service.GetByUser(userID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgClipListFail, err.Error())
		return
	}

	var result []response.ClipResponse
	for _, cl := range clips {
		result = append(result, response.ClipResponse{
			ID:         cl.ID,
			UserID:     cl.UserID,
			AudioID:    cl.AudioID,
			StartTime:  cl.StartTime,
			EndTime:    cl.EndTime,
			ClipPath:   cl.ClipPath,
			ShareToken: cl.ShareToken,
			ShareLink:  h.shareSvc.GenerateClipShareLink(cl.ShareToken),
			CreatedAt:  cl.CreatedAt,
		})
	}

	resp.Success(c, http.StatusOK, constant.MsgClipListOK, result)
}

func (h *ClipHandler) GetByShareToken(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, nil)
		return
	}

	clip, err := h.service.GetByShareToken(token)
	if err != nil {
		resp.Error(c, http.StatusNotFound, constant.MsgClipNotFound, nil)
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgClipGetOK, response.ClipResponse{
		ID:         clip.ID,
		UserID:     clip.UserID,
		AudioID:    clip.AudioID,
		StartTime:  clip.StartTime,
		EndTime:    clip.EndTime,
		ClipPath:   clip.ClipPath,
		ShareToken: clip.ShareToken,
		ShareLink:  h.shareSvc.GenerateClipShareLink(clip.ShareToken),
		CreatedAt:  clip.CreatedAt,
	})
}

func (h *ClipHandler) Delete(c *gin.Context) {
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
		resp.Error(c, http.StatusBadRequest, constant.MsgClipDeleteFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgClipDeleteOK, nil)
}
