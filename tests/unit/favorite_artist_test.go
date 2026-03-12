package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"mqfm-backend/internal/app/service"
	"mqfm-backend/internal/domain/entity"
	favoritemock "mqfm-backend/tests/mocks/favorite"
)

// ──────────────────────── Add ────────────────────────

func TestFavoriteArtist_Add_Success(t *testing.T) {
	var created *entity.FavoriteArtist
	repo := &favoritemock.MockFavoriteArtistRepository{
		ExistsFn: func(userID uint, artistName string) (bool, error) {
			assert.Equal(t, uint(1), userID)
			assert.Equal(t, "Ust. Adi Hidayat", artistName)
			return false, nil
		},
		CreateFn: func(fav *entity.FavoriteArtist) error {
			created = fav
			assert.Equal(t, uint(1), fav.UserID)
			assert.Equal(t, "Ust. Adi Hidayat", fav.ArtistName)
			return nil
		},
	}
	svc := service.NewFavoriteArtistService(repo)

	err := svc.Add(1, "Ust. Adi Hidayat")
	assert.NoError(t, err)
	assert.NotNil(t, created)
}

func TestFavoriteArtist_Add_Duplicate(t *testing.T) {
	repo := &favoritemock.MockFavoriteArtistRepository{
		ExistsFn: func(userID uint, artistName string) (bool, error) {
			return true, nil
		},
	}
	svc := service.NewFavoriteArtistService(repo)

	err := svc.Add(1, "Ust. Adi Hidayat")
	assert.Error(t, err)
	assert.Equal(t, "artist already in favorites", err.Error())
}

func TestFavoriteArtist_Add_ExistsCheckError(t *testing.T) {
	repo := &favoritemock.MockFavoriteArtistRepository{
		ExistsFn: func(userID uint, artistName string) (bool, error) {
			return false, errors.New("db error")
		},
	}
	svc := service.NewFavoriteArtistService(repo)

	err := svc.Add(1, "Ust. Adi Hidayat")
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

func TestFavoriteArtist_Add_CreateError(t *testing.T) {
	repo := &favoritemock.MockFavoriteArtistRepository{
		ExistsFn: func(userID uint, artistName string) (bool, error) {
			return false, nil
		},
		CreateFn: func(fav *entity.FavoriteArtist) error {
			return errors.New("create failed")
		},
	}
	svc := service.NewFavoriteArtistService(repo)

	err := svc.Add(1, "Ust. Adi Hidayat")
	assert.Error(t, err)
	assert.Equal(t, "create failed", err.Error())
}

// ──────────────────────── Remove ────────────────────────

func TestFavoriteArtist_Remove_Success(t *testing.T) {
	repo := &favoritemock.MockFavoriteArtistRepository{
		DeleteFn: func(userID uint, artistName string) error {
			assert.Equal(t, uint(1), userID)
			assert.Equal(t, "Ust. Abdul Somad", artistName)
			return nil
		},
	}
	svc := service.NewFavoriteArtistService(repo)

	err := svc.Remove(1, "Ust. Abdul Somad")
	assert.NoError(t, err)
}

func TestFavoriteArtist_Remove_Error(t *testing.T) {
	repo := &favoritemock.MockFavoriteArtistRepository{
		DeleteFn: func(userID uint, artistName string) error {
			return errors.New("delete failed")
		},
	}
	svc := service.NewFavoriteArtistService(repo)

	err := svc.Remove(1, "Ust. Abdul Somad")
	assert.Error(t, err)
}

// ──────────────────────── GetByUser ────────────────────────

func TestFavoriteArtist_GetByUser_Success(t *testing.T) {
	expected := []entity.FavoriteArtist{
		{ID: 1, UserID: 1, ArtistName: "Ust. Adi Hidayat"},
		{ID: 2, UserID: 1, ArtistName: "Ust. Abdul Somad"},
	}
	repo := &favoritemock.MockFavoriteArtistRepository{
		FindByUserFn: func(userID uint) ([]entity.FavoriteArtist, error) {
			assert.Equal(t, uint(1), userID)
			return expected, nil
		},
	}
	svc := service.NewFavoriteArtistService(repo)

	result, err := svc.GetByUser(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Ust. Adi Hidayat", result[0].ArtistName)
	assert.Equal(t, "Ust. Abdul Somad", result[1].ArtistName)
}

func TestFavoriteArtist_GetByUser_Empty(t *testing.T) {
	repo := &favoritemock.MockFavoriteArtistRepository{
		FindByUserFn: func(userID uint) ([]entity.FavoriteArtist, error) {
			return []entity.FavoriteArtist{}, nil
		},
	}
	svc := service.NewFavoriteArtistService(repo)

	result, err := svc.GetByUser(1)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestFavoriteArtist_GetByUser_Error(t *testing.T) {
	repo := &favoritemock.MockFavoriteArtistRepository{
		FindByUserFn: func(userID uint) ([]entity.FavoriteArtist, error) {
			return nil, errors.New("db error")
		},
	}
	svc := service.NewFavoriteArtistService(repo)

	result, err := svc.GetByUser(1)
	assert.Error(t, err)
	assert.Nil(t, result)
}
