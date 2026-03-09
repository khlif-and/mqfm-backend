package response

import "time"

type LikeResponse struct {
	ID         uint        `json:"id"`
	UserID     uint        `json:"user_id"`
	TargetType string      `json:"target_type"`
	TargetID   uint        `json:"target_id"`
	Audio      interface{} `json:"audio,omitempty"`
	Playlist   interface{} `json:"playlist,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
}
