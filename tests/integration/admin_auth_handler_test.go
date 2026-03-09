package integration_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	admin "mqfm-backend/internal/adapter/handler/admin"
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/infrastructure/middleware"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/tests/mocks"
	"mqfm-backend/tests/testutil"
)

func setupAdminAuthRouter(service *mocks.MockAdminAuthService) *gin.Engine {
	r := testutil.SetupRouter()
	handler := admin.NewAuthHandler(service)

	r.POST("/admin/register", handler.Register)
	r.POST("/admin/login", handler.Login)
	r.GET("/admin/me", middleware.JWTAuth(), handler.Me)
	r.PUT("/admin/:id", middleware.JWTAuth(), handler.Update)
	r.POST("/admin/logout", middleware.JWTAuth(), handler.Logout)
	return r
}

func TestAdminAuthHandler_Register_Success(t *testing.T) {
	service := &mocks.MockAdminAuthService{
		RegisterFn: func(req request.AdminRegisterRequest) (*entity.Admin, error) {
			return &entity.Admin{
				ID: 1, Username: req.Username, Email: req.Email,
				Role: "admin", CreatedAt: time.Now(), UpdatedAt: time.Now(),
			}, nil
		},
	}

	r := setupAdminAuthRouter(service)
	req := testutil.MakeRequest("POST", "/admin/register", map[string]string{
		"username": "admin1", "email": "admin@test.com", "password": "password123",
	})
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := testutil.ParseResponse(w)
	assert.Equal(t, float64(http.StatusCreated), resp["status"])
}

func TestAdminAuthHandler_Register_InvalidInput(t *testing.T) {
	service := &mocks.MockAdminAuthService{}
	r := setupAdminAuthRouter(service)

	req := testutil.MakeRequest("POST", "/admin/register", map[string]string{
		"email": "invalid",
	})
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminAuthHandler_Login_Success(t *testing.T) {
	service := &mocks.MockAdminAuthService{
		LoginFn: func(req request.AdminLoginRequest) (string, *entity.Admin, error) {
			return "test-token", &entity.Admin{
				ID: 1, Username: "admin1", Email: req.Email,
				Role: "admin", CreatedAt: time.Now(), UpdatedAt: time.Now(),
			}, nil
		},
	}

	r := setupAdminAuthRouter(service)
	req := testutil.MakeRequest("POST", "/admin/login", map[string]string{
		"email": "admin@test.com", "password": "password123",
	})
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutil.ParseResponse(w)
	data := resp["data"].(map[string]interface{})
	assert.NotEmpty(t, data["token"])
}

func TestAdminAuthHandler_Login_Fail(t *testing.T) {
	service := &mocks.MockAdminAuthService{
		LoginFn: func(req request.AdminLoginRequest) (string, *entity.Admin, error) {
			return "", nil, errors.New("invalid credentials")
		},
	}

	r := setupAdminAuthRouter(service)
	req := testutil.MakeRequest("POST", "/admin/login", map[string]string{
		"email": "admin@test.com", "password": "wrong",
	})
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminAuthHandler_Me_Success(t *testing.T) {
	service := &mocks.MockAdminAuthService{
		GetByIDFn: func(id uint) (*entity.Admin, error) {
			return &entity.Admin{
				ID: id, Username: "admin1", Email: "admin@test.com",
				Role: "admin", CreatedAt: time.Now(), UpdatedAt: time.Now(),
			}, nil
		},
	}

	r := setupAdminAuthRouter(service)
	req := testutil.MakeAuthRequest("GET", "/admin/me", nil, 1, "admin")
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutil.ParseResponse(w)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "admin1", data["username"])
}

func TestAdminAuthHandler_Me_Unauthorized(t *testing.T) {
	service := &mocks.MockAdminAuthService{}
	r := setupAdminAuthRouter(service)

	req := testutil.MakeRequest("GET", "/admin/me", nil)
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminAuthHandler_Logout(t *testing.T) {
	service := &mocks.MockAdminAuthService{}
	r := setupAdminAuthRouter(service)

	req := testutil.MakeAuthRequest("POST", "/admin/logout", nil, 1, "admin")
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminAuthHandler_Update_Success(t *testing.T) {
	service := &mocks.MockAdminAuthService{
		UpdateFn: func(id uint, updates map[string]interface{}) (*entity.Admin, error) {
			return &entity.Admin{
				ID: id, Username: "updated", Email: "admin@test.com",
				Role: "admin", CreatedAt: time.Now(), UpdatedAt: time.Now(),
			}, nil
		},
	}

	r := setupAdminAuthRouter(service)
	req := testutil.MakeAuthRequest("PUT", "/admin/1", map[string]interface{}{
		"username": "updated",
	}, 1, "admin")
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminAuthHandler_Update_InvalidID(t *testing.T) {
	service := &mocks.MockAdminAuthService{}
	r := setupAdminAuthRouter(service)

	req := testutil.MakeAuthRequest("PUT", "/admin/abc", map[string]interface{}{
		"username": "updated",
	}, 1, "admin")
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
