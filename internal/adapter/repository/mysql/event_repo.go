package mysql

import (
	"errors"

	"gorm.io/gorm"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type eventRepo struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) port.EventRepository {
	return &eventRepo{db: db}
}

func (r *eventRepo) Create(event *entity.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepo) FindAll() ([]entity.Event, error) {
	var events []entity.Event
	err := r.db.Order("event_date DESC").Find(&events).Error
	return events, err
}

func (r *eventRepo) FindByID(id uint) (*entity.Event, error) {
	var event entity.Event
	if err := r.db.First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepo) Update(id uint, updates map[string]interface{}) error {
	result := r.db.Model(&entity.Event{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("event not found")
	}
	return nil
}

func (r *eventRepo) Delete(id uint) error {
	result := r.db.Delete(&entity.Event{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("event not found")
	}
	return nil
}

func (r *eventRepo) FindUpcoming(limit int) ([]entity.Event, error) {
	var events []entity.Event
	err := r.db.Where("event_date >= NOW()").Order("event_date ASC").Limit(limit).Find(&events).Error
	return events, err
}

func (r *eventRepo) CreateRSVP(rsvp *entity.EventRSVP) error {
	return r.db.Create(rsvp).Error
}

func (r *eventRepo) DeleteRSVP(userID, eventID uint) error {
	return r.db.Where("user_id = ? AND event_id = ?", userID, eventID).Delete(&entity.EventRSVP{}).Error
}

func (r *eventRepo) FindRSVPsByEvent(eventID uint) ([]entity.EventRSVP, error) {
	var rsvps []entity.EventRSVP
	err := r.db.Where("event_id = ?", eventID).Find(&rsvps).Error
	return rsvps, err
}

func (r *eventRepo) FindRSVPsByUser(userID uint) ([]entity.EventRSVP, error) {
	var rsvps []entity.EventRSVP
	err := r.db.Where("user_id = ?", userID).Preload("Event").Find(&rsvps).Error
	return rsvps, err
}

func (r *eventRepo) ExistsRSVP(userID, eventID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.EventRSVP{}).Where("user_id = ? AND event_id = ?", userID, eventID).Count(&count).Error
	return count > 0, err
}

func (r *eventRepo) CountRSVP(eventID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.EventRSVP{}).Where("event_id = ?", eventID).Count(&count).Error
	return count, err
}
