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
	var created *entity.Category
	repo := &categorymock.MockCategoryRepository{
		CreateFn: func(category *entity.Category) error {
			created = category
			category.ID = 1
			return nil
		},
	}

	svc := service.NewCategoryService(repo)
	cat, err := svc.Create(request.CreateCategoryRequest{Name: "Fiqih"}, nil)

	assert.NoError(t, err)
	assert.NotNil(t, cat)
	assert.Equal(t, uint(1), cat.ID)
	assert.Equal(t, "Fiqih", cat.Name)
	assert.Same(t, created, cat, "returned pointer must be same as persisted object")
}

func TestCategoryCreate_RepoError(t *testing.T) {
	dbErr := errors.New("duplicate entry")
	repo := &categorymock.MockCategoryRepository{
		CreateFn: func(category *entity.Category) error { return dbErr },
	}

	svc := service.NewCategoryService(repo)
	cat, err := svc.Create(request.CreateCategoryRequest{Name: "Fiqih"}, nil)

	assert.ErrorIs(t, err, dbErr, "must propagate exact repo error")
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
	assert.Equal(t, "Fiqih", cats[0].Name)
	assert.Equal(t, "Tafsir", cats[1].Name)
}

func TestCategoryFindAll_RepoError(t *testing.T) {
	repoErr := errors.New("timeout")
	repo := &categorymock.MockCategoryRepository{
		FindAllFn: func() ([]entity.Category, error) { return nil, repoErr },
	}

	svc := service.NewCategoryService(repo)
	cats, err := svc.FindAll()

	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, cats)
}

func TestCategoryFindByID_Success(t *testing.T) {
	repo := &categorymock.MockCategoryRepository{
		FindByIDFn: func(id uint) (*entity.Category, error) {
			assert.Equal(t, uint(1), id)
			return &entity.Category{ID: id, Name: "Fiqih"}, nil
		},
	}

	svc := service.NewCategoryService(repo)
	cat, err := svc.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, cat)
	assert.Equal(t, uint(1), cat.ID)
	assert.Equal(t, "Fiqih", cat.Name)
}

func TestCategoryFindByID_NotFound(t *testing.T) {
	notFoundErr := errors.New("not found")
	repo := &categorymock.MockCategoryRepository{
		FindByIDFn: func(id uint) (*entity.Category, error) { return nil, notFoundErr },
	}

	svc := service.NewCategoryService(repo)
	cat, err := svc.FindByID(999)

	assert.ErrorIs(t, err, notFoundErr)
	assert.Nil(t, cat)
}

func TestCategoryUpdate_Success(t *testing.T) {
	var updatedID uint
	var updateMap map[string]interface{}
	repo := &categorymock.MockCategoryRepository{
		UpdateFn: func(id uint, updates map[string]interface{}) error {
			updatedID = id
			updateMap = updates
			return nil
		},
		FindByIDFn: func(id uint) (*entity.Category, error) {
			return &entity.Category{ID: id, Name: "Updated"}, nil
		},
	}

	svc := service.NewCategoryService(repo)
	cat, err := svc.Update(1, request.UpdateCategoryRequest{Name: "Updated"}, nil)

	assert.NoError(t, err)
	assert.NotNil(t, cat)
	assert.Equal(t, "Updated", cat.Name)
	assert.Equal(t, uint(1), updatedID, "must pass correct ID to repo.Update")
	assert.Equal(t, "Updated", updateMap["name"], "must pass name in updates map")
}

func TestCategoryUpdate_NoChanges(t *testing.T) {
	repo := &categorymock.MockCategoryRepository{}

	svc := service.NewCategoryService(repo)
	cat, err := svc.Update(1, request.UpdateCategoryRequest{}, nil)

	assert.Error(t, err)
	assert.Nil(t, cat)
	assert.Equal(t, "no updates provided", err.Error())
}

func TestCategoryUpdate_RepoUpdateError(t *testing.T) {
	repoErr := errors.New("constraint violation")
	repo := &categorymock.MockCategoryRepository{
		UpdateFn: func(id uint, updates map[string]interface{}) error { return repoErr },
	}

	svc := service.NewCategoryService(repo)
	cat, err := svc.Update(1, request.UpdateCategoryRequest{Name: "New"}, nil)

	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, cat)
}

func TestCategoryDelete_Success(t *testing.T) {
	var deletedID uint
	repo := &categorymock.MockCategoryRepository{
		DeleteFn: func(id uint) error {
			deletedID = id
			return nil
		},
	}

	svc := service.NewCategoryService(repo)
	err := svc.Delete(1)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), deletedID)
}

func TestCategoryDelete_RepoError(t *testing.T) {
	repoErr := errors.New("foreign key constraint")
	repo := &categorymock.MockCategoryRepository{
		DeleteFn: func(id uint) error { return repoErr },
	}

	svc := service.NewCategoryService(repo)
	err := svc.Delete(1)

	assert.ErrorIs(t, err, repoErr)
}

func TestCategorySearch_Success(t *testing.T) {
	repo := &categorymock.MockCategoryRepository{
		SearchFn: func(query string) ([]entity.Category, error) {
			assert.Equal(t, "Fiq", query, "must pass exact search query")
			return []entity.Category{{ID: 1, Name: "Fiqih"}}, nil
		},
	}

	svc := service.NewCategoryService(repo)
	cats, err := svc.Search("Fiq")

	assert.NoError(t, err)
	assert.Len(t, cats, 1)
	assert.Equal(t, "Fiqih", cats[0].Name)
}

func TestCategorySearch_Empty(t *testing.T) {
	repo := &categorymock.MockCategoryRepository{
		SearchFn: func(query string) ([]entity.Category, error) { return []entity.Category{}, nil },
	}

	svc := service.NewCategoryService(repo)
	cats, err := svc.Search("nonexistent")

	assert.NoError(t, err)
	assert.Empty(t, cats)
}

func TestCategorySearch_RepoError(t *testing.T) {
	repoErr := errors.New("search failed")
	repo := &categorymock.MockCategoryRepository{
		SearchFn: func(query string) ([]entity.Category, error) { return nil, repoErr },
	}

	svc := service.NewCategoryService(repo)
	cats, err := svc.Search("test")

	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, cats)
}
