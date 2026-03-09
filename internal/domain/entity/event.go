package entity

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null;size:255" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	EventDate   time.Time      `gorm:"not null;index" json:"event_date"`
	Location    string         `gorm:"size:500" json:"location"`
	Image       string         `gorm:"size:500" json:"image"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Event) TableName() string {
	return "events"
}

type EventRSVP struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_rsvp_user_event" json:"user_id"`
	EventID   uint      `gorm:"not null;uniqueIndex:idx_rsvp_user_event" json:"event_id"`
	Event     *Event    `gorm:"foreignKey:EventID" json:"event,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (EventRSVP) TableName() string {
	return "event_rsvps"
}
