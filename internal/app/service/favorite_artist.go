package service

import (
	"errors"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type favoriteArtistService struct {
	repo port.FavoriteArtistRepository
}

func NewFavoriteArtistService(repo port.FavoriteArtistRepository) port.FavoriteArtistService {
	return &favoriteArtistService{repo: repo}
}

func (s *favoriteArtistService) Add(userID uint, artistName string) error {
	exists, err := s.repo.Exists(userID, artistName)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("artist already in favorites")
	}

	return s.repo.Create(&entity.FavoriteArtist{
		UserID:     userID,
		ArtistName: artistName,
	})
}

func (s *favoriteArtistService) Remove(userID uint, artistName string) error {
	return s.repo.Delete(userID, artistName)
}

func (s *favoriteArtistService) GetByUser(userID uint) ([]entity.FavoriteArtist, error) {
	return s.repo.FindByUser(userID)
}
