package notificationmock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type MockNotificationRepository struct {
	CreateFn             func(notification *entity.Notification) error
	BulkCreateFn         func(notifications []entity.Notification) error
	FindByUserFn         func(userID uint, limit, offset int) ([]entity.Notification, error)
	MarkAsReadFn         func(id, userID uint) error
	MarkAllAsReadFn      func(userID uint) error
	CountUnreadFn        func(userID uint) (int64, error)
	GetSettingFn         func(userID uint) (*entity.NotificationSetting, error)
	UpsertSettingFn      func(setting *entity.NotificationSetting) error
	FindUserIDsSettingFn func(field string, value bool) ([]uint, error)
}

func (m *MockNotificationRepository) Create(notification *entity.Notification) error {
	return m.CreateFn(notification)
}
func (m *MockNotificationRepository) BulkCreate(notifications []entity.Notification) error {
	return m.BulkCreateFn(notifications)
}
func (m *MockNotificationRepository) FindByUser(userID uint, limit, offset int) ([]entity.Notification, error) {
	return m.FindByUserFn(userID, limit, offset)
}
func (m *MockNotificationRepository) MarkAsRead(id, userID uint) error {
	return m.MarkAsReadFn(id, userID)
}
func (m *MockNotificationRepository) MarkAllAsRead(userID uint) error {
	return m.MarkAllAsReadFn(userID)
}
func (m *MockNotificationRepository) CountUnread(userID uint) (int64, error) {
	return m.CountUnreadFn(userID)
}
func (m *MockNotificationRepository) GetSetting(userID uint) (*entity.NotificationSetting, error) {
	return m.GetSettingFn(userID)
}
func (m *MockNotificationRepository) UpsertSetting(setting *entity.NotificationSetting) error {
	return m.UpsertSettingFn(setting)
}
func (m *MockNotificationRepository) FindUserIDsWithSetting(field string, value bool) ([]uint, error) {
	return m.FindUserIDsSettingFn(field, value)
}

var _ port.NotificationRepository = (*MockNotificationRepository)(nil)
