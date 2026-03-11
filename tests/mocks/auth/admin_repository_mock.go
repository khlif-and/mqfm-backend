package authmock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type MockAdminRepository struct {
	CreateFn      func(admin *entity.Admin) error
	FindByEmailFn func(email string) (*entity.Admin, error)
	FindByIDFn    func(id uint) (*entity.Admin, error)
	UpdateFn      func(id uint, updates map[string]interface{}) error
}

func (m *MockAdminRepository) Create(admin *entity.Admin) error      { return m.CreateFn(admin) }
func (m *MockAdminRepository) FindByEmail(email string) (*entity.Admin, error) {
	return m.FindByEmailFn(email)
}
func (m *MockAdminRepository) FindByID(id uint) (*entity.Admin, error) { return m.FindByIDFn(id) }
func (m *MockAdminRepository) Update(id uint, updates map[string]interface{}) error {
	return m.UpdateFn(id, updates)
}

var _ port.AdminRepository = (*MockAdminRepository)(nil)
