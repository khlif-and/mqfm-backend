package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"mqfm-backend/internal/app/service"
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/internal/shared/security"
	"mqfm-backend/tests/mocks"
)

type stubTokenStore struct{}

func (s *stubTokenStore) StoreToken(_ context.Context, _ uint, _ string, _ string, _ time.Duration) error {
	return nil
}
func (s *stubTokenStore) GetToken(_ context.Context, _ uint, _ string) (string, error) {
	return "", nil
}
func (s *stubTokenStore) DeleteToken(_ context.Context, _ uint, _ string) error { return nil }
func (s *stubTokenStore) RefreshToken(_ context.Context, _ uint, _ string) (string, error) {
	return "", nil
}

func TestAdminRegister_Success(t *testing.T) {
	repo := &mocks.MockAdminRepository{
		CreateFn: func(admin *entity.Admin) error { admin.ID = 1; return nil },
	}
	svc := service.NewAdminAuthService(repo, &stubTokenStore{})
	admin, err := svc.Register(request.AdminRegisterRequest{Username: "admin1", Email: "admin@test.com", Password: "password123"})
	assert.NoError(t, err)
	assert.NotNil(t, admin)
	assert.Equal(t, "admin1", admin.Username)
	assert.Equal(t, "admin@test.com", admin.Email)
	assert.Equal(t, "admin", admin.Role)
	assert.True(t, security.CheckPassword("password123", admin.Password))
}

func TestAdminRegister_DuplicateEmail(t *testing.T) {
	repo := &mocks.MockAdminRepository{
		CreateFn: func(admin *entity.Admin) error { return errors.New("duplicate entry") },
	}
	svc := service.NewAdminAuthService(repo, &stubTokenStore{})
	admin, err := svc.Register(request.AdminRegisterRequest{Username: "admin1", Email: "admin@test.com", Password: "password123"})
	assert.Error(t, err)
	assert.Nil(t, admin)
}

func TestAdminLogin_Success(t *testing.T) {
	hashed, _ := security.HashPassword("password123")
	repo := &mocks.MockAdminRepository{
		FindByEmailFn: func(email string) (*entity.Admin, error) {
			return &entity.Admin{ID: 1, Username: "admin1", Email: email, Password: hashed, Role: "admin"}, nil
		},
	}
	svc := service.NewAdminAuthService(repo, &stubTokenStore{})
	token, admin, err := svc.Login(request.AdminLoginRequest{Email: "admin@test.com", Password: "password123"})
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, uint(1), admin.ID)
}

func TestAdminLogin_WrongPassword(t *testing.T) {
	hashed, _ := security.HashPassword("password123")
	repo := &mocks.MockAdminRepository{
		FindByEmailFn: func(email string) (*entity.Admin, error) {
			return &entity.Admin{ID: 1, Password: hashed}, nil
		},
	}
	svc := service.NewAdminAuthService(repo, &stubTokenStore{})
	token, admin, err := svc.Login(request.AdminLoginRequest{Email: "admin@test.com", Password: "wrongpassword"})
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Nil(t, admin)
}

func TestAdminLogin_UserNotFound(t *testing.T) {
	repo := &mocks.MockAdminRepository{
		FindByEmailFn: func(email string) (*entity.Admin, error) { return nil, errors.New("not found") },
	}
	svc := service.NewAdminAuthService(repo, &stubTokenStore{})
	token, admin, err := svc.Login(request.AdminLoginRequest{Email: "unknown@test.com", Password: "password123"})
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Nil(t, admin)
}

func TestAdminUpdateAdmin_Success(t *testing.T) {
	repo := &mocks.MockAdminRepository{
		UpdateFn:   func(id uint, updates map[string]interface{}) error { return nil },
		FindByIDFn: func(id uint) (*entity.Admin, error) { return &entity.Admin{ID: id, Username: "updated", Email: "admin@test.com"}, nil },
	}
	svc := service.NewAdminAuthService(repo, &stubTokenStore{})
	admin, err := svc.UpdateAdmin(1, map[string]interface{}{"username": "updated"})
	assert.NoError(t, err)
	assert.Equal(t, "updated", admin.Username)
}

func TestAdminUpdateAdmin_WithPassword(t *testing.T) {
	var savedUpdates map[string]interface{}
	repo := &mocks.MockAdminRepository{
		UpdateFn:   func(id uint, updates map[string]interface{}) error { savedUpdates = updates; return nil },
		FindByIDFn: func(id uint) (*entity.Admin, error) { return &entity.Admin{ID: id}, nil },
	}
	svc := service.NewAdminAuthService(repo, &stubTokenStore{})
	_, err := svc.UpdateAdmin(1, map[string]interface{}{"password": "newpassword"})
	assert.NoError(t, err)
	hashedPwd, ok := savedUpdates["password"].(string)
	assert.True(t, ok)
	assert.True(t, security.CheckPassword("newpassword", hashedPwd))
}

func TestAdminGetByID_Success(t *testing.T) {
	repo := &mocks.MockAdminRepository{
		FindByIDFn: func(id uint) (*entity.Admin, error) { return &entity.Admin{ID: id, Username: "admin1"}, nil },
	}
	svc := service.NewAdminAuthService(repo, &stubTokenStore{})
	admin, err := svc.GetAdminByID(1)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), admin.ID)
}

func TestAdminGetByID_NotFound(t *testing.T) {
	repo := &mocks.MockAdminRepository{
		FindByIDFn: func(id uint) (*entity.Admin, error) { return nil, errors.New("not found") },
	}
	svc := service.NewAdminAuthService(repo, &stubTokenStore{})
	admin, err := svc.GetAdminByID(999)
	assert.Error(t, err)
	assert.Nil(t, admin)
}
