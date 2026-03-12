package service

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"time"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/internal/shared/helper"
)

const (
	radioListCacheKey = "cache:radios:active"
	radioDetailPrefix = "cache:radios:"
	radioCacheTTL     = 30 * time.Minute
)

type radioService struct {
	repo           port.RadioRepository
	cache          port.CacheRepository
	colorExtractor port.ColorExtractorService
}

func NewRadioService(repo port.RadioRepository, cache port.CacheRepository, colorExtractor port.ColorExtractorService) port.RadioService {
	return &radioService{repo: repo, cache: cache, colorExtractor: colorExtractor}
}

func (s *radioService) Create(req request.CreateRadioRequest, file *multipart.FileHeader) (*entity.Radio, error) {
	var thumbnailPath, dominantColor string
	if file != nil {
		filename := helper.GenerateUniqueFilename(file.Filename)
		path := "uploads/radios/" + filename
		if err := helper.SaveUploadedFile(file, path); err == nil {
			thumbnailPath = path
			if color, err := s.colorExtractor.ExtractDominantColor(path); err == nil {
				dominantColor = color
			}
		}
	}

	radio := &entity.Radio{
		Title:         req.Title,
		Description:   req.Description,
		Thumbnail:     thumbnailPath,
		DominantColor: dominantColor,
		IsActive:      true,
	}

	if err := s.repo.Create(radio); err != nil {
		return nil, err
	}

	s.invalidateListCache()
	return radio, nil
}

func (s *radioService) FindAll() ([]entity.Radio, error) {
	return s.repo.FindAll()
}

func (s *radioService) FindActive() ([]entity.Radio, error) {
	if s.cache != nil {
		ctx := context.Background()
		cached, err := s.cache.Get(ctx, radioListCacheKey)
		if err == nil && cached != "" {
			var radios []entity.Radio
			if json.Unmarshal([]byte(cached), &radios) == nil {
				return radios, nil
			}
		}
	}

	radios, err := s.repo.FindActive()
	if err != nil {
		return nil, err
	}

	if s.cache != nil && len(radios) > 0 {
		ctx := context.Background()
		_ = s.cache.Set(ctx, radioListCacheKey, radios, radioCacheTTL)
	}

	return radios, nil
}

func (s *radioService) FindByID(id uint) (*entity.Radio, error) {
	if s.cache != nil {
		ctx := context.Background()
		key := radioDetailKey(id)
		cached, err := s.cache.Get(ctx, key)
		if err == nil && cached != "" {
			var radio entity.Radio
			if json.Unmarshal([]byte(cached), &radio) == nil {
				return &radio, nil
			}
		}
	}

	radio, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		ctx := context.Background()
		_ = s.cache.Set(ctx, radioDetailKey(id), radio, radioCacheTTL)
	}

	return radio, nil
}

func (s *radioService) Update(id uint, req request.UpdateRadioRequest, file *multipart.FileHeader) (*entity.Radio, error) {
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if file != nil {
		filename := helper.GenerateUniqueFilename(file.Filename)
		path := "uploads/radios/" + filename
		if err := helper.SaveUploadedFile(file, path); err == nil {
			updates["thumbnail"] = path
			if color, err := s.colorExtractor.ExtractDominantColor(path); err == nil {
				updates["dominant_color"] = color
			}
		}
	}

	if len(updates) == 0 {
		return s.repo.FindByID(id)
	}

	if err := s.repo.Update(id, updates); err != nil {
		return nil, err
	}

	s.invalidateCache(id)
	return s.repo.FindByID(id)
}

func (s *radioService) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	s.invalidateCache(id)
	return nil
}

func (s *radioService) AddAudio(radioID, audioID uint) error {
	count, _ := s.repo.CountAudios(radioID)
	if err := s.repo.AddAudio(radioID, audioID, count+1); err != nil {
		return err
	}
	s.invalidateCache(radioID)
	s.updateDescription(radioID)
	return nil
}

func (s *radioService) RemoveAudio(radioID, audioID uint) error {
	if err := s.repo.RemoveAudio(radioID, audioID); err != nil {
		return err
	}
	s.invalidateCache(radioID)
	s.updateDescription(radioID)
	return nil
}

func (s *radioService) updateDescription(radioID uint) {
	audios, err := s.repo.FindAudios(radioID)
	if err != nil || len(audios) == 0 {
		return
	}

	seen := make(map[string]bool)
	var artists []string
	for _, a := range audios {
		if a.Artist != "" && !seen[a.Artist] {
			seen[a.Artist] = true
			artists = append(artists, a.Artist)
		}
	}

	desc := ""
	for i, name := range artists {
		if i > 0 {
			desc += ", "
		}
		if i >= 3 {
			desc += "dan lainnya..."
			break
		}
		desc += name
	}

	_ = s.repo.Update(radioID, map[string]interface{}{"description": desc})
}

func (s *radioService) invalidateCache(radioID uint) {
	if s.cache == nil {
		return
	}
	ctx := context.Background()
	_ = s.cache.Delete(ctx, radioListCacheKey)
	_ = s.cache.Delete(ctx, radioDetailKey(radioID))
}

func (s *radioService) invalidateListCache() {
	if s.cache == nil {
		return
	}
	_ = s.cache.Delete(context.Background(), radioListCacheKey)
}

func radioDetailKey(id uint) string {
	return fmt.Sprintf("%s%d", radioDetailPrefix, id)
}
