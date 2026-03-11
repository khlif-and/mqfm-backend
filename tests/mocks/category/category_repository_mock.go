package categorymock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type MockCategoryRepository struct {
	FindAllFn  func() ([]entity.Category, error)
	FindByIDFn func(id uint) (*entity.Category, error)
	CreateFn   func(category *entity.Category) error
	UpdateFn   func(id uint, updates map[string]interface{}) error
	DeleteFn   func(id uint) error
	SearchFn   func(query string) ([]entity.Category, error)
}

func (m *MockCategoryRepository) FindAll() ([]entity.Category, error) { return m.FindAllFn() }
func (m *MockCategoryRepository) FindByID(id uint) (*entity.Category, error) {
	return m.FindByIDFn(id)
}
func (m *MockCategoryRepository) Create(category *entity.Category) error {
	return m.CreateFn(category)
}
func (m *MockCategoryRepository) Update(id uint, updates map[string]interface{}) error {
	return m.UpdateFn(id, updates)
}
func (m *MockCategoryRepository) Delete(id uint) error            { return m.DeleteFn(id) }
func (m *MockCategoryRepository) Search(query string) ([]entity.Category, error) {
	return m.SearchFn(query)
}

var _ port.CategoryRepository = (*MockCategoryRepository)(nil)
