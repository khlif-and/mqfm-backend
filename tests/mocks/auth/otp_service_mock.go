package authmock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type MockOTPService struct {
	SendOTPFn   func(email string) error
	VerifyOTPFn func(email string, code string) (*entity.User, error)
}

func (m *MockOTPService) SendOTP(email string) error              { return m.SendOTPFn(email) }
func (m *MockOTPService) VerifyOTP(email, code string) (*entity.User, error) {
	return m.VerifyOTPFn(email, code)
}

var _ port.OTPService = (*MockOTPService)(nil)
