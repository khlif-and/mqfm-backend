package audiomock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"time"
)

type MockAudioRepository struct {
	FindAllFn          func() ([]entity.Audio, error)
	FindByIDFn         func(id uint) (*entity.Audio, error)
	CreateFn           func(audio *entity.Audio) error
	UpdateFn           func(id uint, updates map[string]interface{}) error
	DeleteFn           func(id uint) error
	SearchFn           func(query string) ([]entity.Audio, error)
	FindByIDsFn        func(ids []uint) ([]entity.Audio, error)
	FindByArtistFn     func(artist string, limit int) ([]entity.Audio, error)
	FindByCategoryFn   func(categoryID uint, limit int) ([]entity.Audio, error)
	FindAllActiveFn    func() ([]entity.Audio, error)
	FindBySeriesFn     func(seriesID uint) ([]entity.Audio, error)
	FindNewByArtistsFn func(artists []string, since time.Time) ([]entity.Audio, error)
	CountAllFn         func() (int64, error)
}

func (m *MockAudioRepository) FindAll() ([]entity.Audio, error)   { return m.FindAllFn() }
func (m *MockAudioRepository) FindByID(id uint) (*entity.Audio, error) { return m.FindByIDFn(id) }
func (m *MockAudioRepository) Create(audio *entity.Audio) error    { return m.CreateFn(audio) }
func (m *MockAudioRepository) Update(id uint, updates map[string]interface{}) error {
	return m.UpdateFn(id, updates)
}
func (m *MockAudioRepository) Delete(id uint) error               { return m.DeleteFn(id) }
func (m *MockAudioRepository) Search(query string) ([]entity.Audio, error) {
	return m.SearchFn(query)
}
func (m *MockAudioRepository) FindByIDs(ids []uint) ([]entity.Audio, error) {
	return m.FindByIDsFn(ids)
}
func (m *MockAudioRepository) FindByArtist(artist string, limit int) ([]entity.Audio, error) {
	return m.FindByArtistFn(artist, limit)
}
func (m *MockAudioRepository) FindByCategoryID(categoryID uint, limit int) ([]entity.Audio, error) {
	return m.FindByCategoryFn(categoryID, limit)
}
func (m *MockAudioRepository) FindAllActive() ([]entity.Audio, error) { return m.FindAllActiveFn() }
func (m *MockAudioRepository) FindBySeriesID(seriesID uint) ([]entity.Audio, error) {
	return m.FindBySeriesFn(seriesID)
}
func (m *MockAudioRepository) FindNewByArtists(artists []string, since time.Time) ([]entity.Audio, error) {
	return m.FindNewByArtistsFn(artists, since)
}
func (m *MockAudioRepository) CountAll() (int64, error) { return m.CountAllFn() }

var _ port.AudioRepository = (*MockAudioRepository)(nil)
