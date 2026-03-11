package downloadmock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
)

type MockDownloadService struct {
	RecordDownloadFn func(userID uint, req request.DownloadRequest) (*entity.Download, error)
	GetDownloadsFn   func(userID uint) ([]entity.Download, error)
	DeleteDownloadFn func(id, userID uint) error
	GetStorageUsageFn func(userID uint) (int64, error)
	GetNewFromFavsFn func(userID uint) ([]entity.Audio, error)
	CleanupExpiredFn func() (int64, error)
}

func (m *MockDownloadService) RecordDownload(userID uint, req request.DownloadRequest) (*entity.Download, error) {
	return m.RecordDownloadFn(userID, req)
}
func (m *MockDownloadService) GetDownloads(userID uint) ([]entity.Download, error) {
	return m.GetDownloadsFn(userID)
}
func (m *MockDownloadService) DeleteDownload(id, userID uint) error {
	return m.DeleteDownloadFn(id, userID)
}
func (m *MockDownloadService) GetStorageUsage(userID uint) (int64, error) {
	return m.GetStorageUsageFn(userID)
}
func (m *MockDownloadService) GetNewFromFavorites(userID uint) ([]entity.Audio, error) {
	return m.GetNewFromFavsFn(userID)
}
func (m *MockDownloadService) CleanupExpired() (int64, error) {
	return m.CleanupExpiredFn()
}

var _ port.DownloadService = (*MockDownloadService)(nil)
