package mysql

import (
	"gorm.io/gorm"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type downloadRepo struct {
	db *gorm.DB
}

func NewDownloadRepository(db *gorm.DB) port.DownloadRepository {
	return &downloadRepo{db: db}
}

func (r *downloadRepo) Create(download *entity.Download) error {
	return r.db.Create(download).Error
}

func (r *downloadRepo) FindByUser(userID uint) ([]entity.Download, error) {
	var downloads []entity.Download
	err := r.db.Where("user_id = ?", userID).Preload("Audio").Order("created_at DESC").Find(&downloads).Error
	return downloads, err
}

func (r *downloadRepo) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&entity.Download{}).Error
}

func (r *downloadRepo) Exists(userID, audioID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Download{}).Where("user_id = ? AND audio_id = ?", userID, audioID).Count(&count).Error
	return count > 0, err
}

func (r *downloadRepo) SumSizeByUser(userID uint) (int64, error) {
	var total int64
	err := r.db.Model(&entity.Download{}).Where("user_id = ?", userID).
		Select("COALESCE(SUM(file_size), 0)").Scan(&total).Error
	return total, err
}
