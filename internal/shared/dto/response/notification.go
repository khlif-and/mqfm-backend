package response

import "time"

type NotificationResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	Type        string    `json:"type"`
	ReferenceID uint      `json:"reference_id"`
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}

type NotificationSettingResponse struct {
	DailyReminder bool `json:"daily_reminder"`
	NewContent    bool `json:"new_content"`
	EventReminder bool `json:"event_reminder"`
}

type UnreadCountResponse struct {
	Count int64 `json:"count"`
}
