package response

import "time"

type RankingResponse struct {
	Rank      int            `json:"rank"`
	AudioID   uint           `json:"audio_id"`
	Audio     *AudioResponse `json:"audio,omitempty"`
	Likes     int64          `json:"likes"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type VoteStatusResponse struct {
	AudioID   uint  `json:"audio_id"`
	HasVoted  bool  `json:"has_voted"`
	TotalVotes int64 `json:"total_votes"`
}
