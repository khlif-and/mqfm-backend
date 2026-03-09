package entity

import (
	"time"

	"gorm.io/gorm"
)

type Playlist struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Name      string         `gorm:"not null" json:"name"`
	ImageURL  string         `json:"image_url"`
	ShareToken string        `gorm:"size:64;uniqueIndex" json:"share_token"`
	IsPublic  bool           `gorm:"default:false" json:"is_public"`
	Audios    []*Audio       `gorm:"many2many:playlist_audios;" json:"audios"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Playlist) TableName() string {
	return "playlists"
}
