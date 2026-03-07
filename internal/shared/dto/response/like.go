package response

import "time"

type LikeResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	AudioID   uint      `json:"audio_id"`
	CreatedAt time.Time `json:"created_at"`
}
