package user

import (
	"errors"

	"gorm.io/gorm"

	likeModel "mqfm-backend/internal/models/likes/user"

)

type UserLikeService struct {
	db *gorm.DB
}

func NewUserLikeService(db *gorm.DB) *UserLikeService {
	return &UserLikeService{db: db}
}

func (s *UserLikeService) LikeAudio(userID uint, audioID uint) error {
	var count int64
	s.db.Model(&likeModel.Like{}).Where("user_id = ? AND audio_id = ?", userID, audioID).Count(&count)
	if count > 0 {
		return errors.New("audio already liked")
	}

	like := likeModel.Like{
		UserID:  userID,
		AudioID: audioID,
	}
	return s.db.Create(&like).Error
}

func (s *UserLikeService) UnlikeAudio(userID uint, audioID uint) error {
	result := s.db.Where("user_id = ? AND audio_id = ?", userID, audioID).Delete(&likeModel.Like{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("like not found")
	}
	return nil
}

func (s *UserLikeService) GetLikedAudios(userID uint) ([]likeModel.Like, error) {
	var likes []likeModel.Like
	err := s.db.Where("user_id = ?", userID).Preload("Audio").Find(&likes).Error
	return likes, err
}