package request

type VoteRequest struct {
	AudioID uint `json:"audio_id" binding:"required"`
}
