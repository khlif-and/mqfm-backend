package mysql

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type userPreferenceRepo struct {
	db *gorm.DB
}

func NewUserPreferenceRepository(db *gorm.DB) port.UserPreferenceRepository {
	return &userPreferenceRepo{db: db}
}

func (r *userPreferenceRepo) Upsert(pref *entity.UserPreference) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"playback_speed", "sleep_timer_minutes", "auto_download_wifi"}),
	}).Create(pref).Error
}

func (r *userPreferenceRepo) FindByUser(userID uint) (*entity.UserPreference, error) {
	var pref entity.UserPreference
	err := r.db.Where("user_id = ?", userID).First(&pref).Error
	if err != nil {
		return nil, err
	}
	return &pref, nil
}
