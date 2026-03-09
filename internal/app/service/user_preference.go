package service

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
)

type userPreferenceService struct {
	repo port.UserPreferenceRepository
}

func NewUserPreferenceService(repo port.UserPreferenceRepository) port.UserPreferenceService {
	return &userPreferenceService{repo: repo}
}

func (s *userPreferenceService) GetOrCreate(userID uint) (*entity.UserPreference, error) {
	pref, err := s.repo.FindByUser(userID)
	if err != nil {
		pref = &entity.UserPreference{
			UserID:        userID,
			PlaybackSpeed: 1.0,
		}
		if err := s.repo.Upsert(pref); err != nil {
			return nil, err
		}
		return pref, nil
	}
	return pref, nil
}

func (s *userPreferenceService) Update(userID uint, req request.UpdatePreferenceRequest) (*entity.UserPreference, error) {
	pref, _ := s.GetOrCreate(userID)

	if req.PlaybackSpeed != nil {
		pref.PlaybackSpeed = *req.PlaybackSpeed
	}
	if req.SleepTimerMinutes != nil {
		pref.SleepTimerMinutes = *req.SleepTimerMinutes
	}
	if req.AutoDownloadWifi != nil {
		pref.AutoDownloadWifi = *req.AutoDownloadWifi
	}

	if err := s.repo.Upsert(pref); err != nil {
		return nil, err
	}
	return pref, nil
}
