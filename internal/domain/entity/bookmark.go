package entity

import "time"

type Bookmark struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"not null;index:idx_bookmark_user_audio,unique" json:"user_id"`
	AudioID         uint      `gorm:"not null;index:idx_bookmark_user_audio,unique" json:"audio_id"`
	Audio           *Audio    `gorm:"foreignKey:AudioID" json:"audio,omitempty"`
	PositionSeconds int       `gorm:"not null;index:idx_bookmark_user_audio,unique" json:"position_seconds"`
	Label           string    `gorm:"size:255" json:"label"`
	CreatedAt       time.Time `json:"created_at"`
}

func (Bookmark) TableName() string {
	return "bookmarks"
}
