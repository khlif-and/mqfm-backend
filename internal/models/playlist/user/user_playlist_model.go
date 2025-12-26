package user

import (
	"time"

	"gorm.io/gorm"

	adminAudioModel "mqfm-backend/internal/models/podcast/audio/admin"

)

type Playlist struct {
	ID        uint                     `gorm:"primaryKey" json:"id"`
	UserID    uint                     `gorm:"not null;index" json:"user_id"`
	Name      string                   `gorm:"not null" json:"name"`
	ImageURL  string                   `json:"image_url"`
	Audios    []*adminAudioModel.Audio `gorm:"many2many:playlist_audios;" json:"audios"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
	DeletedAt gorm.DeletedAt           `gorm:"index" json:"-"`
}

func (Playlist) TableName() string {
	return "playlists"
}