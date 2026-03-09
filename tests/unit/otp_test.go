package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"mqfm-backend/internal/app/service"
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/tests/mocks"
)

func TestOTPSendOTP_Success(t *testing.T) {
	var sentEmail string
	otpRepo := &mocks.MockOTPRepository{
		CountRecentFn: func(email string, since time.Time) (int64, error) { return 0, nil },
		CreateFn:      func(otp *entity.OTP) error { return nil },
	}
	userRepo := &mocks.MockUserRepository{}
	emailSvc := &mocks.MockEmailService{
		SendAsyncFn: func(to, subject, body string) { sentEmail = to },
	}

	svc := service.NewOTPService(otpRepo, userRepo, emailSvc)
	err := svc.SendOTP("test@test.com")

	assert.NoError(t, err)
	assert.Equal(t, "test@test.com", sentEmail)
}

func TestOTPSendOTP_RateLimited(t *testing.T) {
	otpRepo := &mocks.MockOTPRepository{
		CountRecentFn: func(email string, since time.Time) (int64, error) { return 3, nil },
	}
	userRepo := &mocks.MockUserRepository{}
	emailSvc := &mocks.MockEmailService{}

	svc := service.NewOTPService(otpRepo, userRepo, emailSvc)
	err := svc.SendOTP("test@test.com")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Too many")
}

func TestOTPVerifyOTP_Success(t *testing.T) {
	otpRepo := &mocks.MockOTPRepository{
		FindLatestFn: func(email string) (*entity.OTP, error) {
			return &entity.OTP{
				ID:        1,
				Email:     email,
				Code:      "123456",
				ExpiresAt: time.Now().Add(5 * time.Minute),
			}, nil
		},
		MarkVerifiedFn: func(id uint) error { return nil },
	}
	userRepo := &mocks.MockUserRepository{
		FindByEmailFn: func(email string) (*entity.User, error) {
			return &entity.User{ID: 1, Email: email}, nil
		},
		UpdateFn: func(id uint, updates map[string]interface{}) error { return nil },
	}
	emailSvc := &mocks.MockEmailService{}

	svc := service.NewOTPService(otpRepo, userRepo, emailSvc)
	user, err := svc.VerifyOTP("test@test.com", "123456")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.True(t, user.EmailVerified)
}

func TestOTPVerifyOTP_Invalid(t *testing.T) {
	otpRepo := &mocks.MockOTPRepository{
		FindLatestFn: func(email string) (*entity.OTP, error) {
			return &entity.OTP{
				Code:      "123456",
				ExpiresAt: time.Now().Add(5 * time.Minute),
			}, nil
		},
	}
	userRepo := &mocks.MockUserRepository{}
	emailSvc := &mocks.MockEmailService{}

	svc := service.NewOTPService(otpRepo, userRepo, emailSvc)
	user, err := svc.VerifyOTP("test@test.com", "000000")

	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestOTPVerifyOTP_Expired(t *testing.T) {
	otpRepo := &mocks.MockOTPRepository{
		FindLatestFn: func(email string) (*entity.OTP, error) {
			return &entity.OTP{
				Code:      "123456",
				ExpiresAt: time.Now().Add(-1 * time.Minute),
			}, nil
		},
	}
	userRepo := &mocks.MockUserRepository{}
	emailSvc := &mocks.MockEmailService{}

	svc := service.NewOTPService(otpRepo, userRepo, emailSvc)
	user, err := svc.VerifyOTP("test@test.com", "123456")

	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestOTPVerifyOTP_NotFound(t *testing.T) {
	otpRepo := &mocks.MockOTPRepository{
		FindLatestFn: func(email string) (*entity.OTP, error) {
			return nil, errors.New("not found")
		},
	}
	userRepo := &mocks.MockUserRepository{}
	emailSvc := &mocks.MockEmailService{}

	svc := service.NewOTPService(otpRepo, userRepo, emailSvc)
	user, err := svc.VerifyOTP("test@test.com", "123456")

	assert.Error(t, err)
	assert.Nil(t, user)
}
