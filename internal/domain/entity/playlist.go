package entity

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Playlist struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	CreatorRole   string         `gorm:"size:10;default:user" json:"creator_role"`
	Name          string         `gorm:"not null" json:"name"`
	ImageURL      string         `json:"image_url"`
	DominantColor string         `gorm:"size:7" json:"dominant_color"`
	ShareToken    string         `gorm:"size:64;uniqueIndex" json:"share_token"`
	IsPublic      bool           `gorm:"default:false" json:"is_public"`
	Audios        []*Audio       `gorm:"many2many:playlist_audios;" json:"audios"`
	User          *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Playlist) TableName() string {
	return "playlists"
}

func (p *Playlist) TimeSinceCreated() string {
	diff := time.Since(p.CreatedAt)
	days := int(diff.Hours() / 24)
	years := days / 365
	months := (days % 365) / 30
	remaining := (days % 365) % 30
	return fmt.Sprintf("%dy %dm %dd", years, months, remaining)
}
