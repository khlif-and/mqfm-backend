package admin

import (
	"time"

	"gorm.io/gorm"

)

type Audio struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	AudioURL    string         `json:"audio_url"` 
	Thumbnail   string         `json:"thumbnail"`
	CategoryID  uint           `json:"category_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Audio) TableName() string {
	return "audios"
}