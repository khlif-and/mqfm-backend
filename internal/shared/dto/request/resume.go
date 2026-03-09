package request

type UpdateResumeRequest struct {
	AudioID         uint `json:"audio_id" binding:"required"`
	PlaylistID      uint `json:"playlist_id"`
	PositionSeconds int  `json:"position_seconds" binding:"min=0"`
}
