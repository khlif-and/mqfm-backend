package livestream

import "time"

type LiveStream struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	IsLive      bool      `json:"is_live"`
	VideoID     string    `json:"video_id"`
	Title       string    `json:"title"`
	Thumbnail   string    `json:"thumbnail"`
	LastChecked time.Time `json:"last_checked"`
}