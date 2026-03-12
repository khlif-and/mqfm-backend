package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"mqfm-backend/internal/app/service"
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/shared/dto/request"
	notificationmock "mqfm-backend/tests/mocks/notification"
)

// ──────────────────────── GetByUser ────────────────────────

func TestNotification_GetByUser_PaginationOffset(t *testing.T) {
	repo := &notificationmock.MockNotificationRepository{
		FindByUserFn: func(userID uint, limit, offset int) ([]entity.Notification, error) {
			assert.Equal(t, uint(1), userID)
			assert.Equal(t, 10, limit)
			assert.Equal(t, 20, offset, "page 3, limit 10 → offset = (3-1)*10 = 20")
			return []entity.Notification{{ID: 21, Title: "Test"}}, nil
		},
	}
	svc := service.NewNotificationService(repo)

	result, err := svc.GetByUser(1, 3, 10)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestNotification_GetByUser_FirstPage(t *testing.T) {
	repo := &notificationmock.MockNotificationRepository{
		FindByUserFn: func(userID uint, limit, offset int) ([]entity.Notification, error) {
			assert.Equal(t, 0, offset, "page 1 → offset = 0")
			return []entity.Notification{
				{ID: 1, Title: "Notif 1"},
				{ID: 2, Title: "Notif 2"},
			}, nil
		},
	}
	svc := service.NewNotificationService(repo)

	result, err := svc.GetByUser(1, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestNotification_GetByUser_Error(t *testing.T) {
	repo := &notificationmock.MockNotificationRepository{
		FindByUserFn: func(userID uint, limit, offset int) ([]entity.Notification, error) {
			return nil, errors.New("db error")
		},
	}
	svc := service.NewNotificationService(repo)

	result, err := svc.GetByUser(1, 1, 10)
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ──────────────────────── MarkAsRead ────────────────────────

func TestNotification_MarkAsRead_Success(t *testing.T) {
	repo := &notificationmock.MockNotificationRepository{
		MarkAsReadFn: func(id, userID uint) error {
			assert.Equal(t, uint(5), id)
			assert.Equal(t, uint(1), userID)
			return nil
		},
	}
	svc := service.NewNotificationService(repo)

	err := svc.MarkAsRead(5, 1)
	assert.NoError(t, err)
}

func TestNotification_MarkAsRead_Error(t *testing.T) {
	repo := &notificationmock.MockNotificationRepository{
		MarkAsReadFn: func(id, userID uint) error { return errors.New("not found") },
	}
	svc := service.NewNotificationService(repo)

	err := svc.MarkAsRead(999, 1)
	assert.Error(t, err)
}

// ──────────────────────── MarkAllAsRead ────────────────────────

func TestNotification_MarkAllAsRead_Success(t *testing.T) {
	repo := &notificationmock.MockNotificationRepository{
		MarkAllAsReadFn: func(userID uint) error {
			assert.Equal(t, uint(1), userID)
			return nil
		},
	}
	svc := service.NewNotificationService(repo)

	err := svc.MarkAllAsRead(1)
	assert.NoError(t, err)
}

// ──────────────────────── CountUnread ────────────────────────

func TestNotification_CountUnread_Success(t *testing.T) {
	repo := &notificationmock.MockNotificationRepository{
		CountUnreadFn: func(userID uint) (int64, error) {
			return 7, nil
		},
	}
	svc := service.NewNotificationService(repo)

	count, err := svc.CountUnread(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(7), count)
}

// ──────────────────────── GetSetting ────────────────────────

func TestNotification_GetSetting_ExistingReturned(t *testing.T) {
	existing := &entity.NotificationSetting{
		ID: 1, UserID: 1, DailyReminder: false, NewContent: true, EventReminder: false,
	}
	repo := &notificationmock.MockNotificationRepository{
		GetSettingFn: func(userID uint) (*entity.NotificationSetting, error) {
			return existing, nil
		},
	}
	svc := service.NewNotificationService(repo)

	result, err := svc.GetSetting(1)
	assert.NoError(t, err)
	assert.Equal(t, existing, result)
	assert.False(t, result.DailyReminder)
	assert.True(t, result.NewContent)
}

func TestNotification_GetSetting_CreatesDefault(t *testing.T) {
	var upserted *entity.NotificationSetting
	repo := &notificationmock.MockNotificationRepository{
		GetSettingFn: func(userID uint) (*entity.NotificationSetting, error) {
			return nil, errors.New("not found")
		},
		UpsertSettingFn: func(setting *entity.NotificationSetting) error {
			upserted = setting
			return nil
		},
	}
	svc := service.NewNotificationService(repo)

	result, err := svc.GetSetting(1)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.UserID)
	assert.True(t, result.DailyReminder, "default daily_reminder = true")
	assert.True(t, result.NewContent, "default new_content = true")
	assert.True(t, result.EventReminder, "default event_reminder = true")
	assert.NotNil(t, upserted)
}

// ──────────────────────── UpdateSetting ────────────────────────

func TestNotification_UpdateSetting_AllFields(t *testing.T) {
	daily := false
	content := false
	event := true

	repo := &notificationmock.MockNotificationRepository{
		GetSettingFn: func(userID uint) (*entity.NotificationSetting, error) {
			return &entity.NotificationSetting{
				ID: 1, UserID: 1, DailyReminder: true, NewContent: true, EventReminder: false,
			}, nil
		},
		UpsertSettingFn: func(setting *entity.NotificationSetting) error { return nil },
	}
	svc := service.NewNotificationService(repo)

	result, err := svc.UpdateSetting(1, request.UpdateNotificationSettingRequest{
		DailyReminder: &daily,
		NewContent:    &content,
		EventReminder: &event,
	})
	assert.NoError(t, err)
	assert.False(t, result.DailyReminder)
	assert.False(t, result.NewContent)
	assert.True(t, result.EventReminder)
}

func TestNotification_UpdateSetting_PartialFields(t *testing.T) {
	daily := false
	repo := &notificationmock.MockNotificationRepository{
		GetSettingFn: func(userID uint) (*entity.NotificationSetting, error) {
			return &entity.NotificationSetting{
				ID: 1, UserID: 1, DailyReminder: true, NewContent: true, EventReminder: true,
			}, nil
		},
		UpsertSettingFn: func(setting *entity.NotificationSetting) error { return nil },
	}
	svc := service.NewNotificationService(repo)

	result, err := svc.UpdateSetting(1, request.UpdateNotificationSettingRequest{
		DailyReminder: &daily,
	})
	assert.NoError(t, err)
	assert.False(t, result.DailyReminder, "updated field")
	assert.True(t, result.NewContent, "untouched field stays same")
	assert.True(t, result.EventReminder, "untouched field stays same")
}

func TestNotification_UpdateSetting_UpsertError(t *testing.T) {
	daily := false
	repo := &notificationmock.MockNotificationRepository{
		GetSettingFn: func(userID uint) (*entity.NotificationSetting, error) {
			return &entity.NotificationSetting{ID: 1, UserID: 1}, nil
		},
		UpsertSettingFn: func(setting *entity.NotificationSetting) error {
			return errors.New("upsert failed")
		},
	}
	svc := service.NewNotificationService(repo)

	result, err := svc.UpdateSetting(1, request.UpdateNotificationSettingRequest{DailyReminder: &daily})
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ──────────────────────── NotifyNewAudio ────────────────────────

func TestNotification_NotifyNewAudio_Success(t *testing.T) {
	bulkCalled := make(chan bool, 1)
	repo := &notificationmock.MockNotificationRepository{
		FindUserIDsSettingFn: func(field string, value bool) ([]uint, error) {
			assert.Equal(t, "new_content", field)
			assert.True(t, value)
			return []uint{1, 2, 3}, nil
		},
		BulkCreateFn: func(notifications []entity.Notification) error {
			assert.Len(t, notifications, 3)
			for _, n := range notifications {
				assert.Equal(t, "new_audio", n.Type)
				assert.Equal(t, uint(42), n.ReferenceID)
				assert.Contains(t, n.Title, "Tafsir Surah Al-Kahfi")
				assert.Contains(t, n.Body, "Ust. Adi Hidayat")
			}
			bulkCalled <- true
			return nil
		},
	}
	svc := service.NewNotificationService(repo)

	audio := &entity.Audio{ID: 42, Title: "Tafsir Surah Al-Kahfi", Artist: "Ust. Adi Hidayat"}
	err := svc.NotifyNewAudio(audio)
	assert.NoError(t, err)

	// Wait for async goroutine
	select {
	case <-bulkCalled:
	case <-time.After(2 * time.Second):
		t.Fatal("BulkCreate was not called within timeout")
	}
}

func TestNotification_NotifyNewAudio_NoUsers(t *testing.T) {
	repo := &notificationmock.MockNotificationRepository{
		FindUserIDsSettingFn: func(field string, value bool) ([]uint, error) {
			return []uint{}, nil
		},
	}
	svc := service.NewNotificationService(repo)

	err := svc.NotifyNewAudio(&entity.Audio{ID: 1, Title: "Test"})
	assert.NoError(t, err)
}

func TestNotification_NotifyNewAudio_FindUsersError(t *testing.T) {
	repo := &notificationmock.MockNotificationRepository{
		FindUserIDsSettingFn: func(field string, value bool) ([]uint, error) {
			return nil, errors.New("db error")
		},
	}
	svc := service.NewNotificationService(repo)

	err := svc.NotifyNewAudio(&entity.Audio{ID: 1, Title: "Test"})
	assert.Error(t, err)
}

// ──────────────────────── NotifyDailyReminder ────────────────────────

func TestNotification_NotifyDailyReminder_Success(t *testing.T) {
	repo := &notificationmock.MockNotificationRepository{
		FindUserIDsSettingFn: func(field string, value bool) ([]uint, error) {
			assert.Equal(t, "daily_reminder", field)
			return []uint{1, 2}, nil
		},
		BulkCreateFn: func(notifications []entity.Notification) error {
			assert.Len(t, notifications, 2)
			for _, n := range notifications {
				assert.Equal(t, "reminder", n.Type)
				assert.Equal(t, "Pengingat Harian", n.Title)
			}
			return nil
		},
	}
	svc := service.NewNotificationService(repo)

	err := svc.NotifyDailyReminder()
	assert.NoError(t, err)
}

func TestNotification_NotifyDailyReminder_NoUsers(t *testing.T) {
	repo := &notificationmock.MockNotificationRepository{
		FindUserIDsSettingFn: func(field string, value bool) ([]uint, error) {
			return []uint{}, nil
		},
	}
	svc := service.NewNotificationService(repo)

	err := svc.NotifyDailyReminder()
	assert.NoError(t, err)
}

func TestNotification_NotifyDailyReminder_BulkCreateError(t *testing.T) {
	repo := &notificationmock.MockNotificationRepository{
		FindUserIDsSettingFn: func(field string, value bool) ([]uint, error) {
			return []uint{1}, nil
		},
		BulkCreateFn: func(notifications []entity.Notification) error {
			return errors.New("bulk failed")
		},
	}
	svc := service.NewNotificationService(repo)

	err := svc.NotifyDailyReminder()
	assert.Error(t, err)
	assert.Equal(t, "bulk failed", err.Error())
}

// ──────────────────────── NotifyEvent ────────────────────────

func TestNotification_NotifyEvent_Success(t *testing.T) {
	bulkCalled := make(chan bool, 1)
	repo := &notificationmock.MockNotificationRepository{
		FindUserIDsSettingFn: func(field string, value bool) ([]uint, error) {
			assert.Equal(t, "event_reminder", field)
			return []uint{1, 2}, nil
		},
		BulkCreateFn: func(notifications []entity.Notification) error {
			assert.Len(t, notifications, 2)
			for _, n := range notifications {
				assert.Equal(t, "event", n.Type)
				assert.Equal(t, uint(10), n.ReferenceID)
				assert.Contains(t, n.Title, "Kajian Akbar")
			}
			bulkCalled <- true
			return nil
		},
	}
	svc := service.NewNotificationService(repo)

	event := &entity.Event{ID: 10, Title: "Kajian Akbar"}
	err := svc.NotifyEvent(event)
	assert.NoError(t, err)

	select {
	case <-bulkCalled:
	case <-time.After(2 * time.Second):
		t.Fatal("BulkCreate was not called within timeout")
	}
}

func TestNotification_NotifyEvent_NoUsers(t *testing.T) {
	repo := &notificationmock.MockNotificationRepository{
		FindUserIDsSettingFn: func(field string, value bool) ([]uint, error) {
			return []uint{}, nil
		},
	}
	svc := service.NewNotificationService(repo)

	err := svc.NotifyEvent(&entity.Event{ID: 1, Title: "Test"})
	assert.NoError(t, err)
}

func TestNotification_NotifyEvent_FindUsersError(t *testing.T) {
	repo := &notificationmock.MockNotificationRepository{
		FindUserIDsSettingFn: func(field string, value bool) ([]uint, error) {
			return nil, errors.New("db error")
		},
	}
	svc := service.NewNotificationService(repo)

	err := svc.NotifyEvent(&entity.Event{ID: 1, Title: "Test"})
	assert.Error(t, err)
}
