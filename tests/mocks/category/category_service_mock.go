package categorymock

import (
	"mime/multipart"
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
)

type MockCategoryService struct {
	CreateFn   func(req request.CreateCategoryRequest, file *multipart.FileHeader) (*entity.Category, error)
	FindAllFn  func() ([]entity.Category, error)
	FindByIDFn func(id uint) (*entity.Category, error)
	UpdateFn   func(id uint, req request.UpdateCategoryRequest, file *multipart.FileHeader) (*entity.Category, error)
	DeleteFn   func(id uint) error
	SearchFn   func(query string) ([]entity.Category, error)
}

func (m *MockCategoryService) Create(req request.CreateCategoryRequest, file *multipart.FileHeader) (*entity.Category, error) {
	return m.CreateFn(req, file)
}
func (m *MockCategoryService) FindAll() ([]entity.Category, error)           { return m.FindAllFn() }
func (m *MockCategoryService) FindByID(id uint) (*entity.Category, error)    { return m.FindByIDFn(id) }
func (m *MockCategoryService) Update(id uint, req request.UpdateCategoryRequest, file *multipart.FileHeader) (*entity.Category, error) {
	return m.UpdateFn(id, req, file)
}
func (m *MockCategoryService) Delete(id uint) error               { return m.DeleteFn(id) }
func (m *MockCategoryService) Search(query string) ([]entity.Category, error) {
	return m.SearchFn(query)
}

var _ port.CategoryService = (*MockCategoryService)(nil)
