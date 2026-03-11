package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"mqfm-backend/internal/app/service"
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/shared/dto/request"
	categorymock "mqfm-backend/tests/mocks/category"
)

func TestCategoryCreate_Success(t *testing.T) {
	repo := &categorymock.MockCategoryRepository{
		CreateFn: func(category *entity.Category) error {
			category.ID = 1
			return nil
		},
	}

	svc := service.NewCategoryService(repo)
	cat, err := svc.Create(request.CreateCategoryRequest{Name: "Fiqih"}, nil)

	assert.NoError(t, err)
	assert.NotNil(t, cat)
	assert.Equal(t, "Fiqih", cat.Name)
}

func TestCategoryCreate_DuplicateName(t *testing.T) {
	repo := &categorymock.MockCategoryRepository{
		CreateFn: func(category *entity.Category) error {
			return errors.New("duplicate entry")
		},
	}

	svc := service.NewCategoryService(repo)
	cat, err := svc.Create(request.CreateCategoryRequest{Name: "Fiqih"}, nil)

	assert.Error(t, err)
	assert.Nil(t, cat)
}

func TestCategoryFindAll_Success(t *testing.T) {
	repo := &categorymock.MockCategoryRepository{
		FindAllFn: func() ([]entity.Category, error) {
			return []entity.Category{
				{ID: 1, Name: "Fiqih"},
				{ID: 2, Name: "Tafsir"},
			}, nil
		},
	}

	svc := service.NewCategoryService(repo)
	cats, err := svc.FindAll()

	assert.NoError(t, err)
	assert.Len(t, cats, 2)
}

func TestCategoryFindByID_Success(t *testing.T) {
	repo := &categorymock.MockCategoryRepository{
		FindByIDFn: func(id uint) (*entity.Category, error) {
			return &entity.Category{ID: id, Name: "Fiqih"}, nil
		},
	}

	svc := service.NewCategoryService(repo)
	cat, err := svc.FindByID(1)

	assert.NoError(t, err)
	assert.Equal(t, "Fiqih", cat.Name)
}

func TestCategoryFindByID_NotFound(t *testing.T) {
	repo := &categorymock.MockCategoryRepository{
		FindByIDFn: func(id uint) (*entity.Category, error) {
			return nil, errors.New("not found")
		},
	}

	svc := service.NewCategoryService(repo)
	cat, err := svc.FindByID(999)

	assert.Error(t, err)
	assert.Nil(t, cat)
}

func TestCategoryUpdate_Success(t *testing.T) {
	repo := &categorymock.MockCategoryRepository{
		UpdateFn: func(id uint, updates map[string]interface{}) error { return nil },
		FindByIDFn: func(id uint) (*entity.Category, error) {
			return &entity.Category{ID: id, Name: "Updated"}, nil
		},
	}

	svc := service.NewCategoryService(repo)
	cat, err := svc.Update(1, request.UpdateCategoryRequest{Name: "Updated"}, nil)

	assert.NoError(t, err)
	assert.Equal(t, "Updated", cat.Name)
}

func TestCategoryUpdate_NoChanges(t *testing.T) {
	repo := &categorymock.MockCategoryRepository{} 

	svc := service.NewCategoryService(repo)
	cat, err := svc.Update(1, request.UpdateCategoryRequest{}, nil)

	assert.Error(t, err)
	assert.Nil(t, cat)
	assert.Equal(t, "no updates provided", err.Error())
}

func TestCategoryDelete_Success(t *testing.T) {
	repo := &categorymock.MockCategoryRepository{
		DeleteFn: func(id uint) error { return nil },
	}

	svc := service.NewCategoryService(repo)
	err := svc.Delete(1)

	assert.NoError(t, err)
}

func TestCategorySearch_Success(t *testing.T) {
	repo := &categorymock.MockCategoryRepository{
		SearchFn: func(query string) ([]entity.Category, error) {
			return []entity.Category{{ID: 1, Name: "Fiqih"}}, nil
		},
	}

	svc := service.NewCategoryService(repo)
	cats, err := svc.Search("Fiq")

	assert.NoError(t, err)
	assert.Len(t, cats, 1)
}
