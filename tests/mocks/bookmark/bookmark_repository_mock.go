package bookmarkmock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type MockBookmarkRepository struct {
	CreateFn               func(bookmark *entity.Bookmark) error
	FindByUserFn           func(userID uint) ([]entity.Bookmark, error)
	FindByUserAndAudioFn   func(userID, audioID uint) ([]entity.Bookmark, error)
	DeleteFn               func(id, userID uint) error
	DeleteAllByUserAudioFn func(userID, audioID uint) error
}

func (m *MockBookmarkRepository) Create(bookmark *entity.Bookmark) error {
	return m.CreateFn(bookmark)
}
func (m *MockBookmarkRepository) FindByUser(userID uint) ([]entity.Bookmark, error) {
	return m.FindByUserFn(userID)
}
func (m *MockBookmarkRepository) FindByUserAndAudio(userID, audioID uint) ([]entity.Bookmark, error) {
	return m.FindByUserAndAudioFn(userID, audioID)
}
func (m *MockBookmarkRepository) Delete(id, userID uint) error { return m.DeleteFn(id, userID) }
func (m *MockBookmarkRepository) DeleteAllByUserAndAudio(userID, audioID uint) error {
	return m.DeleteAllByUserAudioFn(userID, audioID)
}

var _ port.BookmarkRepository = (*MockBookmarkRepository)(nil)
