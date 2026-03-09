package response

import "time"

type EventResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date"`
	Location    string    `json:"location"`
	Image       string    `json:"image"`
	RSVPCount   int64     `json:"rsvp_count"`
	HasRSVP     bool      `json:"has_rsvp"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EventRSVPResponse struct {
	ID        uint           `json:"id"`
	UserID    uint           `json:"user_id"`
	EventID   uint           `json:"event_id"`
	Event     *EventResponse `json:"event,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}
