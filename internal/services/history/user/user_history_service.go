package user

import (
	historyDto "mqfm-backend/internal/dto/history"
	historyModel "mqfm-backend/internal/models/history/user"
	historyRepo "mqfm-backend/internal/repositories/history/user"
)

type UserHistoryService struct {
	repo historyRepo.HistoryRepository
}

func NewUserHistoryService(repo historyRepo.HistoryRepository) *UserHistoryService {
	return &UserHistoryService{repo: repo}
}

func (s *UserHistoryService) RecordPlay(userID uint, req historyDto.HistoryRequest) error {
	history := historyModel.History{
		UserID:  userID,
		AudioID: req.AudioID,
	}
	return s.repo.Upsert(&history)
}

func (s *UserHistoryService) GetHistory(userID uint) ([]historyModel.History, error) {
	return s.repo.FindByUser(userID)
}

func (s *UserHistoryService) DeleteHistory(userID uint, audioID uint) error {
	return s.repo.DeleteByUserAndAudio(userID, audioID)
}

func (s *UserHistoryService) ClearHistory(userID uint) error {
	return s.repo.DeleteAllByUser(userID)
}
