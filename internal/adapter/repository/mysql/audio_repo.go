package mysql

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type audioRepo struct {
	db *gorm.DB
}

func NewAudioRepository(db *gorm.DB) port.AudioRepository {
	return &audioRepo{db: db}
}

func (r *audioRepo) FindAll() ([]entity.Audio, error) {
	var audios []entity.Audio
	err := r.db.Find(&audios).Error
	return audios, err
}

func (r *audioRepo) FindByID(id uint) (*entity.Audio, error) {
	var audio entity.Audio
	if err := r.db.First(&audio, id).Error; err != nil {
		return nil, err
	}
	return &audio, nil
}

func (r *audioRepo) Create(audio *entity.Audio) error {
	return r.db.Create(audio).Error
}

func (r *audioRepo) Update(id uint, updates map[string]interface{}) error {
	result := r.db.Model(&entity.Audio{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no records updated")
	}
	return nil
}

func (r *audioRepo) Delete(id uint) error {
	result := r.db.Delete(&entity.Audio{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("audio not found")
	}
	return nil
}

func (r *audioRepo) Search(query string) ([]entity.Audio, error) {
	var audios []entity.Audio
	err := r.db.Where("title LIKE ? OR artist LIKE ?", "%"+query+"%", "%"+query+"%").Find(&audios).Error
	return audios, err
}

func (r *audioRepo) FindByIDs(ids []uint) ([]entity.Audio, error) {
	var audios []entity.Audio
	if len(ids) == 0 {
		return audios, nil
	}
	err := r.db.Where("id IN ?", ids).Find(&audios).Error
	return audios, err
}

func (r *audioRepo) FindByArtist(artist string, limit int) ([]entity.Audio, error) {
	var audios []entity.Audio
	err := r.db.Where("artist LIKE ? AND status = 'active'", "%"+artist+"%").Limit(limit).Find(&audios).Error
	return audios, err
}

func (r *audioRepo) FindByCategoryID(categoryID uint, limit int) ([]entity.Audio, error) {
	var audios []entity.Audio
	err := r.db.Where("category_id = ? AND status = 'active'", categoryID).Limit(limit).Find(&audios).Error
	return audios, err
}

func (r *audioRepo) FindAllActive() ([]entity.Audio, error) {
	var audios []entity.Audio
	err := r.db.Where("status = 'active'").Find(&audios).Error
	return audios, err
}

func (r *audioRepo) FindBySeriesID(seriesID uint) ([]entity.Audio, error) {
	var audios []entity.Audio
	err := r.db.Where("series_id = ? AND status = 'active'", seriesID).Find(&audios).Error
	return audios, err
}

func (r *audioRepo) FindNewByArtists(artists []string, since time.Time) ([]entity.Audio, error) {
	var audios []entity.Audio
	if len(artists) == 0 {
		return audios, nil
	}
	err := r.db.Where("artist IN ? AND status = 'active' AND created_at > ?", artists, since).
		Order("created_at DESC").Find(&audios).Error
	return audios, err
}

func (r *audioRepo) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&entity.Audio{}).Count(&count).Error
	return count, err
}
