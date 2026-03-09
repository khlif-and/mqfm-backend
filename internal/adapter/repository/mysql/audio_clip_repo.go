package mysql

import (
	"gorm.io/gorm"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type audioClipRepo struct {
	db *gorm.DB
}

func NewAudioClipRepository(db *gorm.DB) port.AudioClipRepository {
	return &audioClipRepo{db: db}
}

func (r *audioClipRepo) Create(clip *entity.AudioClip) error {
	return r.db.Create(clip).Error
}

func (r *audioClipRepo) FindByUser(userID uint) ([]entity.AudioClip, error) {
	var clips []entity.AudioClip
	err := r.db.Where("user_id = ?", userID).Preload("Audio").Order("created_at DESC").Find(&clips).Error
	return clips, err
}

func (r *audioClipRepo) FindByShareToken(token string) (*entity.AudioClip, error) {
	var clip entity.AudioClip
	err := r.db.Where("share_token = ?", token).Preload("Audio").First(&clip).Error
	if err != nil {
		return nil, err
	}
	return &clip, nil
}

func (r *audioClipRepo) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&entity.AudioClip{}).Error
}
