package admin

import (
	"time"

	"gorm.io/gorm"

)

type Audio struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Artist      string         `json:"artist"`
	Description string         `json:"description"`
	FilePath    string         `json:"file_path"` // Renamed from AudioURL to match Service/DTO
	Duration    int            `json:"duration"`
	Status      string         `gorm:"default:active" json:"status"`
	Thumbnail   string         `json:"thumbnail"`
	CategoryID  uint           `json:"category_id"`
	// Category    categoryModel.Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"` // Circular import risk, handle carefully
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Audio) TableName() string {
	return "audios"
}