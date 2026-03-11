package downloadmock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type MockDownloadRepository struct {
	CreateFn        func(download *entity.Download) error
	FindByUserFn    func(userID uint) ([]entity.Download, error)
	DeleteFn        func(id, userID uint) error
	ExistsFn        func(userID, audioID uint) (bool, error)
	SumSizeByUserFn func(userID uint) (int64, error)
	DeleteExpiredFn func() (int64, error)
}

func (m *MockDownloadRepository) Create(download *entity.Download) error {
	return m.CreateFn(download)
}
func (m *MockDownloadRepository) FindByUser(userID uint) ([]entity.Download, error) {
	return m.FindByUserFn(userID)
}
func (m *MockDownloadRepository) Delete(id, userID uint) error { return m.DeleteFn(id, userID) }
func (m *MockDownloadRepository) Exists(userID, audioID uint) (bool, error) {
	return m.ExistsFn(userID, audioID)
}
func (m *MockDownloadRepository) SumSizeByUser(userID uint) (int64, error) {
	return m.SumSizeByUserFn(userID)
}
func (m *MockDownloadRepository) DeleteExpired() (int64, error) { return m.DeleteExpiredFn() }

var _ port.DownloadRepository = (*MockDownloadRepository)(nil)
