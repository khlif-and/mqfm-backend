package admin

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	adminModel "mqfm-backend/internal/models/auth/admin"
	"mqfm-backend/internal/utils"

)

type AdminAuthService struct {
	db *gorm.DB
}

func NewAdminAuthService(db *gorm.DB) *AdminAuthService {
	return &AdminAuthService{db: db}
}

func (s *AdminAuthService) Register(admin *adminModel.Admin) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Log.Error("Failed to hash admin password")
		return err
	}
	admin.Password = string(hashedPassword)
	admin.Role = "admin"
	return s.db.Create(admin).Error
}

func (s *AdminAuthService) Login(email, password string) (string, error) {
	var admin adminModel.Admin
	if err := s.db.Where("email = ?", email).First(&admin).Error; err != nil {
		utils.Log.Warn("Admin login attempt failed: email not found")
		return "", errors.New("invalid admin credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		utils.Log.Warn("Admin login attempt failed: incorrect password")
		return "", errors.New("invalid admin credentials")
	}

	token, err := utils.GenerateToken(admin.ID, "admin")
	if err != nil {
		utils.Log.Error("Failed to generate admin JWT token: " + err.Error())
		return "", err
	}

	return token, nil
}

func (s *AdminAuthService) UpdateAdmin(id uint, updates map[string]interface{}) error {
	if pwd, ok := updates["password"].(string); ok && pwd != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			utils.Log.Error("Failed to hash new admin password")
			return err
		}
		updates["password"] = string(hashed)
	}

	return s.db.Model(&adminModel.Admin{}).Where("id = ?", id).Updates(updates).Error
}