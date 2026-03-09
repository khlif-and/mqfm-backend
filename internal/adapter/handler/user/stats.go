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

type StatsHandler struct {
	service port.ListeningStatService
}

func NewStatsHandler(s port.ListeningStatService) *StatsHandler {
	return &StatsHandler{service: s}
}

func (h *StatsHandler) RecordStat(c *gin.Context) {
	var input request.RecordStatRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	if err := h.service.RecordStat(userID, input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgStatRecordFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgStatRecordOK, nil)
}

func (h *StatsHandler) GetRecap(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	recap, err := h.service.GetRecap(userID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgStatRecapFail, err.Error())
		return
	}

	var cats []response.CategoryStatResponse
	for _, c := range recap.TopCategories {
		cats = append(cats, response.CategoryStatResponse{
			CategoryID: c.CategoryID,
			Name:       c.Name,
			TotalTime:  c.TotalTime,
		})
	}

	var artists []response.ArtistStatResponse
	for _, a := range recap.TopArtists {
		artists = append(artists, response.ArtistStatResponse{
			Artist:    a.Artist,
			TotalTime: a.TotalTime,
		})
	}

	var daily []response.DailyStatResponse
	for _, d := range recap.DailyStats {
		daily = append(daily, response.DailyStatResponse{
			Date:      d.Date,
			TotalTime: d.TotalTime,
		})
	}

	resp.Success(c, http.StatusOK, constant.MsgStatRecapOK, response.StatsRecapResponse{
		WeeklyMinutes:  recap.WeeklyMinutes,
		MonthlyMinutes: recap.MonthlyMinutes,
		TopCategories:  cats,
		TopArtists:     artists,
		DailyStats:     daily,
	})
}
