package request

type DownloadRequest struct {
	AudioID    uint  `json:"audio_id" binding:"required"`
	PlaylistID uint  `json:"playlist_id"`
	FileSize   int64 `json:"file_size"`
}
