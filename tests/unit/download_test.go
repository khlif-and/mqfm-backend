package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"mqfm-backend/internal/app/service"
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/shared/dto/request"
	audiomock "mqfm-backend/tests/mocks/audio"
	downloadmock "mqfm-backend/tests/mocks/download"
	favoritemock "mqfm-backend/tests/mocks/favorite"
)

// ──────────────────────── RecordDownload ────────────────────────

func TestDownload_RecordDownload_Success(t *testing.T) {
	var created *entity.Download
	dlRepo := &downloadmock.MockDownloadRepository{
		ExistsFn: func(userID, audioID uint) (bool, error) {
			assert.Equal(t, uint(1), userID)
			assert.Equal(t, uint(10), audioID)
			return false, nil
		},
		CreateFn: func(download *entity.Download) error {
			created = download
			download.ID = 1
			return nil
		},
	}
	audioRepo := &audiomock.MockAudioRepository{
		FindByIDFn: func(id uint) (*entity.Audio, error) {
			return &entity.Audio{ID: 10, Title: "Tafsir Al-Quran", Artist: "Ust. Adi Hidayat", FileSize: 5000000}, nil
		},
	}
	favRepo := &favoritemock.MockFavoriteArtistRepository{}
	svc := service.NewDownloadService(dlRepo, audioRepo, favRepo)

	result, err := svc.RecordDownload(1, request.DownloadRequest{AudioID: 10, PlaylistID: 5})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, created)
	assert.Equal(t, uint(1), result.UserID)
	assert.Equal(t, uint(10), result.AudioID)
	assert.Equal(t, uint(5), result.PlaylistID)
	assert.Equal(t, int64(5000000), result.FileSize, "should use audio.FileSize")
	assert.NotNil(t, result.Audio, "audio must be attached")
	assert.Equal(t, "Tafsir Al-Quran", result.Audio.Title)
	assert.True(t, result.ExpiresAt.After(time.Now().Add(29*24*time.Hour)), "expiry must be ~30 days")
}

func TestDownload_RecordDownload_FileSizeOverride(t *testing.T) {
	dlRepo := &downloadmock.MockDownloadRepository{
		ExistsFn: func(userID, audioID uint) (bool, error) { return false, nil },
		CreateFn: func(download *entity.Download) error { return nil },
	}
	audioRepo := &audiomock.MockAudioRepository{
		FindByIDFn: func(id uint) (*entity.Audio, error) {
			return &entity.Audio{ID: 10, FileSize: 5000000}, nil
		},
	}
	favRepo := &favoritemock.MockFavoriteArtistRepository{}
	svc := service.NewDownloadService(dlRepo, audioRepo, favRepo)

	result, err := svc.RecordDownload(1, request.DownloadRequest{AudioID: 10, FileSize: 9999999})
	assert.NoError(t, err)
	assert.Equal(t, int64(9999999), result.FileSize, "request FileSize > 0 should override audio FileSize")
}

func TestDownload_RecordDownload_Duplicate(t *testing.T) {
	dlRepo := &downloadmock.MockDownloadRepository{
		ExistsFn: func(userID, audioID uint) (bool, error) { return true, nil },
	}
	audioRepo := &audiomock.MockAudioRepository{}
	favRepo := &favoritemock.MockFavoriteArtistRepository{}
	svc := service.NewDownloadService(dlRepo, audioRepo, favRepo)

	result, err := svc.RecordDownload(1, request.DownloadRequest{AudioID: 10})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "audio already downloaded", err.Error())
}

func TestDownload_RecordDownload_AudioNotFound(t *testing.T) {
	dlRepo := &downloadmock.MockDownloadRepository{
		ExistsFn: func(userID, audioID uint) (bool, error) { return false, nil },
	}
	audioRepo := &audiomock.MockAudioRepository{
		FindByIDFn: func(id uint) (*entity.Audio, error) {
			return nil, errors.New("not found")
		},
	}
	favRepo := &favoritemock.MockFavoriteArtistRepository{}
	svc := service.NewDownloadService(dlRepo, audioRepo, favRepo)

	result, err := svc.RecordDownload(1, request.DownloadRequest{AudioID: 999})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "audio not found", err.Error())
}

func TestDownload_RecordDownload_CreateError(t *testing.T) {
	dlRepo := &downloadmock.MockDownloadRepository{
		ExistsFn: func(userID, audioID uint) (bool, error) { return false, nil },
		CreateFn: func(download *entity.Download) error { return errors.New("create failed") },
	}
	audioRepo := &audiomock.MockAudioRepository{
		FindByIDFn: func(id uint) (*entity.Audio, error) {
			return &entity.Audio{ID: 10, FileSize: 1000}, nil
		},
	}
	favRepo := &favoritemock.MockFavoriteArtistRepository{}
	svc := service.NewDownloadService(dlRepo, audioRepo, favRepo)

	result, err := svc.RecordDownload(1, request.DownloadRequest{AudioID: 10})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "create failed", err.Error())
}

// ──────────────────────── GetDownloads ────────────────────────

func TestDownload_GetDownloads_Success(t *testing.T) {
	expected := []entity.Download{
		{ID: 1, UserID: 1, AudioID: 10},
		{ID: 2, UserID: 1, AudioID: 20},
	}
	dlRepo := &downloadmock.MockDownloadRepository{
		FindByUserFn: func(userID uint) ([]entity.Download, error) {
			assert.Equal(t, uint(1), userID)
			return expected, nil
		},
	}
	audioRepo := &audiomock.MockAudioRepository{}
	favRepo := &favoritemock.MockFavoriteArtistRepository{}
	svc := service.NewDownloadService(dlRepo, audioRepo, favRepo)

	result, err := svc.GetDownloads(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestDownload_GetDownloads_Error(t *testing.T) {
	dlRepo := &downloadmock.MockDownloadRepository{
		FindByUserFn: func(userID uint) ([]entity.Download, error) {
			return nil, errors.New("db error")
		},
	}
	audioRepo := &audiomock.MockAudioRepository{}
	favRepo := &favoritemock.MockFavoriteArtistRepository{}
	svc := service.NewDownloadService(dlRepo, audioRepo, favRepo)

	result, err := svc.GetDownloads(1)
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ──────────────────────── DeleteDownload ────────────────────────

func TestDownload_DeleteDownload_Success(t *testing.T) {
	dlRepo := &downloadmock.MockDownloadRepository{
		DeleteFn: func(id, userID uint) error {
			assert.Equal(t, uint(5), id)
			assert.Equal(t, uint(1), userID)
			return nil
		},
	}
	audioRepo := &audiomock.MockAudioRepository{}
	favRepo := &favoritemock.MockFavoriteArtistRepository{}
	svc := service.NewDownloadService(dlRepo, audioRepo, favRepo)

	err := svc.DeleteDownload(5, 1)
	assert.NoError(t, err)
}

func TestDownload_DeleteDownload_Error(t *testing.T) {
	dlRepo := &downloadmock.MockDownloadRepository{
		DeleteFn: func(id, userID uint) error { return errors.New("delete failed") },
	}
	audioRepo := &audiomock.MockAudioRepository{}
	favRepo := &favoritemock.MockFavoriteArtistRepository{}
	svc := service.NewDownloadService(dlRepo, audioRepo, favRepo)

	err := svc.DeleteDownload(5, 1)
	assert.Error(t, err)
}

// ──────────────────────── GetStorageUsage ────────────────────────

func TestDownload_GetStorageUsage_Success(t *testing.T) {
	dlRepo := &downloadmock.MockDownloadRepository{
		SumSizeByUserFn: func(userID uint) (int64, error) {
			assert.Equal(t, uint(1), userID)
			return 150000000, nil
		},
	}
	audioRepo := &audiomock.MockAudioRepository{}
	favRepo := &favoritemock.MockFavoriteArtistRepository{}
	svc := service.NewDownloadService(dlRepo, audioRepo, favRepo)

	size, err := svc.GetStorageUsage(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(150000000), size)
}

// ──────────────────────── GetNewFromFavorites ────────────────────────

func TestDownload_GetNewFromFavorites_WithFavorites(t *testing.T) {
	favRepo := &favoritemock.MockFavoriteArtistRepository{
		FindByUserFn: func(userID uint) ([]entity.FavoriteArtist, error) {
			return []entity.FavoriteArtist{
				{ArtistName: "Ust. Adi Hidayat"},
				{ArtistName: "Ust. Abdul Somad"},
			}, nil
		},
	}
	audioRepo := &audiomock.MockAudioRepository{
		FindNewByArtistsFn: func(artists []string, since time.Time) ([]entity.Audio, error) {
			assert.Len(t, artists, 2)
			assert.Contains(t, artists, "Ust. Adi Hidayat")
			assert.Contains(t, artists, "Ust. Abdul Somad")
			assert.True(t, since.After(time.Now().AddDate(0, 0, -8)), "since should be within ~7 days ago")
			return []entity.Audio{{ID: 100, Title: "Baru"}}, nil
		},
	}
	dlRepo := &downloadmock.MockDownloadRepository{}
	svc := service.NewDownloadService(dlRepo, audioRepo, favRepo)

	result, err := svc.GetNewFromFavorites(1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Baru", result[0].Title)
}

func TestDownload_GetNewFromFavorites_NoFavorites(t *testing.T) {
	favRepo := &favoritemock.MockFavoriteArtistRepository{
		FindByUserFn: func(userID uint) ([]entity.FavoriteArtist, error) {
			return []entity.FavoriteArtist{}, nil
		},
	}
	audioRepo := &audiomock.MockAudioRepository{}
	dlRepo := &downloadmock.MockDownloadRepository{}
	svc := service.NewDownloadService(dlRepo, audioRepo, favRepo)

	result, err := svc.GetNewFromFavorites(1)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestDownload_GetNewFromFavorites_FavError(t *testing.T) {
	favRepo := &favoritemock.MockFavoriteArtistRepository{
		FindByUserFn: func(userID uint) ([]entity.FavoriteArtist, error) {
			return nil, errors.New("db error")
		},
	}
	audioRepo := &audiomock.MockAudioRepository{}
	dlRepo := &downloadmock.MockDownloadRepository{}
	svc := service.NewDownloadService(dlRepo, audioRepo, favRepo)

	result, err := svc.GetNewFromFavorites(1)
	assert.NoError(t, err, "service returns nil,nil on fav error")
	assert.Nil(t, result)
}

// ──────────────────────── CleanupExpired ────────────────────────

func TestDownload_CleanupExpired_Success(t *testing.T) {
	dlRepo := &downloadmock.MockDownloadRepository{
		DeleteExpiredFn: func() (int64, error) { return 5, nil },
	}
	audioRepo := &audiomock.MockAudioRepository{}
	favRepo := &favoritemock.MockFavoriteArtistRepository{}
	svc := service.NewDownloadService(dlRepo, audioRepo, favRepo)

	count, err := svc.CleanupExpired()
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
}

func TestDownload_CleanupExpired_Error(t *testing.T) {
	dlRepo := &downloadmock.MockDownloadRepository{
		DeleteExpiredFn: func() (int64, error) { return 0, errors.New("cleanup failed") },
	}
	audioRepo := &audiomock.MockAudioRepository{}
	favRepo := &favoritemock.MockFavoriteArtistRepository{}
	svc := service.NewDownloadService(dlRepo, audioRepo, favRepo)

	count, err := svc.CleanupExpired()
	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
}
