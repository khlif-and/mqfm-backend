package entity

import "time"

type SmartResume struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	AudioID         uint      `gorm:"not null" json:"audio_id"`
	Audio           *Audio    `gorm:"foreignKey:AudioID" json:"audio,omitempty"`
	PlaylistID      uint      `gorm:"default:0" json:"playlist_id"`
	PositionSeconds int       `gorm:"default:0" json:"position_seconds"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (SmartResume) TableName() string {
	return "smart_resumes"
}
