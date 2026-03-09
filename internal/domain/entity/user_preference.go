package entity

import "time"

type UserPreference struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	UserID            uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	PlaybackSpeed     float64   `gorm:"default:1.0" json:"playback_speed"`
	SleepTimerMinutes int       `gorm:"default:0" json:"sleep_timer_minutes"`
	AutoDownloadWifi  bool      `gorm:"default:false" json:"auto_download_wifi"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (UserPreference) TableName() string {
	return "user_preferences"
}
