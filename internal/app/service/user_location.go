package service

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
)

type userLocationService struct {
	repo port.UserLocationRepository
}

func NewUserLocationService(repo port.UserLocationRepository) port.UserLocationService {
	return &userLocationService{repo: repo}
}

func (s *userLocationService) Update(userID uint, req request.UpdateLocationRequest) (*entity.UserLocation, error) {
	loc := &entity.UserLocation{
		UserID:    userID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		City:      req.City,
	}

	if err := s.repo.Upsert(loc); err != nil {
		return nil, err
	}
	return s.repo.FindByUser(userID)
}

func (s *userLocationService) Get(userID uint) (*entity.UserLocation, error) {
	return s.repo.FindByUser(userID)
}
