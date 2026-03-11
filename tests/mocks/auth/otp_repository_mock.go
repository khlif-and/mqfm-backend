package authmock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"time"
)

type MockOTPRepository struct {
	CreateFn        func(otp *entity.OTP) error
	FindLatestFn    func(email string) (*entity.OTP, error)
	MarkVerifiedFn  func(id uint) error
	CountRecentFn   func(email string, since time.Time) (int64, error)
	DeleteExpiredFn func() error
}

func (m *MockOTPRepository) Create(otp *entity.OTP) error { return m.CreateFn(otp) }
func (m *MockOTPRepository) FindLatestByEmail(email string) (*entity.OTP, error) {
	return m.FindLatestFn(email)
}
func (m *MockOTPRepository) MarkVerified(id uint) error { return m.MarkVerifiedFn(id) }
func (m *MockOTPRepository) CountRecentByEmail(email string, since time.Time) (int64, error) {
	return m.CountRecentFn(email, since)
}
func (m *MockOTPRepository) DeleteExpired() error { return m.DeleteExpiredFn() }

var _ port.OTPRepository = (*MockOTPRepository)(nil)
