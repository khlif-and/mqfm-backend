package service

import (
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/helper"
)

type colorExtractorService struct{}

func NewColorExtractorService() port.ColorExtractorService {
	return &colorExtractorService{}
}

func (s *colorExtractorService) ExtractDominantColor(imagePath string) (string, error) {
	return helper.ExtractDominantColor(imagePath)
}
