package service

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
)

type historyService struct {
	repo port.HistoryRepository
}

func NewHistoryService(repo port.HistoryRepository) port.HistoryService {
	return &historyService{repo: repo}
}

func (s *historyService) RecordPlay(userID uint, req request.HistoryRequest) error {
	history := entity.History{
		UserID:  userID,
		AudioID: req.AudioID,
	}
	return s.repo.Upsert(&history)
}

func (s *historyService) GetHistory(userID uint) ([]entity.History, error) {
	return s.repo.FindByUser(userID)
}

func (s *historyService) DeleteHistory(userID, audioID uint) error {
	return s.repo.DeleteByUserAndAudio(userID, audioID)
}

func (s *historyService) ClearHistory(userID uint) error {
	return s.repo.DeleteAllByUser(userID)
}
