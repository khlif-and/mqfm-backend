package historymock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
)

type MockHistoryService struct {
	RecordPlayFn    func(userID uint, req request.HistoryRequest) error
	GetHistoryFn    func(userID uint) ([]entity.History, error)
	DeleteHistoryFn func(userID, audioID uint) error
	ClearHistoryFn  func(userID uint) error
}

func (m *MockHistoryService) RecordPlay(userID uint, req request.HistoryRequest) error {
	return m.RecordPlayFn(userID, req)
}
func (m *MockHistoryService) GetHistory(userID uint) ([]entity.History, error) {
	return m.GetHistoryFn(userID)
}
func (m *MockHistoryService) DeleteHistory(userID, audioID uint) error {
	return m.DeleteHistoryFn(userID, audioID)
}
func (m *MockHistoryService) ClearHistory(userID uint) error { return m.ClearHistoryFn(userID) }

var _ port.HistoryService = (*MockHistoryService)(nil)
