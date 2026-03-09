package request

type UpdateProgressRequest struct {
	AudioID      uint    `json:"audio_id" binding:"required"`
	LastPosition int     `json:"last_position" binding:"min=0"`
	Duration     int     `json:"duration" binding:"min=0"`
	Percentage   float64 `json:"percentage" binding:"min=0,max=100"`
	Completed    bool    `json:"completed"`
}

type RecordStatRequest struct {
	AudioID         uint `json:"audio_id" binding:"required"`
	DurationSeconds int  `json:"duration_seconds" binding:"required,min=1"`
}
