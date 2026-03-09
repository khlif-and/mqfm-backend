package mysql

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type notificationRepo struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) port.NotificationRepository {
	return &notificationRepo{db: db}
}

func (r *notificationRepo) Create(notification *entity.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepo) BulkCreate(notifications []entity.Notification) error {
	if len(notifications) == 0 {
		return nil
	}
	return r.db.CreateInBatches(notifications, 100).Error
}

func (r *notificationRepo) FindByUser(userID uint, limit, offset int) ([]entity.Notification, error) {
	var notifications []entity.Notification
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").Limit(limit).Offset(offset).Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepo) MarkAsRead(id, userID uint) error {
	return r.db.Model(&entity.Notification{}).Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true).Error
}

func (r *notificationRepo) MarkAllAsRead(userID uint) error {
	return r.db.Model(&entity.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

func (r *notificationRepo) CountUnread(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

func (r *notificationRepo) GetSetting(userID uint) (*entity.NotificationSetting, error) {
	var setting entity.NotificationSetting
	err := r.db.Where("user_id = ?", userID).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *notificationRepo) UpsertSetting(setting *entity.NotificationSetting) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"daily_reminder", "new_content", "event_reminder"}),
	}).Create(setting).Error
}

func (r *notificationRepo) FindUserIDsWithSetting(field string, value bool) ([]uint, error) {
	var ids []uint
	err := r.db.Model(&entity.NotificationSetting{}).Where(field+" = ?", value).Pluck("user_id", &ids).Error
	return ids, err
}
