package entity

import "time"

type UserLocation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	Latitude  float64   `gorm:"not null" json:"latitude"`
	Longitude float64   `gorm:"not null" json:"longitude"`
	City      string    `gorm:"size:255" json:"city"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (UserLocation) TableName() string {
	return "user_locations"
}
