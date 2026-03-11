package historymock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type MockHistoryRepository struct {
	UpsertFn            func(history *entity.History) error
	FindByUserFn        func(userID uint) ([]entity.History, error)
	FindByUsersFn       func(userIDs []uint) ([]entity.History, error)
	DeleteByUserAudioFn func(userID, audioID uint) error
	DeleteAllFn         func(userID uint) error
	CountByAudioFn      func(audioID uint) (int64, error)
	FindFrequentFn      func(userID uint, minPlays int, limit int) ([]entity.History, error)
	AggregateFn         func() (map[uint]int64, error)
}

func (m *MockHistoryRepository) Upsert(history *entity.History) error { return m.UpsertFn(history) }
func (m *MockHistoryRepository) FindByUser(userID uint) ([]entity.History, error) {
	return m.FindByUserFn(userID)
}
func (m *MockHistoryRepository) FindByUsers(userIDs []uint) ([]entity.History, error) {
	return m.FindByUsersFn(userIDs)
}
func (m *MockHistoryRepository) DeleteByUserAndAudio(userID, audioID uint) error {
	return m.DeleteByUserAudioFn(userID, audioID)
}
func (m *MockHistoryRepository) DeleteAllByUser(userID uint) error { return m.DeleteAllFn(userID) }
func (m *MockHistoryRepository) CountByAudio(audioID uint) (int64, error) {
	return m.CountByAudioFn(audioID)
}
func (m *MockHistoryRepository) FindFrequentByUser(userID uint, minPlays int, limit int) ([]entity.History, error) {
	return m.FindFrequentFn(userID, minPlays, limit)
}
func (m *MockHistoryRepository) AggregatePlayCounts() (map[uint]int64, error) {
	return m.AggregateFn()
}

var _ port.HistoryRepository = (*MockHistoryRepository)(nil)
