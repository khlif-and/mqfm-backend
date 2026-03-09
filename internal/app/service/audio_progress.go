package service

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
)

type audioProgressService struct {
	repo port.AudioProgressRepository
}

func NewAudioProgressService(repo port.AudioProgressRepository) port.AudioProgressService {
	return &audioProgressService{repo: repo}
}

func (s *audioProgressService) UpdateProgress(userID uint, req request.UpdateProgressRequest) (*entity.AudioProgress, error) {
	progress := &entity.AudioProgress{
		UserID:       userID,
		AudioID:      req.AudioID,
		LastPosition: req.LastPosition,
		Duration:     req.Duration,
		Percentage:   req.Percentage,
		Completed:    req.Completed || req.Percentage >= 95,
	}
	if err := s.repo.Upsert(progress); err != nil {
		return nil, err
	}
	return progress, nil
}

func (s *audioProgressService) GetProgress(userID, audioID uint) (*entity.AudioProgress, error) {
	return s.repo.FindByUserAndAudio(userID, audioID)
}

func (s *audioProgressService) GetAllProgress(userID uint) ([]entity.AudioProgress, error) {
	return s.repo.FindByUser(userID)
}

func (s *audioProgressService) GetCompleted(userID uint) ([]entity.AudioProgress, error) {
	return s.repo.FindCompletedByUser(userID)
}
