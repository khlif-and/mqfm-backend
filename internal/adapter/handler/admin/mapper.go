package admin

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/shared/dto/response"
)

func toCategoryResponse(c *entity.Category) response.CategoryResponse {
	return response.CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		Image:     c.Image,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func toCategoryResponseVal(c entity.Category) response.CategoryResponse {
	return response.CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		Image:     c.Image,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func toAudioResponse(a *entity.Audio) response.AudioResponse {
	return response.AudioResponse{
		ID:            a.ID,
		Title:         a.Title,
		Artist:        a.Artist,
		Description:   a.Description,
		FilePath:      a.FilePath,
		Duration:      a.Duration,
		Status:        a.Status,
		CategoryID:    a.CategoryID,
		Thumbnail:     a.Thumbnail,
		DominantColor: a.DominantColor,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}

func toAudioResponseVal(a entity.Audio) response.AudioResponse {
	return response.AudioResponse{
		ID:            a.ID,
		Title:         a.Title,
		Artist:        a.Artist,
		Description:   a.Description,
		FilePath:      a.FilePath,
		Duration:      a.Duration,
		Status:        a.Status,
		CategoryID:    a.CategoryID,
		Thumbnail:     a.Thumbnail,
		DominantColor: a.DominantColor,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}
