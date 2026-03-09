package response

import "time"

type DownloadResponse struct {
	ID        uint           `json:"id"`
	UserID    uint           `json:"user_id"`
	AudioID   uint           `json:"audio_id"`
	Audio     *AudioResponse `json:"audio,omitempty"`
	FileSize  int64          `json:"file_size"`
	CreatedAt time.Time      `json:"created_at"`
}

type StorageUsageResponse struct {
	TotalBytes int64 `json:"total_bytes"`
	TotalMB    int64 `json:"total_mb"`
}
