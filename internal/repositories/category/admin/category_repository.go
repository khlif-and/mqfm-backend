package admin

import (
	"errors"
	"gorm.io/gorm"
	categoryModel "mqfm-backend/internal/models/category/admin"
)

type CategoryRepository interface {
	FindAll() ([]categoryModel.Category, error)
	FindByID(id uint) (*categoryModel.Category, error)
	Create(category *categoryModel.Category) error
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	Search(query string) ([]categoryModel.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll() ([]categoryModel.Category, error) {
	var categories []categoryModel.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) FindByID(id uint) (*categoryModel.Category, error) {
	var category categoryModel.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Create(category *categoryModel.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(id uint, updates map[string]interface{}) error {
	result := r.db.Model(&categoryModel.Category{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no records updated")
	}
	return nil
}

func (r *categoryRepository) Delete(id uint) error {
	result := r.db.Delete(&categoryModel.Category{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("category not found")
	}
	return nil
}

func (r *categoryRepository) Search(query string) ([]categoryModel.Category, error) {
	var categories []categoryModel.Category
	err := r.db.Where("name LIKE ?", "%"+query+"%").Find(&categories).Error
	return categories, err
}
