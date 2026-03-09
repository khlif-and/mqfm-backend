package entity

import "time"

type PlaylistCollaborator struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PlaylistID uint      `gorm:"not null;uniqueIndex:idx_collab_playlist_user" json:"playlist_id"`
	UserID     uint      `gorm:"not null;uniqueIndex:idx_collab_playlist_user" json:"user_id"`
	Role       string    `gorm:"size:20;default:contributor" json:"role"`
	CreatedAt  time.Time `json:"created_at"`
}

func (PlaylistCollaborator) TableName() string {
	return "playlist_collaborators"
}
