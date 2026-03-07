package mysql

import (
	"errors"

	"gorm.io/gorm"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type categoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) port.CategoryRepository {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) FindAll() ([]entity.Category, error) {
	var categories []entity.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *categoryRepo) FindByID(id uint) (*entity.Category, error) {
	var category entity.Category
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepo) Create(category *entity.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepo) Update(id uint, updates map[string]interface{}) error {
	result := r.db.Model(&entity.Category{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no records updated")
	}
	return nil
}

func (r *categoryRepo) Delete(id uint) error {
	result := r.db.Delete(&entity.Category{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("category not found")
	}
	return nil
}

func (r *categoryRepo) Search(query string) ([]entity.Category, error) {
	var categories []entity.Category
	err := r.db.Where("name LIKE ?", "%"+query+"%").Find(&categories).Error
	return categories, err
}
