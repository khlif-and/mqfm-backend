package user

import (
	"time"

	adminAudioModel "mqfm-backend/internal/models/podcast/audio/admin"
)

type History struct {
	ID        uint                   `gorm:"primaryKey" json:"id"`
	UserID    uint                   `gorm:"not null;index:idx_user_audio_history,unique" json:"user_id"`
	AudioID   uint                   `gorm:"not null;index:idx_user_audio_history,unique" json:"audio_id"`
	Audio     *adminAudioModel.Audio `gorm:"foreignKey:AudioID" json:"audio"`
	PlayCount int                    `gorm:"default:1" json:"play_count"`
	PlayedAt  time.Time              `gorm:"autoUpdateTime" json:"played_at"`
	CreatedAt time.Time              `json:"created_at"`
}

func (History) TableName() string {
	return "histories"
}
