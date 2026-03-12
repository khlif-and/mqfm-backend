package mysql

import (
	"gorm.io/gorm"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type radioRepo struct {
	db *gorm.DB
}

func NewRadioRepository(db *gorm.DB) port.RadioRepository {
	return &radioRepo{db: db}
}

func (r *radioRepo) Create(radio *entity.Radio) error {
	return r.db.Create(radio).Error
}

func (r *radioRepo) FindAll() ([]entity.Radio, error) {
	var radios []entity.Radio
	err := r.db.Preload("Audios").Order("created_at DESC").Find(&radios).Error
	return radios, err
}

func (r *radioRepo) FindActive() ([]entity.Radio, error) {
	var radios []entity.Radio
	err := r.db.Preload("Audios").Where("is_active = ?", true).Order("created_at DESC").Find(&radios).Error
	return radios, err
}

func (r *radioRepo) FindByID(id uint) (*entity.Radio, error) {
	var radio entity.Radio
	err := r.db.Preload("Audios").First(&radio, id).Error
	return &radio, err
}

func (r *radioRepo) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&entity.Radio{}).Where("id = ?", id).Updates(updates).Error
}

func (r *radioRepo) Delete(id uint) error {
	return r.db.Delete(&entity.Radio{}, id).Error
}

func (r *radioRepo) AddAudio(radioID, audioID uint, orderNum int) error {
	return r.db.Create(&entity.RadioAudio{
		RadioID:  radioID,
		AudioID:  audioID,
		OrderNum: orderNum,
	}).Error
}

func (r *radioRepo) RemoveAudio(radioID, audioID uint) error {
	return r.db.Where("radio_id = ? AND audio_id = ?", radioID, audioID).Delete(&entity.RadioAudio{}).Error
}

func (r *radioRepo) FindAudios(radioID uint) ([]*entity.Audio, error) {
	var audios []*entity.Audio
	err := r.db.Joins("JOIN radio_audios ON radio_audios.audio_id = audios.id").
		Where("radio_audios.radio_id = ?", radioID).
		Order("radio_audios.order_num ASC").
		Find(&audios).Error
	return audios, err
}

func (r *radioRepo) CountAudios(radioID uint) (int, error) {
	var count int64
	err := r.db.Model(&entity.RadioAudio{}).Where("radio_id = ?", radioID).Count(&count).Error
	return int(count), err
}
