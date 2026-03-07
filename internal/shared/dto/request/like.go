package request

type LikeRequest struct {
	AudioID uint `json:"audio_id" binding:"required"`
}
