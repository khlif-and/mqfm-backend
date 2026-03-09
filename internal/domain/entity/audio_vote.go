package entity

import "time"

type AudioVote struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_vote_user_audio" json:"user_id"`
	AudioID   uint      `gorm:"not null;uniqueIndex:idx_vote_user_audio;index" json:"audio_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (AudioVote) TableName() string {
	return "audio_votes"
}
