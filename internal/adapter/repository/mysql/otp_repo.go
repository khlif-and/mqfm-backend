package mysql

import (
	"time"

	"gorm.io/gorm"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type otpRepo struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) port.OTPRepository {
	return &otpRepo{db: db}
}

func (r *otpRepo) Create(otp *entity.OTP) error {
	return r.db.Create(otp).Error
}

func (r *otpRepo) FindLatestByEmail(email string) (*entity.OTP, error) {
	var otp entity.OTP
	err := r.db.Where("email = ? AND verified = false AND expires_at > ?", email, time.Now()).
		Order("created_at DESC").
		First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *otpRepo) MarkVerified(id uint) error {
	return r.db.Model(&entity.OTP{}).Where("id = ?", id).Update("verified", true).Error
}

func (r *otpRepo) CountRecentByEmail(email string, since time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&entity.OTP{}).Where("email = ? AND created_at > ?", email, since).Count(&count).Error
	return count, err
}

func (r *otpRepo) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&entity.OTP{}).Error
}
