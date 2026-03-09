package service

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/internal/shared/logger"
)

type notificationService struct {
	repo port.NotificationRepository
}

func NewNotificationService(repo port.NotificationRepository) port.NotificationService {
	return &notificationService{repo: repo}
}

func (s *notificationService) GetByUser(userID uint, page, limit int) ([]entity.Notification, error) {
	offset := (page - 1) * limit
	return s.repo.FindByUser(userID, limit, offset)
}

func (s *notificationService) MarkAsRead(id, userID uint) error {
	return s.repo.MarkAsRead(id, userID)
}

func (s *notificationService) MarkAllAsRead(userID uint) error {
	return s.repo.MarkAllAsRead(userID)
}

func (s *notificationService) CountUnread(userID uint) (int64, error) {
	return s.repo.CountUnread(userID)
}

func (s *notificationService) GetSetting(userID uint) (*entity.NotificationSetting, error) {
	setting, err := s.repo.GetSetting(userID)
	if err != nil {
		defaultSetting := &entity.NotificationSetting{
			UserID:        userID,
			DailyReminder: true,
			NewContent:    true,
			EventReminder: true,
		}
		_ = s.repo.UpsertSetting(defaultSetting)
		return defaultSetting, nil
	}
	return setting, nil
}

func (s *notificationService) UpdateSetting(userID uint, req request.UpdateNotificationSettingRequest) (*entity.NotificationSetting, error) {
	setting, _ := s.GetSetting(userID)
	if setting == nil {
		setting = &entity.NotificationSetting{UserID: userID, DailyReminder: true, NewContent: true, EventReminder: true}
	}

	if req.DailyReminder != nil {
		setting.DailyReminder = *req.DailyReminder
	}
	if req.NewContent != nil {
		setting.NewContent = *req.NewContent
	}
	if req.EventReminder != nil {
		setting.EventReminder = *req.EventReminder
	}

	if err := s.repo.UpsertSetting(setting); err != nil {
		return nil, err
	}
	return setting, nil
}

func (s *notificationService) NotifyNewAudio(audio *entity.Audio) error {
	userIDs, err := s.repo.FindUserIDsWithSetting("new_content", true)
	if err != nil {
		return err
	}

	notifications := make([]entity.Notification, 0, len(userIDs))
	for _, uid := range userIDs {
		notifications = append(notifications, entity.Notification{
			UserID:      uid,
			Title:       "Kajian Baru: " + audio.Title,
			Body:        "Kajian baru dari " + audio.Artist + " telah tersedia",
			Type:        "new_audio",
			ReferenceID: audio.ID,
		})
	}

	if len(notifications) > 0 {
		go func() {
			if err := s.repo.BulkCreate(notifications); err != nil {
				logger.Error("failed to create notifications: " + err.Error())
			}
		}()
	}
	return nil
}

func (s *notificationService) NotifyDailyReminder() error {
	userIDs, err := s.repo.FindUserIDsWithSetting("daily_reminder", true)
	if err != nil {
		return err
	}

	notifications := make([]entity.Notification, 0, len(userIDs))
	for _, uid := range userIDs {
		notifications = append(notifications, entity.Notification{
			UserID: uid,
			Title:  "Pengingat Harian",
			Body:   "Sudah dengar kajian hari ini? Yuk dengarkan kajian untuk menambah ilmu!",
			Type:   "reminder",
		})
	}

	if len(notifications) > 0 {
		return s.repo.BulkCreate(notifications)
	}
	return nil
}

func (s *notificationService) NotifyEvent(event *entity.Event) error {
	userIDs, err := s.repo.FindUserIDsWithSetting("event_reminder", true)
	if err != nil {
		return err
	}

	notifications := make([]entity.Notification, 0, len(userIDs))
	for _, uid := range userIDs {
		notifications = append(notifications, entity.Notification{
			UserID:      uid,
			Title:       "Event: " + event.Title,
			Body:        "Jangan lewatkan event kajian: " + event.Title,
			Type:        "event",
			ReferenceID: event.ID,
		})
	}

	if len(notifications) > 0 {
		go func() {
			if err := s.repo.BulkCreate(notifications); err != nil {
				logger.Error("failed to create event notifications: " + err.Error())
			}
		}()
	}
	return nil
}
