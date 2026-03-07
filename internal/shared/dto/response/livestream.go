package response

import "time"

type LivestreamResponse struct {
	ID          uint      `json:"id"`
	IsLive      bool      `json:"is_live"`
	VideoID     string    `json:"video_id"`
	Title       string    `json:"title"`
	Thumbnail   string    `json:"thumbnail"`
	LastChecked time.Time `json:"last_checked"`
}
