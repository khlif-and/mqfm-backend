package mysql

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type smartResumeRepo struct {
	db *gorm.DB
}

func NewSmartResumeRepository(db *gorm.DB) port.SmartResumeRepository {
	return &smartResumeRepo{db: db}
}

func (r *smartResumeRepo) Upsert(resume *entity.SmartResume) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"audio_id", "playlist_id", "position_seconds"}),
	}).Create(resume).Error
}

func (r *smartResumeRepo) FindByUser(userID uint) (*entity.SmartResume, error) {
	var resume entity.SmartResume
	err := r.db.Where("user_id = ?", userID).Preload("Audio").First(&resume).Error
	if err != nil {
		return nil, err
	}
	return &resume, nil
}
