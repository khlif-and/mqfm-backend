package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/dto/response"
	"mqfm-backend/internal/shared/helper"
	resp "mqfm-backend/internal/shared/response"
	"mqfm-backend/internal/shared/security"
)

type RecommendationHandler struct {
	service port.RecommendationService
}

func NewRecommendationHandler(svc port.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{service: svc}
}

func (h *RecommendationHandler) GetPopular(c *gin.Context) {
	audios, err := h.service.GetPopular(20)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgRecommendFail, err.Error())
		return
	}
	resp.Success(c, http.StatusOK, constant.MsgRecommendPopularOK, toAudioList(audios))
}

func (h *RecommendationHandler) GetMostListened(c *gin.Context) {
	audios, err := h.service.GetMostListened(20)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgRecommendFail, err.Error())
		return
	}
	resp.Success(c, http.StatusOK, constant.MsgRecommendMostListenedOK, toAudioList(audios))
}

func (h *RecommendationHandler) GetByArtist(c *gin.Context) {
	artist := c.Query("artist")
	if artist == "" {
		resp.Error(c, http.StatusBadRequest, constant.MsgArtistRequired, nil)
		return
	}
	audios, err := h.service.GetByArtist(artist, 20)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgRecommendFail, err.Error())
		return
	}
	resp.Success(c, http.StatusOK, constant.MsgRecommendByArtistOK, toAudioList(audios))
}

func (h *RecommendationHandler) GetSimilar(c *gin.Context) {
	audioID, ok := helper.ParamToUint(c, "audio_id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidAudioID, nil)
		return
	}
	audios, err := h.service.GetSimilar(audioID, 20)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgRecommendFail, err.Error())
		return
	}
	resp.Success(c, http.StatusOK, constant.MsgRecommendSimilarOK, toAudioList(audios))
}

func (h *RecommendationHandler) GetQuickPick(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}
	audios, err := h.service.GetQuickPick(userID, 20)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgRecommendFail, err.Error())
		return
	}
	resp.Success(c, http.StatusOK, constant.MsgRecommendQuickPickOK, toAudioList(audios))
}

func (h *RecommendationHandler) GetOnboarding(c *gin.Context) {
	audios, err := h.service.GetOnboarding(20)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgRecommendFail, err.Error())
		return
	}
	resp.Success(c, http.StatusOK, constant.MsgRecommendOnboardingOK, toAudioList(audios))
}

func (h *RecommendationHandler) GetPersonalized(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}
	audios, err := h.service.GetPersonalized(userID, 20)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgRecommendFail, err.Error())
		return
	}
	resp.Success(c, http.StatusOK, constant.MsgRecommendPersonalizedOK, toAudioList(audios))
}

func (h *RecommendationHandler) GetLocationBased(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}
	audios, err := h.service.GetLocationBased(userID, 20)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgRecommendFail, err.Error())
		return
	}
	resp.Success(c, http.StatusOK, constant.MsgRecommendLocationOK, toAudioList(audios))
}

func (h *RecommendationHandler) GetTimeBased(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}
	audios, err := h.service.GetTimeBasedPersonalized(userID, time.Now().Hour(), 20)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgRecommendFail, err.Error())
		return
	}
	resp.Success(c, http.StatusOK, constant.MsgRecommendTimeBasedOK, toAudioList(audios))
}

func toAudioList(audios []entity.Audio) []response.AudioResponse {
	var result []response.AudioResponse
	for _, a := range audios {
		result = append(result, response.AudioResponse{
			ID:            a.ID,
			Title:         a.Title,
			Artist:        a.Artist,
			Description:   a.Description,
			FilePath:      a.FilePath,
			Duration:      a.Duration,
			Status:        a.Status,
			CategoryID:    a.CategoryID,
			Thumbnail:     a.Thumbnail,
			DominantColor: a.DominantColor,
			CreatedAt:     a.CreatedAt,
			UpdatedAt:     a.UpdatedAt,
		})
	}
	return result
}
