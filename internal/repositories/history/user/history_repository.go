package user

import (
	"time"

	"gorm.io/gorm"

	historyModel "mqfm-backend/internal/models/history/user"
)

type HistoryRepository interface {
	Upsert(history *historyModel.History) error
	FindByUser(userID uint) ([]historyModel.History, error)
	DeleteByUserAndAudio(userID, audioID uint) error
	DeleteAllByUser(userID uint) error
}

type historyRepository struct {
	db *gorm.DB
}

func NewHistoryRepository(db *gorm.DB) HistoryRepository {
	return &historyRepository{db: db}
}

func (r *historyRepository) Upsert(history *historyModel.History) error {
	var existing historyModel.History
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

func (r *historyRepository) FindByUser(userID uint) ([]historyModel.History, error) {
	var histories []historyModel.History
	err := r.db.Preload("Audio").Where("user_id = ?", userID).Order("played_at DESC").Find(&histories).Error
	return histories, err
}

func (r *historyRepository) DeleteByUserAndAudio(userID, audioID uint) error {
	result := r.db.Where("user_id = ? AND audio_id = ?", userID, audioID).Delete(&historyModel.History{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *historyRepository) DeleteAllByUser(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&historyModel.History{}).Error
}
