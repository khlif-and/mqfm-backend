package mysql

import (
	"errors"

	"gorm.io/gorm"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type adminRepo struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) port.AdminRepository {
	return &adminRepo{db: db}
}

func (r *adminRepo) Create(admin *entity.Admin) error {
	return r.db.Create(admin).Error
}

func (r *adminRepo) FindByEmail(email string) (*entity.Admin, error) {
	var admin entity.Admin
	if err := r.db.Where("email = ?", email).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepo) FindByID(id uint) (*entity.Admin, error) {
	var admin entity.Admin
	if err := r.db.First(&admin, id).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepo) Update(id uint, updates map[string]interface{}) error {
	result := r.db.Model(&entity.Admin{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no records updated")
	}
	return nil
}
