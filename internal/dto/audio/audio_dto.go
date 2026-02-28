package audio

import "time"

type CreateAudioRequest struct {
	Title      string `form:"title" binding:"required"`
	Artist     string `form:"artist" binding:"required"`
	Status     string `form:"status" binding:"required,oneof=active inactive"`
	CategoryID uint   `form:"category_id" binding:"required"`
}

type UpdateAudioRequest struct {
	Title      string `form:"title"`
	Artist     string `form:"artist"`
	Status     string `form:"status,oneof=active inactive"`
	CategoryID uint   `form:"category_id"`
}

type AudioResponse struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Artist     string    `json:"artist"`
	FilePath   string    `json:"file_path"`
	Duration   int       `json:"duration"`
	Status     string    `json:"status"`
	CategoryID uint      `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
