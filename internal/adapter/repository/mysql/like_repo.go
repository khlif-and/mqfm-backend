package mysql

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type likeRepo struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) port.LikeRepository {
	return &likeRepo{db: db}
}

func (r *likeRepo) Create(like *entity.Like) error {
	return r.db.Create(like).Error
}

func (r *likeRepo) Delete(userID, audioID uint) error {
	result := r.db.Where("user_id = ? AND audio_id = ?", userID, audioID).Delete(&entity.Like{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("like not found")
	}
	return nil
}

func (r *likeRepo) FindByUser(userID uint) ([]entity.Like, error) {
	var likes []entity.Like
	err := r.db.Preload("Audio").Where("user_id = ?", userID).Find(&likes).Error
	return likes, err
}

func (r *likeRepo) Exists(userID, audioID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Like{}).Where("user_id = ? AND audio_id = ?", userID, audioID).Count(&count).Error
	return count > 0, err
}

type historyRepo struct {
	db *gorm.DB
}

func NewHistoryRepository(db *gorm.DB) port.HistoryRepository {
	return &historyRepo{db: db}
}

func (r *historyRepo) Upsert(history *entity.History) error {
	var existing entity.History
	err := r.db.Where("user_id = ? AND audio_id = ?", history.UserID, history.AudioID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		history.PlayCount = 1
		history.PlayedAt = time.Now()
		return r.db.Create(history).Error
	}
	if err != nil {
		return err
	}

	return r.db.Model(&existing).Updates(map[string]interface{}{
		"play_count": existing.PlayCount + 1,
		"played_at":  time.Now(),
	}).Error
}

func (r *historyRepo) FindByUser(userID uint) ([]entity.History, error) {
	var histories []entity.History
	err := r.db.Preload("Audio").Where("user_id = ?", userID).Order("played_at DESC").Find(&histories).Error
	return histories, err
}

func (r *historyRepo) DeleteByUserAndAudio(userID, audioID uint) error {
	result := r.db.Where("user_id = ? AND audio_id = ?", userID, audioID).Delete(&entity.History{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *historyRepo) DeleteAllByUser(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&entity.History{}).Error
}
