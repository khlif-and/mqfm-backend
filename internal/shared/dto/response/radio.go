package response

import "time"

type RadioResponse struct {
	ID            uint            `json:"id"`
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	Thumbnail     string          `json:"thumbnail,omitempty"`
	DominantColor string          `json:"dominant_color,omitempty"`
	IsActive      bool            `json:"is_active"`
	AudioCount    int             `json:"audio_count"`
	Audios        []AudioResponse `json:"audios,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}
