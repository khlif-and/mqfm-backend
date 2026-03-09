package mysql

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type userLocationRepo struct {
	db *gorm.DB
}

func NewUserLocationRepository(db *gorm.DB) port.UserLocationRepository {
	return &userLocationRepo{db: db}
}

func (r *userLocationRepo) Upsert(location *entity.UserLocation) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"latitude", "longitude", "city"}),
	}).Create(location).Error
}

func (r *userLocationRepo) FindByUser(userID uint) (*entity.UserLocation, error) {
	var location entity.UserLocation
	err := r.db.Where("user_id = ?", userID).First(&location).Error
	if err != nil {
		return nil, err
	}
	return &location, nil
}

func (r *userLocationRepo) FindAll() ([]entity.UserLocation, error) {
	var locations []entity.UserLocation
	err := r.db.Find(&locations).Error
	return locations, err
}
