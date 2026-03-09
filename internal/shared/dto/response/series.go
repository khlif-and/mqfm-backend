package response

import "time"

type SeriesResponse struct {
	ID          uint               `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Artist      string             `json:"artist"`
	Image       string             `json:"image"`
	Items       []SeriesItemResponse `json:"items,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type SeriesItemResponse struct {
	ID       uint           `json:"id"`
	SeriesID uint           `json:"series_id"`
	AudioID  uint           `json:"audio_id"`
	Audio    *AudioResponse `json:"audio,omitempty"`
	OrderNum int            `json:"order_num"`
}
