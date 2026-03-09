package request

type UpdateNotificationSettingRequest struct {
	DailyReminder *bool `json:"daily_reminder"`
	NewContent    *bool `json:"new_content"`
	EventReminder *bool `json:"event_reminder"`
}
