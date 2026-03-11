package user

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/shared/dto/response"
)

func toAudioResponsePtr(a *entity.Audio) *response.AudioResponse {
	if a == nil {
		return nil
	}
	r := toAudioResponseVal(*a)
	return &r
}

func toAudioResponseVal(a entity.Audio) response.AudioResponse {
	return response.AudioResponse{
		ID:            a.ID,
		Title:         a.Title,
		Artist:        a.Artist,
		Description:   a.Description,
		FilePath:      a.FilePath,
		Duration:      a.Duration,
		DurationFmt:   response.FormatDuration(a.Duration),
		FileSize:      a.FileSize,
		Status:        a.Status,
		CategoryID:    a.CategoryID,
		Thumbnail:     a.Thumbnail,
		DominantColor: a.DominantColor,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}

func toProgressResponse(p *entity.AudioProgress) *response.ProgressResponse {
	var pct float64
	if p.Duration > 0 {
		pct = float64(p.LastPosition) / float64(p.Duration) * 100
	}
	return &response.ProgressResponse{
		ID:           p.ID,
		UserID:       p.UserID,
		AudioID:      p.AudioID,
		Audio:        toAudioResponsePtr(p.Audio),
		LastPosition: p.LastPosition,
		Duration:     p.Duration,
		Percentage:   pct,
		Completed:    p.Completed,
		UpdatedAt:    p.UpdatedAt,
	}
}
