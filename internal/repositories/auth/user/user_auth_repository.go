package user

import (
	"errors"
	"gorm.io/gorm"
	userModel "mqfm-backend/internal/models/auth/user"
)

type UserAuthRepository interface {
	Create(user *userModel.User) error
	FindByEmail(email string) (*userModel.User, error)
	FindByID(id uint) (*userModel.User, error)
	Update(id uint, updates map[string]interface{}) error
}

type userAuthRepository struct {
	db *gorm.DB
}

func NewUserAuthRepository(db *gorm.DB) UserAuthRepository {
	return &userAuthRepository{db: db}
}

func (r *userAuthRepository) Create(user *userModel.User) error {
	return r.db.Create(user).Error
}

func (r *userAuthRepository) FindByEmail(email string) (*userModel.User, error) {
	var user userModel.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userAuthRepository) FindByID(id uint) (*userModel.User, error) {
	var user userModel.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userAuthRepository) Update(id uint, updates map[string]interface{}) error {
	result := r.db.Model(&userModel.User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no records updated")
	}
	return nil
}
