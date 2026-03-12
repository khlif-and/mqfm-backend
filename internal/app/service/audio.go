package service

import (
	"errors"
	"mime/multipart"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/internal/shared/helper"
	"mqfm-backend/internal/shared/logger"
)

type audioService struct {
	repo           port.AudioRepository
	colorExtractor port.ColorExtractorService
	converter      port.AudioConverterService
}

func NewAudioService(repo port.AudioRepository, colorExtractor port.ColorExtractorService, converter port.AudioConverterService) port.AudioService {
	return &audioService{repo: repo, colorExtractor: colorExtractor, converter: converter}
}

func (s *audioService) Create(req request.CreateAudioRequest, audioFile *multipart.FileHeader, thumbnailFile *multipart.FileHeader) (*entity.Audio, error) {
	var filePath string
	if audioFile != nil {
		filename := helper.GenerateUniqueFilename(audioFile.Filename)
		path := "uploads/audios/" + filename
		if err := helper.SaveUploadedFile(audioFile, path); err != nil {
			logger.Error("failed to save audio file")
		} else {
			filePath = path
		}
	}

	var thumbnailPath, dominantColor string
	if thumbnailFile != nil {
		filename := helper.GenerateUniqueFilename(thumbnailFile.Filename)
		path := "uploads/thumbnails/" + filename
		if err := helper.SaveUploadedFile(thumbnailFile, path); err != nil {
			logger.Error("failed to save thumbnail")
		} else {
			thumbnailPath = path
			if color, err := s.colorExtractor.ExtractDominantColor(path); err == nil {
				dominantColor = color
			}
		}
	}

	audio := entity.Audio{
		Title:         req.Title,
		Artist:        req.Artist,
		Description:   req.Description,
		Status:        req.Status,
		CategoryID:    req.CategoryID,
		FilePath:      filePath,
		Thumbnail:     thumbnailPath,
		DominantColor: dominantColor,
	}

	if filePath != "" {
		audio.Duration = helper.ExtractAudioDuration(filePath)
		audio.FileSize = helper.GetFileSize(filePath)
	}

	if err := s.repo.Create(&audio); err != nil {
		return nil, err
	}

	if audio.FilePath != "" && s.converter != nil {
		go s.convertToOGG(audio.ID, audio.FilePath)
	}

	return &audio, nil
}

func (s *audioService) FindAll() ([]entity.Audio, error) {
	return s.repo.FindAll()
}

func (s *audioService) FindByID(id uint) (*entity.Audio, error) {
	return s.repo.FindByID(id)
}

func (s *audioService) Update(id uint, req request.UpdateAudioRequest, audioFile *multipart.FileHeader, thumbnailFile *multipart.FileHeader) (*entity.Audio, error) {
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Artist != "" {
		updates["artist"] = req.Artist
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.CategoryID != 0 {
		updates["category_id"] = req.CategoryID
	}

	if audioFile != nil {
		filename := helper.GenerateUniqueFilename(audioFile.Filename)
		path := "uploads/audios/" + filename
		if err := helper.SaveUploadedFile(audioFile, path); err != nil {
			logger.Error("failed to save audio file")
		} else {
			updates["file_path"] = path
			updates["duration"] = helper.ExtractAudioDuration(path)
			updates["file_size"] = helper.GetFileSize(path)
		}
	}

	if thumbnailFile != nil {
		filename := helper.GenerateUniqueFilename(thumbnailFile.Filename)
		path := "uploads/thumbnails/" + filename
		if err := helper.SaveUploadedFile(thumbnailFile, path); err != nil {
			logger.Error("failed to save thumbnail")
		} else {
			updates["thumbnail"] = path
			if color, err := s.colorExtractor.ExtractDominantColor(path); err == nil {
				updates["dominant_color"] = color
			}
		}
	}

	if len(updates) == 0 {
		return nil, errors.New("no updates provided")
	}

	if err := s.repo.Update(id, updates); err != nil {
		return nil, err
	}

	if path, ok := updates["file_path"].(string); ok && s.converter != nil {
		go s.convertToOGG(id, path)
	}

	return s.repo.FindByID(id)
}

func (s *audioService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *audioService) Search(query string) ([]entity.Audio, error) {
	return s.repo.Search(query)
}

func (s *audioService) convertToOGG(audioID uint, filePath string) {
	oggPath, err := s.converter.ConvertToOGG(filePath)
	if err != nil {
		logger.Error("ogg conversion failed for audio " + filePath)
		return
	}
	_ = s.repo.Update(audioID, map[string]interface{}{"ogg_path": oggPath})
}
