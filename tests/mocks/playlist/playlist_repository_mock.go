package playlistmock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type MockPlaylistRepository struct {
	CreateFn        func(playlist *entity.Playlist) error
	FindByUserIDFn  func(userID uint) ([]entity.Playlist, error)
	FindByIDFn      func(id uint) (*entity.Playlist, error)
	FindAllFn       func() ([]entity.Playlist, error)
	SearchFn        func(query string) ([]entity.Playlist, error)
	AddAudioFn      func(playlist *entity.Playlist, audio *entity.Audio) error
	FindAudioByIDFn func(id uint) (*entity.Audio, error)
	RemoveAudioFn   func(playlist *entity.Playlist, audio *entity.Audio) error
	FindByShareFn   func(token string) (*entity.Playlist, error)
	UpdateFn        func(id uint, updates map[string]interface{}) error
	DeleteFn        func(id uint) error
	CountAudiosFn   func(playlistID uint) (int, error)
}

func (m *MockPlaylistRepository) Create(playlist *entity.Playlist) error {
	return m.CreateFn(playlist)
}
func (m *MockPlaylistRepository) FindByUserID(userID uint) ([]entity.Playlist, error) {
	return m.FindByUserIDFn(userID)
}
func (m *MockPlaylistRepository) FindByID(id uint) (*entity.Playlist, error) {
	return m.FindByIDFn(id)
}
func (m *MockPlaylistRepository) FindAll() ([]entity.Playlist, error) {
	return m.FindAllFn()
}
func (m *MockPlaylistRepository) Search(query string) ([]entity.Playlist, error) {
	return m.SearchFn(query)
}
func (m *MockPlaylistRepository) AddAudio(playlist *entity.Playlist, audio *entity.Audio) error {
	return m.AddAudioFn(playlist, audio)
}
func (m *MockPlaylistRepository) FindAudioByID(id uint) (*entity.Audio, error) {
	return m.FindAudioByIDFn(id)
}
func (m *MockPlaylistRepository) RemoveAudio(playlist *entity.Playlist, audio *entity.Audio) error {
	return m.RemoveAudioFn(playlist, audio)
}
func (m *MockPlaylistRepository) FindByShareToken(token string) (*entity.Playlist, error) {
	return m.FindByShareFn(token)
}
func (m *MockPlaylistRepository) Update(id uint, updates map[string]interface{}) error {
	return m.UpdateFn(id, updates)
}
func (m *MockPlaylistRepository) Delete(id uint) error {
	return m.DeleteFn(id)
}
func (m *MockPlaylistRepository) CountAudios(playlistID uint) (int, error) {
	return m.CountAudiosFn(playlistID)
}

var _ port.PlaylistRepository = (*MockPlaylistRepository)(nil)
