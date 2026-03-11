package audiomock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type MockAudioRankingRepository struct {
	UpsertFn         func(ranking *entity.AudioRanking) error
	BulkUpsertFn     func(rankings []entity.AudioRanking) error
	FindTopWeeklyFn  func(limit int) ([]entity.AudioRanking, error)
	FindTopMonthlyFn func(limit int) ([]entity.AudioRanking, error)
	FindByAudioFn    func(audioID uint) (*entity.AudioRanking, error)
	CountAllFn       func() (int64, error)
	FindAllFn        func(limit, offset int) ([]entity.AudioRanking, error)
}

func (m *MockAudioRankingRepository) Upsert(ranking *entity.AudioRanking) error {
	return m.UpsertFn(ranking)
}
func (m *MockAudioRankingRepository) BulkUpsert(rankings []entity.AudioRanking) error {
	return m.BulkUpsertFn(rankings)
}
func (m *MockAudioRankingRepository) FindTopWeekly(limit int) ([]entity.AudioRanking, error) {
	return m.FindTopWeeklyFn(limit)
}
func (m *MockAudioRankingRepository) FindTopMonthly(limit int) ([]entity.AudioRanking, error) {
	return m.FindTopMonthlyFn(limit)
}
func (m *MockAudioRankingRepository) FindByAudioID(audioID uint) (*entity.AudioRanking, error) {
	return m.FindByAudioFn(audioID)
}
func (m *MockAudioRankingRepository) CountAll() (int64, error) { return m.CountAllFn() }
func (m *MockAudioRankingRepository) FindAll(limit, offset int) ([]entity.AudioRanking, error) {
	return m.FindAllFn(limit, offset)
}

var _ port.AudioRankingRepository = (*MockAudioRankingRepository)(nil)
