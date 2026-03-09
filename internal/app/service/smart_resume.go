package service

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
)

type smartResumeService struct {
	repo port.SmartResumeRepository
}

func NewSmartResumeService(repo port.SmartResumeRepository) port.SmartResumeService {
	return &smartResumeService{repo: repo}
}

func (s *smartResumeService) Update(userID uint, req request.UpdateResumeRequest) (*entity.SmartResume, error) {
	resume := &entity.SmartResume{
		UserID:          userID,
		AudioID:         req.AudioID,
		PlaylistID:      req.PlaylistID,
		PositionSeconds: req.PositionSeconds,
	}

	if err := s.repo.Upsert(resume); err != nil {
		return nil, err
	}
	return s.repo.FindByUser(userID)
}

func (s *smartResumeService) Get(userID uint) (*entity.SmartResume, error) {
	return s.repo.FindByUser(userID)
}
