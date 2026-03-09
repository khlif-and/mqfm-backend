package service

import (
	"errors"
	"mime/multipart"
	"time"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/internal/shared/helper"
	"mqfm-backend/internal/shared/logger"
)

type eventService struct {
	repo        port.EventRepository
	notifSvc    port.NotificationService
}

func NewEventService(repo port.EventRepository, notifSvc port.NotificationService) port.EventService {
	return &eventService{repo: repo, notifSvc: notifSvc}
}

func (s *eventService) Create(req request.CreateEventRequest, file *multipart.FileHeader) (*entity.Event, error) {
	eventDate, err := time.Parse("2006-01-02 15:04:05", req.EventDate)
	if err != nil {
		return nil, errors.New("invalid event date format, use YYYY-MM-DD HH:MM:SS")
	}

	var imagePath string
	if file != nil {
		filename := helper.GenerateUniqueFilename(file.Filename)
		path := "uploads/thumbnails/" + filename
		if err := helper.SaveUploadedFile(file, path); err != nil {
			logger.Error("failed to save event image")
		} else {
			imagePath = path
		}
	}

	event := &entity.Event{
		Title:       req.Title,
		Description: req.Description,
		EventDate:   eventDate,
		Location:    req.Location,
		Image:       imagePath,
	}

	if err := s.repo.Create(event); err != nil {
		return nil, err
	}

	go s.notifSvc.NotifyEvent(event)

	return event, nil
}

func (s *eventService) FindAll() ([]entity.Event, error) {
	return s.repo.FindAll()
}

func (s *eventService) FindByID(id uint) (*entity.Event, error) {
	return s.repo.FindByID(id)
}

func (s *eventService) Update(id uint, req request.UpdateEventRequest, file *multipart.FileHeader) (*entity.Event, error) {
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.EventDate != "" {
		eventDate, err := time.Parse("2006-01-02 15:04:05", req.EventDate)
		if err != nil {
			return nil, errors.New("invalid event date format")
		}
		updates["event_date"] = eventDate
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if file != nil {
		filename := helper.GenerateUniqueFilename(file.Filename)
		path := "uploads/thumbnails/" + filename
		if err := helper.SaveUploadedFile(file, path); err == nil {
			updates["image"] = path
		}
	}

	if err := s.repo.Update(id, updates); err != nil {
		return nil, err
	}
	return s.repo.FindByID(id)
}

func (s *eventService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *eventService) GetUpcoming(limit int) ([]entity.Event, error) {
	return s.repo.FindUpcoming(limit)
}

func (s *eventService) RSVP(userID, eventID uint) error {
	exists, err := s.repo.ExistsRSVP(userID, eventID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already RSVP'd")
	}

	_, err = s.repo.FindByID(eventID)
	if err != nil {
		return errors.New("event not found")
	}

	return s.repo.CreateRSVP(&entity.EventRSVP{
		UserID:  userID,
		EventID: eventID,
	})
}

func (s *eventService) CancelRSVP(userID, eventID uint) error {
	return s.repo.DeleteRSVP(userID, eventID)
}

func (s *eventService) GetUserRSVPs(userID uint) ([]entity.EventRSVP, error) {
	return s.repo.FindRSVPsByUser(userID)
}

func (s *eventService) GetRSVPCount(eventID uint) (int64, error) {
	return s.repo.CountRSVP(eventID)
}
