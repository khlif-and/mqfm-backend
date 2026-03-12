package preferencemock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type MockUserPreferenceRepository struct {
	UpsertFn     func(pref *entity.UserPreference) error
	FindByUserFn func(userID uint) (*entity.UserPreference, error)
}

func (m *MockUserPreferenceRepository) Upsert(pref *entity.UserPreference) error {
	return m.UpsertFn(pref)
}
func (m *MockUserPreferenceRepository) FindByUser(userID uint) (*entity.UserPreference, error) {
	return m.FindByUserFn(userID)
}

var _ port.UserPreferenceRepository = (*MockUserPreferenceRepository)(nil)
