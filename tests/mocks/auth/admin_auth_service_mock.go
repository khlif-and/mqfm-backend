package authmock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
)

type MockAdminAuthService struct {
	RegisterFn func(req request.AdminRegisterRequest) (*entity.Admin, error)
	LoginFn    func(req request.AdminLoginRequest) (string, *entity.Admin, error)
	UpdateFn   func(id uint, updates map[string]interface{}) (*entity.Admin, error)
	GetByIDFn  func(id uint) (*entity.Admin, error)
}

func (m *MockAdminAuthService) Register(req request.AdminRegisterRequest) (*entity.Admin, error) {
	return m.RegisterFn(req)
}
func (m *MockAdminAuthService) Login(req request.AdminLoginRequest) (string, *entity.Admin, error) {
	return m.LoginFn(req)
}
func (m *MockAdminAuthService) UpdateAdmin(id uint, updates map[string]interface{}) (*entity.Admin, error) {
	return m.UpdateFn(id, updates)
}
func (m *MockAdminAuthService) GetAdminByID(id uint) (*entity.Admin, error) {
	return m.GetByIDFn(id)
}

var _ port.AdminAuthService = (*MockAdminAuthService)(nil)
