package entity

import "time"

type Download struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index:idx_download_user_audio,unique" json:"user_id"`
	AudioID    uint      `gorm:"not null;index:idx_download_user_audio,unique" json:"audio_id"`
	Audio      *Audio    `gorm:"foreignKey:AudioID" json:"audio,omitempty"`
	PlaylistID uint      `gorm:"default:0" json:"playlist_id"`
	FileSize   int64     `gorm:"default:0" json:"file_size"`
	ExpiresAt  time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
}

func (Download) TableName() string {
	return "downloads"
}
