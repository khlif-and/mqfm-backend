package statmock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type MockListeningStatRepository struct {
	CreateFn            func(stat *entity.ListeningStat) error
	GetWeeklySummaryFn  func(userID uint) (int, error)
	GetMonthlySummaryFn func(userID uint) (int, error)
	GetTopCategoriesFn  func(userID uint, limit int) ([]port.CategoryStat, error)
	GetTopArtistsFn     func(userID uint, limit int) ([]port.ArtistStat, error)
	GetDailySummaryFn   func(userID uint, days int) ([]port.DailyStat, error)
}

func (m *MockListeningStatRepository) Create(stat *entity.ListeningStat) error {
	return m.CreateFn(stat)
}
func (m *MockListeningStatRepository) GetWeeklySummary(userID uint) (int, error) {
	return m.GetWeeklySummaryFn(userID)
}
func (m *MockListeningStatRepository) GetMonthlySummary(userID uint) (int, error) {
	return m.GetMonthlySummaryFn(userID)
}
func (m *MockListeningStatRepository) GetTopCategories(userID uint, limit int) ([]port.CategoryStat, error) {
	return m.GetTopCategoriesFn(userID, limit)
}
func (m *MockListeningStatRepository) GetTopArtists(userID uint, limit int) ([]port.ArtistStat, error) {
	return m.GetTopArtistsFn(userID, limit)
}
func (m *MockListeningStatRepository) GetDailySummary(userID uint, days int) ([]port.DailyStat, error) {
	return m.GetDailySummaryFn(userID, days)
}

var _ port.ListeningStatRepository = (*MockListeningStatRepository)(nil)
