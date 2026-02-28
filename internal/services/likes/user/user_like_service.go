package user

import (
	"errors"

	likesDto "mqfm-backend/internal/dto/likes"
	likeModel "mqfm-backend/internal/models/likes/user"
	likeRepo "mqfm-backend/internal/repositories/likes/user"
)

type UserLikeService struct {
	repo likeRepo.LikeRepository
}

func NewUserLikeService(repo likeRepo.LikeRepository) *UserLikeService {
	return &UserLikeService{repo: repo}
}

func (s *UserLikeService) LikeAudio(userID uint, req likesDto.LikeRequest) (*likeModel.Like, error) {
	exists, err := s.repo.Exists(userID, req.AudioID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("audio already liked")
	}

	like := likeModel.Like{
		UserID:  userID,
		AudioID: req.AudioID,
	}

	if err := s.repo.Create(&like); err != nil {
		return nil, err
	}

	return &like, nil
}

func (s *UserLikeService) UnlikeAudio(userID uint, audioID uint) error {
	return s.repo.Delete(userID, audioID)
}

func (s *UserLikeService) GetLikedAudios(userID uint) ([]likeModel.Like, error) {
	return s.repo.FindByUser(userID)
}