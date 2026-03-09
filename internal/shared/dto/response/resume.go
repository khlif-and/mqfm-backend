package response

import "time"

type ResumeResponse struct {
	AudioID         uint           `json:"audio_id"`
	Audio           *AudioResponse `json:"audio,omitempty"`
	PlaylistID      uint           `json:"playlist_id"`
	PositionSeconds int            `json:"position_seconds"`
	UpdatedAt       time.Time      `json:"updated_at"`
}
