package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	userModel "mqfm-backend/internal/models/auth/user"
	"mqfm-backend/internal/utils"

)

type UserAuthService struct {
	db *gorm.DB
}

func NewUserAuthService(db *gorm.DB) *UserAuthService {
	return &UserAuthService{db: db}
}

func (s *UserAuthService) Register(user *userModel.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Log.Error("Failed to hash user password")
		return err
	}
	user.Password = string(hashedPassword)
	user.Role = "user"
	return s.db.Create(user).Error
}

func (s *UserAuthService) Login(email, password string) (string, *userModel.User, error) {
	var user userModel.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		utils.Log.Warn("User login attempt failed: email not found")
		return "", nil, errors.New("invalid user credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		utils.Log.Warn("User login attempt failed: incorrect password")
		return "", nil, errors.New("invalid user credentials")
	}

	token, err := utils.GenerateToken(user.ID, "user")
	if err != nil {
		utils.Log.Error("Failed to generate user JWT token: " + err.Error())
		return "", nil, err
	}

	return token, &user, nil
}

func (s *UserAuthService) UpdateUser(id uint, updates map[string]interface{}) (*userModel.User, error) {
	if pwd, ok := updates["password"].(string); ok && pwd != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			utils.Log.Error("Failed to hash new user password")
			return nil, err
		}
		updates["password"] = string(hashed)
	}

	if err := s.db.Model(&userModel.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}

	var updatedUser userModel.User
	if err := s.db.First(&updatedUser, id).Error; err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (s *UserAuthService) GetUserByID(id uint) (*userModel.User, error) {
	var user userModel.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}