package request

type HistoryRequest struct {
	AudioID uint `json:"audio_id" binding:"required"`
}
