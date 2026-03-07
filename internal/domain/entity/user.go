package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Username       string         `gorm:"unique;not null" json:"username"`
	Email          string         `gorm:"unique;not null" json:"email"`
	Password       string         `json:"-"`
	ProfilePicture string         `json:"profile_picture"`
	Role           string         `gorm:"default:user" json:"role"`
	Provider       string         `gorm:"default:local" json:"provider"`
	ProviderID     string         `json:"provider_id"`
	EmailVerified  bool           `gorm:"default:false" json:"email_verified"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}
