package mocks

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
	"mime/multipart"
)

type MockOTPService struct {
	SendOTPFn   func(email string) error
	VerifyOTPFn func(email string, code string) (*entity.User, error)
}

func (m *MockOTPService) SendOTP(email string) error              { return m.SendOTPFn(email) }
func (m *MockOTPService) VerifyOTP(email, code string) (*entity.User, error) {
	return m.VerifyOTPFn(email, code)
}

type MockEmailService struct {
	SendAsyncFn func(to, subject, body string)
}

func (m *MockEmailService) SendAsync(to, subject, body string) {
	if m.SendAsyncFn != nil {
		m.SendAsyncFn(to, subject, body)
	}
}

type MockColorExtractorService struct {
	ExtractFn func(imagePath string) (string, error)
}

func (m *MockColorExtractorService) ExtractDominantColor(imagePath string) (string, error) {
	return m.ExtractFn(imagePath)
}

type MockAdminAuthService struct {
	RegisterFn   func(req request.AdminRegisterRequest) (*entity.Admin, error)
	LoginFn      func(req request.AdminLoginRequest) (string, *entity.Admin, error)
	UpdateFn     func(id uint, updates map[string]interface{}) (*entity.Admin, error)
	GetByIDFn    func(id uint) (*entity.Admin, error)
}

func (m *MockAdminAuthService) Register(req request.AdminRegisterRequest) (*entity.Admin, error) {
	return m.RegisterFn(req)
}
func (m *MockAdminAuthService) Login(req request.AdminLoginRequest) (string, *entity.Admin, error) {
	return m.LoginFn(req)
}
func (m *MockAdminAuthService) UpdateAdmin(id uint, updates map[string]interface{}) (*entity.Admin, error) {
	return m.UpdateFn(id, updates)
}
func (m *MockAdminAuthService) GetAdminByID(id uint) (*entity.Admin, error) {
	return m.GetByIDFn(id)
}

type MockCategoryService struct {
	CreateFn  func(req request.CreateCategoryRequest, file *multipart.FileHeader) (*entity.Category, error)
	FindAllFn func() ([]entity.Category, error)
	FindByIDFn func(id uint) (*entity.Category, error)
	UpdateFn  func(id uint, req request.UpdateCategoryRequest, file *multipart.FileHeader) (*entity.Category, error)
	DeleteFn  func(id uint) error
	SearchFn  func(query string) ([]entity.Category, error)
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

type MockAudioService struct {
	CreateFn  func(req request.CreateAudioRequest, audio *multipart.FileHeader, thumb *multipart.FileHeader) (*entity.Audio, error)
	FindAllFn func() ([]entity.Audio, error)
	FindByIDFn func(id uint) (*entity.Audio, error)
	UpdateFn  func(id uint, req request.UpdateAudioRequest, audio *multipart.FileHeader, thumb *multipart.FileHeader) (*entity.Audio, error)
	DeleteFn  func(id uint) error
	SearchFn  func(query string) ([]entity.Audio, error)
}

func (m *MockAudioService) Create(req request.CreateAudioRequest, audio *multipart.FileHeader, thumb *multipart.FileHeader) (*entity.Audio, error) {
	return m.CreateFn(req, audio, thumb)
}
func (m *MockAudioService) FindAll() ([]entity.Audio, error)           { return m.FindAllFn() }
func (m *MockAudioService) FindByID(id uint) (*entity.Audio, error)    { return m.FindByIDFn(id) }
func (m *MockAudioService) Update(id uint, req request.UpdateAudioRequest, audio *multipart.FileHeader, thumb *multipart.FileHeader) (*entity.Audio, error) {
	return m.UpdateFn(id, req, audio, thumb)
}
func (m *MockAudioService) Delete(id uint) error                       { return m.DeleteFn(id) }
func (m *MockAudioService) Search(query string) ([]entity.Audio, error) { return m.SearchFn(query) }

type MockHistoryService struct {
	RecordPlayFn    func(userID uint, req request.HistoryRequest) error
	GetHistoryFn    func(userID uint) ([]entity.History, error)
	DeleteHistoryFn func(userID, audioID uint) error
	ClearHistoryFn  func(userID uint) error
}

func (m *MockHistoryService) RecordPlay(userID uint, req request.HistoryRequest) error {
	return m.RecordPlayFn(userID, req)
}
func (m *MockHistoryService) GetHistory(userID uint) ([]entity.History, error) {
	return m.GetHistoryFn(userID)
}
func (m *MockHistoryService) DeleteHistory(userID, audioID uint) error {
	return m.DeleteHistoryFn(userID, audioID)
}
func (m *MockHistoryService) ClearHistory(userID uint) error { return m.ClearHistoryFn(userID) }

type MockLikeService struct {
	LikeFn          func(userID uint, req request.LikeRequest) (*entity.Like, error)
	UnlikeFn        func(userID uint, req request.UnlikeRequest) error
	GetLikesFn      func(userID string, targetType string) ([]entity.Like, error)
	CountByTargetFn func(targetType string, targetID uint) (int64, error)
}

func (m *MockLikeService) Like(userID uint, req request.LikeRequest) (*entity.Like, error) {
	return m.LikeFn(userID, req)
}
func (m *MockLikeService) Unlike(userID uint, req request.UnlikeRequest) error {
	return m.UnlikeFn(userID, req)
}
func (m *MockLikeService) GetLikes(userID string, targetType string) ([]entity.Like, error) {
	return m.GetLikesFn(userID, targetType)
}
func (m *MockLikeService) CountByTarget(targetType string, targetID uint) (int64, error) {
	return m.CountByTargetFn(targetType, targetID)
}

var _ port.OTPService = (*MockOTPService)(nil)
var _ port.EmailService = (*MockEmailService)(nil)
var _ port.ColorExtractorService = (*MockColorExtractorService)(nil)
var _ port.AdminAuthService = (*MockAdminAuthService)(nil)
var _ port.CategoryService = (*MockCategoryService)(nil)
var _ port.AudioService = (*MockAudioService)(nil)
var _ port.HistoryService = (*MockHistoryService)(nil)
var _ port.LikeService = (*MockLikeService)(nil)
