package entity

import "time"

type FavoriteArtist struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;uniqueIndex:idx_fav_user_artist" json:"user_id"`
	ArtistName string    `gorm:"not null;size:255;uniqueIndex:idx_fav_user_artist" json:"artist_name"`
	CreatedAt  time.Time `json:"created_at"`
}

func (FavoriteArtist) TableName() string {
	return "favorite_artists"
}
