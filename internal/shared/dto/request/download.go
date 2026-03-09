package request

type DownloadRequest struct {
	AudioID  uint  `json:"audio_id" binding:"required"`
	FileSize int64 `json:"file_size"`
}
