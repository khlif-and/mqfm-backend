package admin

import (
	"errors"
	"gorm.io/gorm"
	audioModel "mqfm-backend/internal/models/podcast/audio/admin"
)

type AudioRepository interface {
	FindAll() ([]audioModel.Audio, error)
	FindByID(id uint) (*audioModel.Audio, error)
	Create(audio *audioModel.Audio) error
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	Search(query string) ([]audioModel.Audio, error)
}

type audioRepository struct {
	db *gorm.DB
}

func NewAudioRepository(db *gorm.DB) AudioRepository {
	return &audioRepository{db: db}
}

func (r *audioRepository) FindAll() ([]audioModel.Audio, error) {
	var audios []audioModel.Audio
	err := r.db.Preload("Category").Find(&audios).Error
	return audios, err
}

func (r *audioRepository) FindByID(id uint) (*audioModel.Audio, error) {
	var audio audioModel.Audio
	err := r.db.Preload("Category").First(&audio, id).Error
	if err != nil {
		return nil, err
	}
	return &audio, nil
}

func (r *audioRepository) Create(audio *audioModel.Audio) error {
	return r.db.Create(audio).Error
}

func (r *audioRepository) Update(id uint, updates map[string]interface{}) error {
	result := r.db.Model(&audioModel.Audio{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no records updated")
	}
	return nil
}

func (r *audioRepository) Delete(id uint) error {
	result := r.db.Delete(&audioModel.Audio{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("audio not found")
	}
	return nil
}

func (r *audioRepository) Search(query string) ([]audioModel.Audio, error) {
	var audios []audioModel.Audio
	err := r.db.Preload("Category").Where("title LIKE ? OR artist LIKE ?", "%"+query+"%", "%"+query+"%").Find(&audios).Error
	return audios, err
}
