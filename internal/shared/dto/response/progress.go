package response

import "time"

type ProgressResponse struct {
	ID           uint           `json:"id"`
	UserID       uint           `json:"user_id"`
	AudioID      uint           `json:"audio_id"`
	Audio        *AudioResponse `json:"audio,omitempty"`
	LastPosition int            `json:"last_position"`
	Duration     int            `json:"duration"`
	Percentage   float64        `json:"percentage"`
	Completed    bool           `json:"completed"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
