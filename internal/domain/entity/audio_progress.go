package entity

import "time"

type AudioProgress struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;uniqueIndex:idx_progress_user_audio" json:"user_id"`
	AudioID      uint      `gorm:"not null;uniqueIndex:idx_progress_user_audio" json:"audio_id"`
	Audio        *Audio    `gorm:"foreignKey:AudioID" json:"audio,omitempty"`
	LastPosition int       `gorm:"default:0" json:"last_position"`
	Duration     int       `gorm:"default:0" json:"duration"`
	Percentage   float64   `gorm:"default:0" json:"percentage"`
	Completed    bool      `gorm:"default:false" json:"completed"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (AudioProgress) TableName() string {
	return "audio_progress"
}
