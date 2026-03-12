package audiomock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type MockAudioScoreRepository struct {
	UpsertFn         func(score *entity.AudioScore) error
	FindTopFn        func(limit int) ([]entity.AudioScore, error)
	FindTopByLikesFn func(limit int, maxLikes int64) ([]entity.AudioScore, error)
	FindByAudioFn    func(audioID uint) (*entity.AudioScore, error)
	FindByAudiosFn   func(audioIDs []uint) ([]entity.AudioScore, error)
	DeleteAllFn      func() error
	BulkUpsertFn     func(scores []entity.AudioScore) error
}

func (m *MockAudioScoreRepository) Upsert(score *entity.AudioScore) error {
	return m.UpsertFn(score)
}
func (m *MockAudioScoreRepository) FindTopByScore(limit int) ([]entity.AudioScore, error) {
	return m.FindTopFn(limit)
}
func (m *MockAudioScoreRepository) FindTopByLikes(limit int, maxLikes int64) ([]entity.AudioScore, error) {
	return m.FindTopByLikesFn(limit, maxLikes)
}
func (m *MockAudioScoreRepository) FindByAudioID(audioID uint) (*entity.AudioScore, error) {
	return m.FindByAudioFn(audioID)
}
func (m *MockAudioScoreRepository) FindByAudioIDs(audioIDs []uint) ([]entity.AudioScore, error) {
	return m.FindByAudiosFn(audioIDs)
}
func (m *MockAudioScoreRepository) DeleteAll() error { return m.DeleteAllFn() }
func (m *MockAudioScoreRepository) BulkUpsert(scores []entity.AudioScore) error {
	return m.BulkUpsertFn(scores)
}
func (m *MockAudioScoreRepository) FindTopByWeeklyLikes(limit int) ([]entity.AudioScore, error) {
	return []entity.AudioScore{}, nil
}
func (m *MockAudioScoreRepository) FindTopByMonthlyLikes(limit int) ([]entity.AudioScore, error) {
	return []entity.AudioScore{}, nil
}
func (m *MockAudioScoreRepository) BulkUpdateWeeklyLikes(data map[uint]int64) error { return nil }
func (m *MockAudioScoreRepository) BulkUpdateMonthlyLikes(data map[uint]int64) error { return nil }

var _ port.AudioScoreRepository = (*MockAudioScoreRepository)(nil)
