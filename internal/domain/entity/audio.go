package entity

import (
	"time"

	"gorm.io/gorm"
)

type Audio struct {
	ID            uint           `gorm:"primaryKey" json:"audio_id"`
	Title         string         `gorm:"not null" json:"title"`
	Artist        string         `json:"artist"`
	Description   string         `json:"description"`
	FilePath      string         `json:"file_path"`
	Duration      int            `json:"duration"`
	Status        string         `gorm:"default:active" json:"status"`
	Thumbnail     string         `json:"thumbnail"`
	DominantColor string         `gorm:"size:7" json:"dominant_color"`
	CategoryID    uint           `json:"category_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Audio) TableName() string {
	return "audios"
}
