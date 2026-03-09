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

func toEventResponse(e *entity.Event) response.EventResponse {
	return response.EventResponse{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		EventDate:   e.EventDate,
		Location:    e.Location,
		Image:       e.Image,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func toEventResponseVal(e entity.Event) response.EventResponse {
	return response.EventResponse{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		EventDate:   e.EventDate,
		Location:    e.Location,
		Image:       e.Image,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func toSeriesResponse(s *entity.AudioSeries) response.SeriesResponse {
	r := response.SeriesResponse{
		ID:          s.ID,
		Title:       s.Title,
		Description: s.Description,
		Artist:      s.Artist,
		Image:       s.Image,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
	for _, item := range s.Items {
		si := response.SeriesItemResponse{
			ID:       item.ID,
			SeriesID: item.SeriesID,
			AudioID:  item.AudioID,
			OrderNum: item.OrderNum,
		}
		if item.Audio != nil {
			ar := toAudioResponseVal(*item.Audio)
			si.Audio = &ar
		}
		r.Items = append(r.Items, si)
	}
	return r
}

func toSeriesResponseVal(s entity.AudioSeries) response.SeriesResponse {
	return toSeriesResponse(&s)
}
