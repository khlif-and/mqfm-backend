package user

import (
	"errors"
	"gorm.io/gorm"
	likeModel "mqfm-backend/internal/models/likes/user"
)

type LikeRepository interface {
	Create(like *likeModel.Like) error
	Delete(userID, audioID uint) error
	FindByUser(userID uint) ([]likeModel.Like, error)
	Exists(userID, audioID uint) (bool, error)
}

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) Create(like *likeModel.Like) error {
	return r.db.Create(like).Error
}

func (r *likeRepository) Delete(userID, audioID uint) error {
	result := r.db.Where("user_id = ? AND audio_id = ?", userID, audioID).Delete(&likeModel.Like{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("like not found")
	}
	return nil
}

func (r *likeRepository) FindByUser(userID uint) ([]likeModel.Like, error) {
	var likes []likeModel.Like
	err := r.db.Preload("Audio").Where("user_id = ?", userID).Find(&likes).Error
	return likes, err
}

func (r *likeRepository) Exists(userID, audioID uint) (bool, error) {
	var count int64
	err := r.db.Model(&likeModel.Like{}).Where("user_id = ? AND audio_id = ?", userID, audioID).Count(&count).Error
	return count > 0, err
}
