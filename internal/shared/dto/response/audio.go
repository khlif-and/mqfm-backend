package response

import "time"

type AudioResponse struct {
	ID            uint      `json:"audio_id"`
	Title         string    `json:"title"`
	Artist        string    `json:"artist"`
	Description   string    `json:"description,omitempty"`
	FilePath      string    `json:"file_path"`
	Duration      int       `json:"duration"`
	Status        string    `json:"status"`
	CategoryID    uint      `json:"category_id"`
	Thumbnail     string    `json:"thumbnail,omitempty"`
	DominantColor string    `json:"dominant_color,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
