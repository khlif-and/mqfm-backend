package locationmock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type MockUserLocationRepository struct {
	UpsertFn     func(location *entity.UserLocation) error
	FindByUserFn func(userID uint) (*entity.UserLocation, error)
	FindAllFn    func() ([]entity.UserLocation, error)
}

func (m *MockUserLocationRepository) Upsert(location *entity.UserLocation) error {
	return m.UpsertFn(location)
}
func (m *MockUserLocationRepository) FindByUser(userID uint) (*entity.UserLocation, error) {
	return m.FindByUserFn(userID)
}
func (m *MockUserLocationRepository) FindAll() ([]entity.UserLocation, error) {
	return m.FindAllFn()
}

var _ port.UserLocationRepository = (*MockUserLocationRepository)(nil)
