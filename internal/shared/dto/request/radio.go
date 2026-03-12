package request

type CreateRadioRequest struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description"`
}

type UpdateRadioRequest struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	IsActive    *bool  `form:"is_active"`
}

type RadioAudioRequest struct {
	RadioID uint `json:"radio_id" binding:"required"`
	AudioID uint `json:"audio_id" binding:"required"`
}
