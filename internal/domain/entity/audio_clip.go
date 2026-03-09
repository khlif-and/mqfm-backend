package entity

import "time"

type AudioClip struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	AudioID    uint      `gorm:"not null;index" json:"audio_id"`
	Audio      *Audio    `gorm:"foreignKey:AudioID" json:"audio,omitempty"`
	StartTime  int       `gorm:"not null" json:"start_time"`
	EndTime    int       `gorm:"not null" json:"end_time"`
	ClipPath   string    `gorm:"size:500" json:"clip_path"`
	ShareToken string    `gorm:"size:64;uniqueIndex" json:"share_token"`
	CreatedAt  time.Time `json:"created_at"`
}

func (AudioClip) TableName() string {
	return "audio_clips"
}
