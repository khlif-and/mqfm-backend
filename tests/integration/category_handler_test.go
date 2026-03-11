package integration_test

import (
	"errors"
	"mime/multipart"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	admin "mqfm-backend/internal/adapter/handler/admin"
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/shared/dto/request"
	categorymock "mqfm-backend/tests/mocks/category"
	"mqfm-backend/tests/testutil"
)

func setupCategoryRouter(service *categorymock.MockCategoryService) *gin.Engine {
	r := testutil.SetupRouter()
	handler := admin.NewCategoryHandler(service)

	r.POST("/categories", handler.Create)
	r.GET("/categories", handler.FindAll)
	r.GET("/categories/:id", handler.FindByID)
	r.PUT("/categories/:id", handler.Update)
	r.DELETE("/categories/:id", handler.Delete)
	r.GET("/categories/search", handler.Search)
	return r
}

func TestCategoryHandler_Create_Success(t *testing.T) {
	service := &categorymock.MockCategoryService{
		CreateFn: func(req request.CreateCategoryRequest, file *multipart.FileHeader) (*entity.Category, error) {
			return &entity.Category{
				ID: 1, Name: req.Name,
				CreatedAt: time.Now(), UpdatedAt: time.Now(),
			}, nil
		},
	}

	r := setupCategoryRouter(service)
	req, _ := testutil.MakeMultipartRequest("POST", "/categories", map[string]string{"name": "Fiqih"})
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := testutil.ParseResponse(w)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "Fiqih", data["name"])
}

func TestCategoryHandler_Create_InvalidInput(t *testing.T) {
	service := &categorymock.MockCategoryService{}
	r := setupCategoryRouter(service)

	req, _ := testutil.MakeMultipartRequest("POST", "/categories", map[string]string{})
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCategoryHandler_FindAll_Success(t *testing.T) {
	service := &categorymock.MockCategoryService{
		FindAllFn: func() ([]entity.Category, error) {
			return []entity.Category{
				{ID: 1, Name: "Fiqih", CreatedAt: time.Now(), UpdatedAt: time.Now()},
				{ID: 2, Name: "Tafsir", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			}, nil
		},
	}

	r := setupCategoryRouter(service)
	req := testutil.MakeRequest("GET", "/categories", nil)
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutil.ParseResponse(w)
	data := resp["data"].([]interface{})
	assert.Len(t, data, 2)
}

func TestCategoryHandler_FindAll_Error(t *testing.T) {
	service := &categorymock.MockCategoryService{
		FindAllFn: func() ([]entity.Category, error) {
			return nil, errors.New("db error")
		},
	}

	r := setupCategoryRouter(service)
	req := testutil.MakeRequest("GET", "/categories", nil)
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCategoryHandler_FindByID_Success(t *testing.T) {
	service := &categorymock.MockCategoryService{
		FindByIDFn: func(id uint) (*entity.Category, error) {
			return &entity.Category{
				ID: id, Name: "Fiqih",
				CreatedAt: time.Now(), UpdatedAt: time.Now(),
			}, nil
		},
	}

	r := setupCategoryRouter(service)
	req := testutil.MakeRequest("GET", "/categories/1", nil)
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCategoryHandler_FindByID_InvalidID(t *testing.T) {
	service := &categorymock.MockCategoryService{}
	r := setupCategoryRouter(service)

	req := testutil.MakeRequest("GET", "/categories/abc", nil)
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCategoryHandler_FindByID_NotFound(t *testing.T) {
	service := &categorymock.MockCategoryService{
		FindByIDFn: func(id uint) (*entity.Category, error) {
			return nil, errors.New("not found")
		},
	}

	r := setupCategoryRouter(service)
	req := testutil.MakeRequest("GET", "/categories/999", nil)
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCategoryHandler_Delete_Success(t *testing.T) {
	service := &categorymock.MockCategoryService{
		DeleteFn: func(id uint) error { return nil },
	}

	r := setupCategoryRouter(service)
	req := testutil.MakeRequest("DELETE", "/categories/1", nil)
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCategoryHandler_Delete_Error(t *testing.T) {
	service := &categorymock.MockCategoryService{
		DeleteFn: func(id uint) error { return errors.New("cannot delete") },
	}

	r := setupCategoryRouter(service)
	req := testutil.MakeRequest("DELETE", "/categories/1", nil)
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCategoryHandler_Search_Success(t *testing.T) {
	service := &categorymock.MockCategoryService{
		SearchFn: func(query string) ([]entity.Category, error) {
			return []entity.Category{
				{ID: 1, Name: "Fiqih", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			}, nil
		},
	}

	r := setupCategoryRouter(service)
	req := testutil.MakeRequest("GET", "/categories/search?q=Fiq", nil)
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCategoryHandler_Search_Empty(t *testing.T) {
	service := &categorymock.MockCategoryService{}
	r := setupCategoryRouter(service)

	req := testutil.MakeRequest("GET", "/categories/search", nil)
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
