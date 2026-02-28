package likes

import "time"

type LikeRequest struct {
	AudioID uint `json:"audio_id" binding:"required"`
}

type LikeResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	AudioID   uint      `json:"audio_id"`
	CreatedAt time.Time `json:"created_at"`
}
