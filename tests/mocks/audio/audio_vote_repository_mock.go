package audiomock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type MockAudioVoteRepository struct {
	CreateFn            func(vote *entity.AudioVote) error
	DeleteFn            func(userID, audioID uint) error
	ExistsFn            func(userID, audioID uint) (bool, error)
	CountByAudioFn      func(audioID uint) (int64, error)
	CountWeeklyFn       func(audioID uint) (int64, error)
	CountMonthlyFn      func(audioID uint) (int64, error)
	FindVotedAudioIDsFn func(userID uint) ([]uint, error)
}

func (m *MockAudioVoteRepository) Create(vote *entity.AudioVote) error { return m.CreateFn(vote) }
func (m *MockAudioVoteRepository) Delete(userID, audioID uint) error {
	return m.DeleteFn(userID, audioID)
}
func (m *MockAudioVoteRepository) Exists(userID, audioID uint) (bool, error) {
	return m.ExistsFn(userID, audioID)
}
func (m *MockAudioVoteRepository) CountByAudio(audioID uint) (int64, error) {
	return m.CountByAudioFn(audioID)
}
func (m *MockAudioVoteRepository) CountWeeklyByAudio(audioID uint) (int64, error) {
	return m.CountWeeklyFn(audioID)
}
func (m *MockAudioVoteRepository) CountMonthlyByAudio(audioID uint) (int64, error) {
	return m.CountMonthlyFn(audioID)
}
func (m *MockAudioVoteRepository) FindVotedAudioIDs(userID uint) ([]uint, error) {
	return m.FindVotedAudioIDsFn(userID)
}

var _ port.AudioVoteRepository = (*MockAudioVoteRepository)(nil)
