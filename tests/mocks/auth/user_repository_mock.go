package authmock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type MockUserRepository struct {
	CreateFn         func(user *entity.User) error
	FindByEmailFn    func(email string) (*entity.User, error)
	FindByIDFn       func(id uint) (*entity.User, error)
	FindByProviderFn func(provider string, providerID string) (*entity.User, error)
	UpdateFn         func(id uint, updates map[string]interface{}) error
}

func (m *MockUserRepository) Create(user *entity.User) error { return m.CreateFn(user) }
func (m *MockUserRepository) FindByEmail(email string) (*entity.User, error) {
	return m.FindByEmailFn(email)
}
func (m *MockUserRepository) FindByID(id uint) (*entity.User, error) { return m.FindByIDFn(id) }
func (m *MockUserRepository) FindByProviderID(provider string, providerID string) (*entity.User, error) {
	return m.FindByProviderFn(provider, providerID)
}
func (m *MockUserRepository) Update(id uint, updates map[string]interface{}) error {
	return m.UpdateFn(id, updates)
}

var _ port.UserRepository = (*MockUserRepository)(nil)
