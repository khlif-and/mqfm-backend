package user

import (
	"time"

	adminAudioModel "mqfm-backend/internal/models/podcast/audio/admin"

)

type Like struct {
	ID        uint                   `gorm:"primaryKey" json:"id"`
	UserID    uint                   `gorm:"not null;index:idx_user_audio,unique" json:"user_id"`
	AudioID   uint                   `gorm:"not null;index:idx_user_audio,unique" json:"audio_id"`
	Audio     *adminAudioModel.Audio `gorm:"foreignKey:AudioID" json:"audio"`
	CreatedAt time.Time              `json:"created_at"`
}

func (Like) TableName() string {
	return "likes"
}