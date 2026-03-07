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
	repo port.AudioRepository
}

func NewAudioService(repo port.AudioRepository) port.AudioService {
	return &audioService{repo: repo}
}

func (s *audioService) Create(req request.CreateAudioRequest, file *multipart.FileHeader) (*entity.Audio, error) {
	var filePath string
	if file != nil {
		filename := helper.GenerateUniqueFilename(file.Filename)
		path := "uploads/audios/" + filename
		if err := helper.SaveUploadedFile(file, path); err != nil {
			logger.Error("failed to save audio file")
		} else {
			filePath = path
		}
	}

	audio := entity.Audio{
		Title:      req.Title,
		Artist:     req.Artist,
		Status:     req.Status,
		CategoryID: req.CategoryID,
		FilePath:   filePath,
	}

	if err := s.repo.Create(&audio); err != nil {
		return nil, err
	}

	return &audio, nil
}

func (s *audioService) FindAll() ([]entity.Audio, error) {
	return s.repo.FindAll()
}

func (s *audioService) FindByID(id uint) (*entity.Audio, error) {
	return s.repo.FindByID(id)
}

func (s *audioService) Update(id uint, req request.UpdateAudioRequest, file *multipart.FileHeader) (*entity.Audio, error) {
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Artist != "" {
		updates["artist"] = req.Artist
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.CategoryID != 0 {
		updates["category_id"] = req.CategoryID
	}

	if file != nil {
		filename := helper.GenerateUniqueFilename(file.Filename)
		path := "uploads/audios/" + filename
		if err := helper.SaveUploadedFile(file, path); err != nil {
			logger.Error("failed to save audio file")
		} else {
			updates["file_path"] = path
		}
	}

	if len(updates) == 0 {
		return nil, errors.New("no updates provided")
	}

	if err := s.repo.Update(id, updates); err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

func (s *audioService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *audioService) Search(query string) ([]entity.Audio, error) {
	return s.repo.Search(query)
}
