package mysql

import (
	"errors"

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
