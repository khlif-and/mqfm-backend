package entity

import "time"

type ListeningStat struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"not null;index:idx_stat_user_date" json:"user_id"`
	AudioID         uint      `gorm:"not null;index" json:"audio_id"`
	Audio           *Audio    `gorm:"foreignKey:AudioID" json:"audio,omitempty"`
	DurationSeconds int       `gorm:"default:0" json:"duration_seconds"`
	ListenedAt      time.Time `gorm:"index:idx_stat_user_date" json:"listened_at"`
}

func (ListeningStat) TableName() string {
	return "listening_stats"
}
