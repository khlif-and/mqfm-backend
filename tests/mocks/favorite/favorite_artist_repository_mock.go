package favoritemock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type MockFavoriteArtistRepository struct {
	CreateFn     func(fav *entity.FavoriteArtist) error
	DeleteFn     func(userID uint, artistName string) error
	FindByUserFn func(userID uint) ([]entity.FavoriteArtist, error)
	ExistsFn     func(userID uint, artistName string) (bool, error)
}

func (m *MockFavoriteArtistRepository) Create(fav *entity.FavoriteArtist) error {
	return m.CreateFn(fav)
}
func (m *MockFavoriteArtistRepository) Delete(userID uint, artistName string) error {
	return m.DeleteFn(userID, artistName)
}
func (m *MockFavoriteArtistRepository) FindByUser(userID uint) ([]entity.FavoriteArtist, error) {
	return m.FindByUserFn(userID)
}
func (m *MockFavoriteArtistRepository) Exists(userID uint, artistName string) (bool, error) {
	return m.ExistsFn(userID, artistName)
}

var _ port.FavoriteArtistRepository = (*MockFavoriteArtistRepository)(nil)
