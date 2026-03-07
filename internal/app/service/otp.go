package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/logger"
)

type otpService struct {
	otpRepo      port.OTPRepository
	emailService port.EmailService
}

func NewOTPService(otpRepo port.OTPRepository, emailSvc port.EmailService) port.OTPService {
	return &otpService{
		otpRepo:      otpRepo,
		emailService: emailSvc,
	}
}

func (s *otpService) SendOTP(email string) error {
	count, err := s.otpRepo.CountRecentByEmail(email, time.Now().Add(-1*time.Hour))
	if err != nil {
		return err
	}
	if count >= 3 {
		return errors.New(constant.MsgOTPRateLimit)
	}

	code := generateOTPCode()

	otp := entity.OTP{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := s.otpRepo.Create(&otp); err != nil {
		logger.Error("otp create failed: " + err.Error())
		return err
	}

	subject := "MQFM - Email Verification"
	body := fmt.Sprintf(
		`<div style="font-family:Arial,sans-serif;max-width:400px;margin:0 auto;padding:20px;">
			<h2 style="color:#1a1a2e;">MQFM Verification</h2>
			<p>Your verification code is:</p>
			<div style="background:#f0f0f0;padding:15px;text-align:center;font-size:32px;letter-spacing:8px;font-weight:bold;border-radius:8px;">%s</div>
			<p style="color:#666;font-size:12px;margin-top:15px;">This code expires in 5 minutes.</p>
		</div>`, code)

	s.emailService.SendAsync(email, subject, body)

	return nil
}

func (s *otpService) VerifyOTP(email string, code string) error {
	otp, err := s.otpRepo.FindLatestByEmail(email)
	if err != nil {
		return errors.New(constant.MsgOTPInvalid)
	}

	if time.Now().After(otp.ExpiresAt) {
		return errors.New(constant.MsgOTPExpired)
	}

	if otp.Code != code {
		return errors.New(constant.MsgOTPInvalid)
	}

	return s.otpRepo.MarkVerified(otp.ID)
}

func generateOTPCode() string {
	code := ""
	for i := 0; i < 6; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		code += fmt.Sprintf("%d", n.Int64())
	}
	return code
}
