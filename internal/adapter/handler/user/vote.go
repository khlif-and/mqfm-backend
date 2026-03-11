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
	scoreRepo port.AudioScoreRepository
}

func NewVoteHandler(s port.AudioVoteService, sr port.AudioScoreRepository) *VoteHandler {
	return &VoteHandler{service: s, scoreRepo: sr}
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
	if limit > 20 {
		limit = 20
	}

	scores, err := h.scoreRepo.FindTopByWeeklyLikes(limit)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgRankingFail, err.Error())
		return
	}

	var result []response.RankingResponse
	for i, s := range scores {
		item := response.RankingResponse{
			Rank:      i + 1,
			AudioID:   s.AudioID,
			Likes:     s.TotalLikes,
			UpdatedAt: s.UpdatedAt,
		}
		if s.Audio != nil {
			ar := toAudioResponseVal(*s.Audio)
			item.Audio = &ar
		}
		result = append(result, item)
	}

	resp.Success(c, http.StatusOK, constant.MsgRankingWeeklyOK, result)
}

func (h *VoteHandler) MonthlyRanking(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit > 20 {
		limit = 20
	}

	scores, err := h.scoreRepo.FindTopByLikes(limit, 35000)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgRankingFail, err.Error())
		return
	}

	var result []response.RankingResponse
	for i, s := range scores {
		item := response.RankingResponse{
			Rank:      i + 1,
			AudioID:   s.AudioID,
			Likes:     s.TotalLikes,
			UpdatedAt: s.UpdatedAt,
		}
		if s.Audio != nil {
			ar := toAudioResponseVal(*s.Audio)
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
