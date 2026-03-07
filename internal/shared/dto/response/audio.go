package response

import "time"

type AudioResponse struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Artist     string    `json:"artist"`
	FilePath   string    `json:"file_path"`
	Duration   int       `json:"duration"`
	Status     string    `json:"status"`
	CategoryID uint      `json:"category_id"`
	Thumbnail  string    `json:"thumbnail,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
