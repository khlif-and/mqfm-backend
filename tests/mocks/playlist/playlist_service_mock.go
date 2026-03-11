package playlistmock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type MockPlaylistService struct {
	CreateFn          func(playlist *entity.Playlist, file interface{}) error
	GetByUserIDFn     func(userID uint) ([]entity.Playlist, error)
	GetByIDFn         func(id uint) (*entity.Playlist, error)
	UpdateFn          func(id uint, updates map[string]interface{}) (*entity.Playlist, error)
	DeleteFn          func(id, userID uint) error
	SearchFn          func(query string) ([]entity.Playlist, error)
	AddAudioFn        func(playlistID, audioID uint) error
	RemoveAudioFn     func(playlistID, audioID uint) error
	SharePlaylistFn   func(playlistID uint) (string, error)
	GetByShareTokenFn func(token string) (*entity.Playlist, error)
}

func (m *MockPlaylistService) Create(playlist *entity.Playlist, file interface{}) error {
	return m.CreateFn(playlist, file)
}
func (m *MockPlaylistService) GetByUserID(userID uint) ([]entity.Playlist, error) {
	return m.GetByUserIDFn(userID)
}
func (m *MockPlaylistService) GetByID(id uint) (*entity.Playlist, error) {
	return m.GetByIDFn(id)
}
func (m *MockPlaylistService) Update(id uint, updates map[string]interface{}) (*entity.Playlist, error) {
	return m.UpdateFn(id, updates)
}
func (m *MockPlaylistService) Delete(id, userID uint) error { return m.DeleteFn(id, userID) }
func (m *MockPlaylistService) Search(query string) ([]entity.Playlist, error) {
	return m.SearchFn(query)
}
func (m *MockPlaylistService) AddAudioToPlaylist(playlistID, audioID uint) error {
	return m.AddAudioFn(playlistID, audioID)
}
func (m *MockPlaylistService) RemoveAudioFromPlaylist(playlistID, audioID uint) error {
	return m.RemoveAudioFn(playlistID, audioID)
}
func (m *MockPlaylistService) SharePlaylist(playlistID uint) (string, error) {
	return m.SharePlaylistFn(playlistID)
}
func (m *MockPlaylistService) GetByShareToken(token string) (*entity.Playlist, error) {
	return m.GetByShareTokenFn(token)
}

var _ port.PlaylistService = (*MockPlaylistService)(nil)
