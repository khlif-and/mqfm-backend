package mysql

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type audioProgressRepo struct {
	db *gorm.DB
}

func NewAudioProgressRepository(db *gorm.DB) port.AudioProgressRepository {
	return &audioProgressRepo{db: db}
}

func (r *audioProgressRepo) Upsert(progress *entity.AudioProgress) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "audio_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"last_position", "duration", "percentage", "completed"}),
	}).Create(progress).Error
}

func (r *audioProgressRepo) FindByUserAndAudio(userID, audioID uint) (*entity.AudioProgress, error) {
	var progress entity.AudioProgress
	err := r.db.Where("user_id = ? AND audio_id = ?", userID, audioID).Preload("Audio").First(&progress).Error
	if err != nil {
		return nil, err
	}
	return &progress, nil
}

func (r *audioProgressRepo) FindByUser(userID uint) ([]entity.AudioProgress, error) {
	var progress []entity.AudioProgress
	err := r.db.Where("user_id = ?", userID).Preload("Audio").Order("updated_at DESC").Find(&progress).Error
	return progress, err
}

func (r *audioProgressRepo) FindCompletedByUser(userID uint) ([]entity.AudioProgress, error) {
	var progress []entity.AudioProgress
	err := r.db.Where("user_id = ? AND completed = ?", userID, true).Preload("Audio").Find(&progress).Error
	return progress, err
}

func (r *audioProgressRepo) Delete(userID, audioID uint) error {
	return r.db.Where("user_id = ? AND audio_id = ?", userID, audioID).Delete(&entity.AudioProgress{}).Error
}
