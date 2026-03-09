package service

import (
	"context"
	"errors"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/internal/shared/logger"
	"mqfm-backend/internal/shared/security"
)

type adminAuthService struct {
	repo       port.AdminRepository
	tokenStore port.TokenStore
}

func NewAdminAuthService(repo port.AdminRepository, tokenStore port.TokenStore) port.AdminAuthService {
	return &adminAuthService{repo: repo, tokenStore: tokenStore}
}

func (s *adminAuthService) Register(req request.AdminRegisterRequest) (*entity.Admin, error) {
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		logger.Error("failed to hash admin password")
		return nil, err
	}

	admin := entity.Admin{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     constant.RoleAdmin,
	}

	if err := s.repo.Create(&admin); err != nil {
		return nil, err
	}

	return &admin, nil
}

func (s *adminAuthService) Login(req request.AdminLoginRequest) (string, *entity.Admin, error) {
	admin, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		logger.Warn("admin login: email not found")
		return "", nil, errors.New("invalid admin credentials")
	}

	if !security.CheckPassword(req.Password, admin.Password) {
		logger.Warn("admin login: incorrect password")
		return "", nil, errors.New("invalid admin credentials")
	}

	token, err := security.GenerateToken(admin.ID, constant.RoleAdmin)
	if err != nil {
		logger.Error("failed to generate admin token")
		return "", nil, err
	}

	if s.tokenStore != nil {
		_ = s.tokenStore.StoreToken(context.Background(), admin.ID, constant.RoleAdmin, token, security.TokenTTL)
	}

	return token, admin, nil
}

func (s *adminAuthService) UpdateAdmin(id uint, updates map[string]interface{}) (*entity.Admin, error) {
	if pwd, ok := updates["password"].(string); ok && pwd != "" {
		hashed, err := security.HashPassword(pwd)
		if err != nil {
			return nil, err
		}
		updates["password"] = hashed
	}

	if err := s.repo.Update(id, updates); err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

func (s *adminAuthService) GetAdminByID(id uint) (*entity.Admin, error) {
	return s.repo.FindByID(id)
}
