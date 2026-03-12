package entity

import (
	"time"

	"gorm.io/gorm"
)

type Radio struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Title         string         `gorm:"not null;size:255" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	Thumbnail     string         `gorm:"size:500" json:"thumbnail"`
	DominantColor string         `gorm:"size:7" json:"dominant_color"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	Audios        []*Audio       `gorm:"many2many:radio_audios;" json:"audios,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Radio) TableName() string { return "radios" }

type RadioAudio struct {
	RadioID  uint `gorm:"primaryKey" json:"radio_id"`
	AudioID  uint `gorm:"primaryKey" json:"audio_id"`
	OrderNum int  `gorm:"default:0" json:"order_num"`
}

func (RadioAudio) TableName() string { return "radio_audios" }
