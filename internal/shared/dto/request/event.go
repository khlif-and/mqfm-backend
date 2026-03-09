package request

type CreateEventRequest struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description"`
	EventDate   string `form:"event_date" binding:"required"`
	Location    string `form:"location"`
}

type UpdateEventRequest struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	EventDate   string `form:"event_date"`
	Location    string `form:"location"`
}
