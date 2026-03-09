package request

type LikeRequest struct {
	TargetType string `json:"target_type" binding:"required,oneof=audio playlist"`
	TargetID   uint   `json:"target_id" binding:"required"`
}

type UnlikeRequest struct {
	TargetType string `json:"target_type" binding:"required,oneof=audio playlist"`
	TargetID   uint   `json:"target_id" binding:"required"`
}
