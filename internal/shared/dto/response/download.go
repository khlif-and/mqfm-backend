package response

import (
	"fmt"
	"time"
)

type DownloadResponse struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	AudioID       uint      `json:"audio_id"`
	PlaylistID    uint      `json:"playlist_id,omitempty"`
	Audio         *AudioResponse `json:"audio,omitempty"`
	Title         string    `json:"title"`
	Artist        string    `json:"artist"`
	Thumbnail     string    `json:"thumbnail"`
	DominantColor string    `json:"dominant_color,omitempty"`
	Duration      int       `json:"duration"`
	DurationFmt   string    `json:"duration_fmt"`
	FileSize      int64     `json:"file_size"`
	ExpiresAt     time.Time `json:"expires_at"`
	DaysRemaining int       `json:"days_remaining"`
	CreatedAt     time.Time `json:"created_at"`
}

type StorageUsageResponse struct {
	TotalBytes int64 `json:"total_bytes"`
	TotalMB    int64 `json:"total_mb"`
}

func FormatDuration(seconds int) string {
	m := seconds / 60
	s := seconds % 60
	return fmt.Sprintf("%02d:%02d", m, s)
}
