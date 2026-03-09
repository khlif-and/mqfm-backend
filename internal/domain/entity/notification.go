package entity

import "time"

type Notification struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Title       string    `gorm:"not null;size:255" json:"title"`
	Body        string    `gorm:"type:text" json:"body"`
	Type        string    `gorm:"size:50;index" json:"type"`
	ReferenceID uint      `gorm:"default:0" json:"reference_id"`
	IsRead      bool      `gorm:"default:false" json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Notification) TableName() string {
	return "notifications"
}

type NotificationSetting struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	DailyReminder bool      `gorm:"default:true" json:"daily_reminder"`
	NewContent    bool      `gorm:"default:true" json:"new_content"`
	EventReminder bool      `gorm:"default:true" json:"event_reminder"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (NotificationSetting) TableName() string {
	return "notification_settings"
}
