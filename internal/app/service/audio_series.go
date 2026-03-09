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

type audioSeriesService struct {
	repo port.AudioSeriesRepository
}

func NewAudioSeriesService(repo port.AudioSeriesRepository) port.AudioSeriesService {
	return &audioSeriesService{repo: repo}
}

func (s *audioSeriesService) Create(req request.CreateSeriesRequest, file *multipart.FileHeader) (*entity.AudioSeries, error) {
	var imagePath string
	if file != nil {
		filename := helper.GenerateUniqueFilename(file.Filename)
		path := "uploads/thumbnails/" + filename
		if err := helper.SaveUploadedFile(file, path); err != nil {
			logger.Error("failed to save series image")
		} else {
			imagePath = path
		}
	}

	series := &entity.AudioSeries{
		Title:       req.Title,
		Description: req.Description,
		Artist:      req.Artist,
		Image:       imagePath,
	}

	if err := s.repo.Create(series); err != nil {
		return nil, err
	}
	return series, nil
}

func (s *audioSeriesService) FindAll() ([]entity.AudioSeries, error) {
	return s.repo.FindAll()
}

func (s *audioSeriesService) FindByID(id uint) (*entity.AudioSeries, error) {
	return s.repo.FindByID(id)
}

func (s *audioSeriesService) Update(id uint, req request.UpdateSeriesRequest, file *multipart.FileHeader) (*entity.AudioSeries, error) {
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Artist != "" {
		updates["artist"] = req.Artist
	}
	if file != nil {
		filename := helper.GenerateUniqueFilename(file.Filename)
		path := "uploads/thumbnails/" + filename
		if err := helper.SaveUploadedFile(file, path); err == nil {
			updates["image"] = path
		}
	}

	if err := s.repo.Update(id, updates); err != nil {
		return nil, err
	}
	return s.repo.FindByID(id)
}

func (s *audioSeriesService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *audioSeriesService) AddItem(req request.AddSeriesItemRequest) error {
	_, err := s.repo.FindByID(req.SeriesID)
	if err != nil {
		return errors.New("series not found")
	}

	return s.repo.AddItem(&entity.AudioSeriesItem{
		SeriesID: req.SeriesID,
		AudioID:  req.AudioID,
		OrderNum: req.OrderNum,
	})
}

func (s *audioSeriesService) RemoveItem(seriesID, audioID uint) error {
	return s.repo.RemoveItem(seriesID, audioID)
}

func (s *audioSeriesService) GetItems(seriesID uint) ([]entity.AudioSeriesItem, error) {
	return s.repo.FindItems(seriesID)
}

func (s *audioSeriesService) Search(query string) ([]entity.AudioSeries, error) {
	return s.repo.Search(query)
}
