package entity

import "time"

type Like struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index:idx_like_unique,unique" json:"user_id"`
	TargetType string    `gorm:"not null;size:20;index:idx_like_unique,unique" json:"target_type"`
	TargetID   uint      `gorm:"not null;index:idx_like_unique,unique" json:"target_id"`
	Audio      *Audio    `gorm:"foreignKey:TargetID" json:"audio,omitempty"`
	Playlist   *Playlist `gorm:"foreignKey:TargetID" json:"playlist,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

func (Like) TableName() string {
	return "likes"
}
