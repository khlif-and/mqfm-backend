package request

type CreateBookmarkRequest struct {
	AudioID         uint   `json:"audio_id" binding:"required"`
	PositionSeconds int    `json:"position_seconds" binding:"required,min=0"`
	Label           string `json:"label"`
}
