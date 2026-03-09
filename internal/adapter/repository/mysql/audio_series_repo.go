package mysql

import (
	"gorm.io/gorm"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type audioSeriesRepo struct {
	db *gorm.DB
}

func NewAudioSeriesRepository(db *gorm.DB) port.AudioSeriesRepository {
	return &audioSeriesRepo{db: db}
}

func (r *audioSeriesRepo) Create(series *entity.AudioSeries) error {
	return r.db.Create(series).Error
}

func (r *audioSeriesRepo) FindAll() ([]entity.AudioSeries, error) {
	var series []entity.AudioSeries
	err := r.db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_num ASC").Preload("Audio")
	}).Find(&series).Error
	return series, err
}

func (r *audioSeriesRepo) FindByID(id uint) (*entity.AudioSeries, error) {
	var series entity.AudioSeries
	err := r.db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_num ASC").Preload("Audio")
	}).First(&series, id).Error
	if err != nil {
		return nil, err
	}
	return &series, nil
}

func (r *audioSeriesRepo) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&entity.AudioSeries{}).Where("id = ?", id).Updates(updates).Error
}

func (r *audioSeriesRepo) Delete(id uint) error {
	return r.db.Delete(&entity.AudioSeries{}, id).Error
}

func (r *audioSeriesRepo) AddItem(item *entity.AudioSeriesItem) error {
	return r.db.Create(item).Error
}

func (r *audioSeriesRepo) RemoveItem(seriesID, audioID uint) error {
	return r.db.Where("series_id = ? AND audio_id = ?", seriesID, audioID).Delete(&entity.AudioSeriesItem{}).Error
}

func (r *audioSeriesRepo) FindItems(seriesID uint) ([]entity.AudioSeriesItem, error) {
	var items []entity.AudioSeriesItem
	err := r.db.Where("series_id = ?", seriesID).Preload("Audio").Order("order_num ASC").Find(&items).Error
	return items, err
}

func (r *audioSeriesRepo) Search(query string) ([]entity.AudioSeries, error) {
	var series []entity.AudioSeries
	err := r.db.Where("title LIKE ? OR artist LIKE ?", "%"+query+"%", "%"+query+"%").Find(&series).Error
	return series, err
}
