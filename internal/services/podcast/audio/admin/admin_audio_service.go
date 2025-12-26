package admin

import (
	"errors"

	"gorm.io/gorm"

	audioModel "mqfm-backend/internal/models/podcast/audio/admin"

)

type AdminAudioService struct {
	db *gorm.DB
}

func NewAdminAudioService(db *gorm.DB) *AdminAudioService {
	return &AdminAudioService{db: db}
}

func (s *AdminAudioService) Create(audio *audioModel.Audio) error {
	return s.db.Create(audio).Error
}

func (s *AdminAudioService) FindAll() ([]audioModel.Audio, error) {
	var audios []audioModel.Audio
	if err := s.db.Find(&audios).Error; err != nil {
		return nil, err
	}
	return audios, nil
}

func (s *AdminAudioService) FindByID(id uint) (*audioModel.Audio, error) {
	var audio audioModel.Audio
	if err := s.db.First(&audio, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("audio not found")
		}
		return nil, err
	}
	return &audio, nil
}

func (s *AdminAudioService) Update(id uint, updates map[string]interface{}) (*audioModel.Audio, error) {
	if err := s.db.Model(&audioModel.Audio{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}

	var updatedAudio audioModel.Audio
	if err := s.db.First(&updatedAudio, id).Error; err != nil {
		return nil, err
	}

	return &updatedAudio, nil
}

func (s *AdminAudioService) Delete(id uint) error {
	var audio audioModel.Audio
	if err := s.db.First(&audio, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("audio not found")
		}
		return err
	}

	return s.db.Delete(&audio).Error
}

func (s *AdminAudioService) Search(query string) ([]audioModel.Audio, error) {
	var audios []audioModel.Audio
	// Mencari berdasarkan Title yang mengandung kata kunci (query)
	// Menggunakan query LIKE %...%
	if err := s.db.Where("title LIKE ?", "%"+query+"%").Find(&audios).Error; err != nil {
		return nil, err
	}
	return audios, nil
}