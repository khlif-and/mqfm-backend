package response

import "time"

type BookmarkResponse struct {
	ID              uint          `json:"id"`
	UserID          uint          `json:"user_id"`
	AudioID         uint          `json:"audio_id"`
	Audio           *AudioResponse `json:"audio,omitempty"`
	PositionSeconds int           `json:"position_seconds"`
	Label           string        `json:"label"`
	CreatedAt       time.Time     `json:"created_at"`
}
