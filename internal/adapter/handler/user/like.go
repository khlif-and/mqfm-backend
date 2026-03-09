package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/internal/shared/dto/response"
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

	like, err := h.service.Like(userID, input)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgLikeFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgLikeOK, response.LikeResponse{
		ID:         like.ID,
		UserID:     like.UserID,
		TargetType: like.TargetType,
		TargetID:   like.TargetID,
		CreatedAt:  like.CreatedAt,
	})
}

func (h *LikeHandler) Unlike(c *gin.Context) {
	var input request.UnlikeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	if err := h.service.Unlike(userID, input); err != nil {
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

	targetType := c.DefaultQuery("type", "audio")
	if targetType != "audio" && targetType != "playlist" {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, "type must be audio or playlist")
		return
	}

	likes, err := h.service.GetLikes(fmt.Sprintf("%d", userID), targetType)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgLikeListFail, err.Error())
		return
	}

	var result []response.LikeResponse
	for _, l := range likes {
		lr := response.LikeResponse{
			ID:         l.ID,
			UserID:     l.UserID,
			TargetType: l.TargetType,
			TargetID:   l.TargetID,
			CreatedAt:  l.CreatedAt,
		}
		if l.TargetType == "audio" && l.Audio != nil {
			lr.Audio = response.AudioResponse{
				ID:            l.Audio.ID,
				Title:         l.Audio.Title,
				Artist:        l.Audio.Artist,
				FilePath:      l.Audio.FilePath,
				Duration:      l.Audio.Duration,
				Status:        l.Audio.Status,
				CategoryID:    l.Audio.CategoryID,
				Thumbnail:     l.Audio.Thumbnail,
				DominantColor: l.Audio.DominantColor,
				CreatedAt:     l.Audio.CreatedAt,
				UpdatedAt:     l.Audio.UpdatedAt,
			}
		}
		if l.TargetType == "playlist" && l.Playlist != nil {
			lr.Playlist = response.PlaylistResponse{
				ID:            l.Playlist.ID,
				UserID:        l.Playlist.UserID,
				CreatorRole:   l.Playlist.CreatorRole,
				Name:          l.Playlist.Name,
				ImageURL:      l.Playlist.ImageURL,
				DominantColor: l.Playlist.DominantColor,
				IsPublic:      l.Playlist.IsPublic,
				TimeSince:     l.Playlist.TimeSinceCreated(),
				AudioCount:    len(l.Playlist.Audios),
				CreatedAt:     l.Playlist.CreatedAt,
				UpdatedAt:     l.Playlist.UpdatedAt,
			}
			if l.Playlist.User != nil {
				lr.Playlist = response.PlaylistResponse{
					ID:            l.Playlist.ID,
					UserID:        l.Playlist.UserID,
					CreatorRole:   l.Playlist.CreatorRole,
					Name:          l.Playlist.Name,
					ImageURL:      l.Playlist.ImageURL,
					DominantColor: l.Playlist.DominantColor,
					IsPublic:      l.Playlist.IsPublic,
					TimeSince:     l.Playlist.TimeSinceCreated(),
					CreatorName:   l.Playlist.User.Username,
					AudioCount:    len(l.Playlist.Audios),
					CreatedAt:     l.Playlist.CreatedAt,
					UpdatedAt:     l.Playlist.UpdatedAt,
				}
			}
		}
		result = append(result, lr)
	}

	resp.Success(c, http.StatusOK, constant.MsgLikeListOK, result)
}
