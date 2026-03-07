package response

import "time"

type AudioScoreResponse struct {
	ID            uint      `json:"id"`
	Title         string    `json:"title"`
	Artist        string    `json:"artist"`
	FilePath      string    `json:"file_path"`
	Duration      int       `json:"duration"`
	Status        string    `json:"status"`
	Thumbnail     string    `json:"thumbnail,omitempty"`
	DominantColor string    `json:"dominant_color,omitempty"`
	CategoryID    uint      `json:"category_id"`
	TotalPlays    int64     `json:"total_plays"`
	TotalLikes    int64     `json:"total_likes"`
	WeightScore   float64   `json:"weight_score"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
