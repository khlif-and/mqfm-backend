package service_test

import (
"errors"
"testing"

"github.com/stretchr/testify/assert"

"mqfm-backend/internal/app/service"
"mqfm-backend/internal/domain/entity"
	audiomock "mqfm-backend/tests/mocks/audio"
	playlistmock "mqfm-backend/tests/mocks/playlist"
)

func TestPlaylistCreate_Success(t *testing.T) {
repo := &playlistmock.MockPlaylistRepository{
CreateFn: func(playlist *entity.Playlist) error { playlist.ID = 1; return nil },
}
colorSvc := &audiomock.MockColorExtractorService{ExtractFn: func(p string) (string, error) { return "#fff", nil }}
svc := service.NewPlaylistService(repo, colorSvc)
playlist := &entity.Playlist{UserID: 1, Name: "My Playlist"}
err := svc.Create(playlist, nil)
assert.NoError(t, err)
assert.Equal(t, uint(1), playlist.ID)
}

func TestPlaylistGetByUserID_Success(t *testing.T) {
repo := &playlistmock.MockPlaylistRepository{
FindByUserIDFn: func(userID uint) ([]entity.Playlist, error) {
return []entity.Playlist{{ID: 1, UserID: userID, Name: "P1"}, {ID: 2, UserID: userID, Name: "P2"}}, nil
},
}
colorSvc := &audiomock.MockColorExtractorService{ExtractFn: func(p string) (string, error) { return "", nil }}
svc := service.NewPlaylistService(repo, colorSvc)
playlists, err := svc.GetByUserID(1)
assert.NoError(t, err)
assert.Len(t, playlists, 2)
}

func TestPlaylistAddAudio_Success(t *testing.T) {
audio := &entity.Audio{ID: 10, Title: "Audio 1", Thumbnail: "uploads/thumbnails/a.jpg"}
var updatedFields map[string]interface{}
repo := &playlistmock.MockPlaylistRepository{
FindByIDFn:      func(id uint) (*entity.Playlist, error) { return &entity.Playlist{ID: id, UserID: 1, Audios: []*entity.Audio{}}, nil },
FindAudioByIDFn: func(id uint) (*entity.Audio, error) { return audio, nil },
AddAudioFn:      func(p *entity.Playlist, a *entity.Audio) error { return nil },
CountAudiosFn:   func(pid uint) (int, error) { return 0, nil },
UpdateFn:        func(id uint, u map[string]interface{}) error { updatedFields = u; return nil },
}
colorSvc := &audiomock.MockColorExtractorService{ExtractFn: func(p string) (string, error) { return "#abc", nil }}
svc := service.NewPlaylistService(repo, colorSvc)
err := svc.AddAudioToPlaylist(1, 10)
assert.NoError(t, err)
assert.Equal(t, "uploads/thumbnails/a.jpg", updatedFields["image_url"])
assert.Equal(t, "#abc", updatedFields["dominant_color"])
}

func TestPlaylistAddAudio_NoOverwriteExistingImage(t *testing.T) {
audio := &entity.Audio{ID: 10, Title: "Audio 1", Thumbnail: "uploads/thumbnails/a.jpg"}
var updatedFields map[string]interface{}
repo := &playlistmock.MockPlaylistRepository{
FindByIDFn: func(id uint) (*entity.Playlist, error) {
return &entity.Playlist{ID: id, UserID: 1, ImageURL: "uploads/playlists/custom.jpg", Audios: []*entity.Audio{}}, nil
},
FindAudioByIDFn: func(id uint) (*entity.Audio, error) { return audio, nil },
AddAudioFn:      func(p *entity.Playlist, a *entity.Audio) error { return nil },
CountAudiosFn:   func(pid uint) (int, error) { return 0, nil },
UpdateFn:        func(id uint, u map[string]interface{}) error { updatedFields = u; return nil },
}
colorSvc := &audiomock.MockColorExtractorService{ExtractFn: func(p string) (string, error) { return "#abc", nil }}
svc := service.NewPlaylistService(repo, colorSvc)
err := svc.AddAudioToPlaylist(1, 10)
assert.NoError(t, err)
assert.Nil(t, updatedFields["image_url"])
}

func TestPlaylistAddAudio_PlaylistFull(t *testing.T) {
repo := &playlistmock.MockPlaylistRepository{
FindByIDFn:    func(id uint) (*entity.Playlist, error) { return &entity.Playlist{ID: id, UserID: 1}, nil },
CountAudiosFn: func(pid uint) (int, error) { return 20, nil },
}
colorSvc := &audiomock.MockColorExtractorService{ExtractFn: func(p string) (string, error) { return "", nil }}
svc := service.NewPlaylistService(repo, colorSvc)
err := svc.AddAudioToPlaylist(1, 21)
assert.Error(t, err)
assert.Contains(t, err.Error(), "full")
}

func TestPlaylistAddAudio_AlreadyExists(t *testing.T) {
repo := &playlistmock.MockPlaylistRepository{
FindByIDFn:    func(id uint) (*entity.Playlist, error) {
return &entity.Playlist{ID: id, UserID: 1, Audios: []*entity.Audio{{ID: 10}}}, nil
},
CountAudiosFn: func(pid uint) (int, error) { return 1, nil },
}
colorSvc := &audiomock.MockColorExtractorService{ExtractFn: func(p string) (string, error) { return "", nil }}
svc := service.NewPlaylistService(repo, colorSvc)
err := svc.AddAudioToPlaylist(1, 10)
assert.Error(t, err)
assert.Contains(t, err.Error(), "already")
}

func TestPlaylistAddAudio_NotFound(t *testing.T) {
repo := &playlistmock.MockPlaylistRepository{
FindByIDFn: func(id uint) (*entity.Playlist, error) { return nil, errors.New("not found") },
}
colorSvc := &audiomock.MockColorExtractorService{ExtractFn: func(p string) (string, error) { return "", nil }}
svc := service.NewPlaylistService(repo, colorSvc)
err := svc.AddAudioToPlaylist(999, 10)
assert.Error(t, err)
}

func TestPlaylistRemoveAudio_Success(t *testing.T) {
audio := &entity.Audio{ID: 10}
repo := &playlistmock.MockPlaylistRepository{
FindByIDFn:      func(id uint) (*entity.Playlist, error) { return &entity.Playlist{ID: id, UserID: 1}, nil },
FindAudioByIDFn: func(id uint) (*entity.Audio, error) { return audio, nil },
RemoveAudioFn:   func(p *entity.Playlist, a *entity.Audio) error { return nil },
}
colorSvc := &audiomock.MockColorExtractorService{ExtractFn: func(p string) (string, error) { return "", nil }}
svc := service.NewPlaylistService(repo, colorSvc)
err := svc.RemoveAudioFromPlaylist(1, 10)
assert.NoError(t, err)
}

func TestPlaylistSharePlaylist_NewToken(t *testing.T) {
repo := &playlistmock.MockPlaylistRepository{
FindByIDFn: func(id uint) (*entity.Playlist, error) {
return &entity.Playlist{ID: id, UserID: 1, ShareToken: ""}, nil
},
UpdateFn: func(id uint, u map[string]interface{}) error { return nil },
}
colorSvc := &audiomock.MockColorExtractorService{ExtractFn: func(p string) (string, error) { return "", nil }}
svc := service.NewPlaylistService(repo, colorSvc)
token, err := svc.SharePlaylist(1)
assert.NoError(t, err)
assert.NotEmpty(t, token)
}

func TestPlaylistSharePlaylist_ExistingToken(t *testing.T) {
repo := &playlistmock.MockPlaylistRepository{
FindByIDFn: func(id uint) (*entity.Playlist, error) {
return &entity.Playlist{ID: id, UserID: 1, ShareToken: "existing-token"}, nil
},
}
colorSvc := &audiomock.MockColorExtractorService{ExtractFn: func(p string) (string, error) { return "", nil }}
svc := service.NewPlaylistService(repo, colorSvc)
token, err := svc.SharePlaylist(1)
assert.NoError(t, err)
assert.Equal(t, "existing-token", token)
}

func TestPlaylistGetByShareToken_Success(t *testing.T) {
repo := &playlistmock.MockPlaylistRepository{
FindByShareFn: func(token string) (*entity.Playlist, error) {
return &entity.Playlist{ID: 1, Name: "Shared", ShareToken: token}, nil
},
}
colorSvc := &audiomock.MockColorExtractorService{ExtractFn: func(p string) (string, error) { return "", nil }}
svc := service.NewPlaylistService(repo, colorSvc)
playlist, err := svc.GetByShareToken("test-token")
assert.NoError(t, err)
assert.Equal(t, "Shared", playlist.Name)
}
