package service

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
)

type bookmarkService struct {
	repo port.BookmarkRepository
}

func NewBookmarkService(repo port.BookmarkRepository) port.BookmarkService {
	return &bookmarkService{repo: repo}
}

func (s *bookmarkService) Create(userID uint, req request.CreateBookmarkRequest) (*entity.Bookmark, error) {
	bookmark := &entity.Bookmark{
		UserID:          userID,
		AudioID:         req.AudioID,
		PositionSeconds: req.PositionSeconds,
		Label:           req.Label,
	}
	if err := s.repo.Create(bookmark); err != nil {
		return nil, err
	}
	return bookmark, nil
}

func (s *bookmarkService) GetByUser(userID uint) ([]entity.Bookmark, error) {
	return s.repo.FindByUser(userID)
}

func (s *bookmarkService) GetByUserAndAudio(userID, audioID uint) ([]entity.Bookmark, error) {
	return s.repo.FindByUserAndAudio(userID, audioID)
}

func (s *bookmarkService) Delete(id, userID uint) error {
	return s.repo.Delete(id, userID)
}
