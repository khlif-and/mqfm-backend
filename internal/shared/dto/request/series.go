package request

type CreateSeriesRequest struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description"`
	Artist      string `form:"artist"`
}

type UpdateSeriesRequest struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	Artist      string `form:"artist"`
}

type AddSeriesItemRequest struct {
	SeriesID uint `json:"series_id" binding:"required"`
	AudioID  uint `json:"audio_id" binding:"required"`
	OrderNum int  `json:"order_num" binding:"required,min=1"`
}
