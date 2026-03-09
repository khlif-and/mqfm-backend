package request

type CreateClipRequest struct {
	AudioID   uint `json:"audio_id" binding:"required"`
	StartTime int  `json:"start_time" binding:"min=0"`
	EndTime   int  `json:"end_time" binding:"required,min=1"`
}
