package response

import "time"

type RankingResponse struct {
	Rank         int            `json:"rank"`
	AudioID      uint           `json:"audio_id"`
	Audio        *AudioResponse `json:"audio,omitempty"`
	WeeklyVotes  int64          `json:"weekly_votes"`
	MonthlyVotes int64          `json:"monthly_votes"`
	TotalVotes   int64          `json:"total_votes"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type VoteStatusResponse struct {
	AudioID   uint  `json:"audio_id"`
	HasVoted  bool  `json:"has_voted"`
	TotalVotes int64 `json:"total_votes"`
}
