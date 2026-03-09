package response

import "time"

type ClipResponse struct {
	ID         uint           `json:"id"`
	UserID     uint           `json:"user_id"`
	AudioID    uint           `json:"audio_id"`
	Audio      *AudioResponse `json:"audio,omitempty"`
	StartTime  int            `json:"start_time"`
	EndTime    int            `json:"end_time"`
	ClipPath   string         `json:"clip_path"`
	ShareToken string         `json:"share_token"`
	ShareLink  string         `json:"share_link"`
	CreatedAt  time.Time      `json:"created_at"`
}
