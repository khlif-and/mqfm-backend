package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/internal/shared/dto/response"
	"mqfm-backend/internal/shared/helper"
	resp "mqfm-backend/internal/shared/response"
	"mqfm-backend/internal/shared/security"
)

type VoteHandler struct {
	service port.AudioVoteService
}

func NewVoteHandler(s port.AudioVoteService) *VoteHandler {
	return &VoteHandler{service: s}
}

func (h *VoteHandler) Vote(c *gin.Context) {
	var input request.VoteRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	if err := h.service.Vote(userID, input.AudioID); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgVoteFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgVoteOK, nil)
}

func (h *VoteHandler) Unvote(c *gin.Context) {
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

	if err := h.service.Unvote(userID, audioID); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgUnvoteFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgUnvoteOK, nil)
}

func (h *VoteHandler) WeeklyRanking(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	rankings, err := h.service.GetWeeklyRanking(limit)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgRankingFail, err.Error())
		return
	}

	var result []response.RankingResponse
	for _, r := range rankings {
		item := response.RankingResponse{
			Rank:         r.WeeklyRank,
			AudioID:      r.AudioID,
			WeeklyVotes:  r.WeeklyVotes,
			MonthlyVotes: r.MonthlyVotes,
			TotalVotes:   r.TotalVotes,
			UpdatedAt:    r.UpdatedAt,
		}
		if r.Audio != nil {
			ar := toAudioResponseVal(*r.Audio)
			item.Audio = &ar
		}
		result = append(result, item)
	}

	resp.Success(c, http.StatusOK, constant.MsgRankingWeeklyOK, result)
}

func (h *VoteHandler) MonthlyRanking(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	rankings, err := h.service.GetMonthlyRanking(limit)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgRankingFail, err.Error())
		return
	}

	var result []response.RankingResponse
	for _, r := range rankings {
		item := response.RankingResponse{
			Rank:         r.MonthlyRank,
			AudioID:      r.AudioID,
			WeeklyVotes:  r.WeeklyVotes,
			MonthlyVotes: r.MonthlyVotes,
			TotalVotes:   r.TotalVotes,
			UpdatedAt:    r.UpdatedAt,
		}
		if r.Audio != nil {
			ar := toAudioResponseVal(*r.Audio)
			item.Audio = &ar
		}
		result = append(result, item)
	}

	resp.Success(c, http.StatusOK, constant.MsgRankingMonthlyOK, result)
}

func (h *VoteHandler) Status(c *gin.Context) {
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

	hasVoted, _ := h.service.HasVoted(userID, audioID)
	rankings, _ := h.service.GetWeeklyRanking(1)
	var total int64
	if len(rankings) > 0 {
		total = rankings[0].TotalVotes
	}

	resp.Success(c, http.StatusOK, constant.MsgVoteStatusOK, response.VoteStatusResponse{
		AudioID:    audioID,
		HasVoted:   hasVoted,
		TotalVotes: total,
	})
}
