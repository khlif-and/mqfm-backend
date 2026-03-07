package service

import (
	"errors"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/dto/request"
)

type likeService struct {
	repo port.LikeRepository
}

func NewLikeService(repo port.LikeRepository) port.LikeService {
	return &likeService{repo: repo}
}

func (s *likeService) LikeAudio(userID uint, req request.LikeRequest) (*entity.Like, error) {
	exists, err := s.repo.Exists(userID, req.AudioID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constant.MsgAlreadyLiked)
	}

	like := entity.Like{
		UserID:  userID,
		AudioID: req.AudioID,
	}

	if err := s.repo.Create(&like); err != nil {
		return nil, err
	}

	return &like, nil
}

func (s *likeService) UnlikeAudio(userID, audioID uint) error {
	return s.repo.Delete(userID, audioID)
}

func (s *likeService) GetLikedAudios(userID uint) ([]entity.Like, error) {
	return s.repo.FindByUser(userID)
}
