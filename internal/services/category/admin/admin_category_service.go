package admin

import (
	"errors"

	"gorm.io/gorm"

	categoryModel "mqfm-backend/internal/models/category/admin"

)

type AdminCategoryService struct {
	db *gorm.DB
}

func NewAdminCategoryService(db *gorm.DB) *AdminCategoryService {
	return &AdminCategoryService{db: db}
}

func (s *AdminCategoryService) Create(category *categoryModel.Category) error {
	return s.db.Create(category).Error
}

func (s *AdminCategoryService) FindAll() ([]categoryModel.Category, error) {
	var categories []categoryModel.Category
	if err := s.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *AdminCategoryService) FindByID(id uint) (*categoryModel.Category, error) {
	var category categoryModel.Category
	if err := s.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

func (s *AdminCategoryService) Update(id uint, updates map[string]interface{}) (*categoryModel.Category, error) {
	if err := s.db.Model(&categoryModel.Category{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}

	var updatedCategory categoryModel.Category
	if err := s.db.First(&updatedCategory, id).Error; err != nil {
		return nil, err
	}

	return &updatedCategory, nil
}

func (s *AdminCategoryService) Delete(id uint) error {
	var category categoryModel.Category
	if err := s.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return err
	}

	return s.db.Delete(&category).Error
}