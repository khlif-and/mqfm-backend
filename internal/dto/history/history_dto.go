package history

import "time"

type HistoryRequest struct {
	AudioID uint `json:"audio_id" binding:"required"`
}

type HistoryResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	AudioID   uint      `json:"audio_id"`
	PlayCount int       `json:"play_count"`
	PlayedAt  time.Time `json:"played_at"`
	CreatedAt time.Time `json:"created_at"`
}
