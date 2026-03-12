package resumemock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type MockSmartResumeRepository struct {
	UpsertFn     func(resume *entity.SmartResume) error
	FindByUserFn func(userID uint) (*entity.SmartResume, error)
}

func (m *MockSmartResumeRepository) Upsert(resume *entity.SmartResume) error {
	return m.UpsertFn(resume)
}
func (m *MockSmartResumeRepository) FindByUser(userID uint) (*entity.SmartResume, error) {
	return m.FindByUserFn(userID)
}

var _ port.SmartResumeRepository = (*MockSmartResumeRepository)(nil)
