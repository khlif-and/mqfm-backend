package response

import "time"

type PlaylistResponse struct {
	ID            uint            `json:"id"`
	UserID        uint            `json:"user_id"`
	CreatorRole   string          `json:"creator_role"`
	Name          string          `json:"name"`
	ImageURL      string          `json:"image_url,omitempty"`
	DominantColor string          `json:"dominant_color,omitempty"`
	ShareToken    string          `json:"share_token,omitempty"`
	IsPublic      bool            `json:"is_public"`
	TimeSince     string          `json:"time_since_created"`
	CreatorName   string          `json:"creator_name"`
	AudioCount    int             `json:"audio_count"`
	Audios        []AudioResponse `json:"audios,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}
