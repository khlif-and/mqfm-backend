package request

type CreateAudioRequest struct {
	Title       string `form:"title" binding:"required"`
	Artist      string `form:"artist" binding:"required"`
	Description string `form:"description"`
	Status      string `form:"status" binding:"required,oneof=active inactive"`
	CategoryID  uint   `form:"category_id" binding:"required"`
}

type UpdateAudioRequest struct {
	Title       string `form:"title"`
	Artist      string `form:"artist"`
	Description string `form:"description"`
	Status      string `form:"status"`
	CategoryID  uint   `form:"category_id"`
}
