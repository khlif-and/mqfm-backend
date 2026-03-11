package audiomock

import (
	"mime/multipart"
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
)

type MockAudioService struct {
	CreateFn   func(req request.CreateAudioRequest, audio *multipart.FileHeader, thumb *multipart.FileHeader) (*entity.Audio, error)
	FindAllFn  func() ([]entity.Audio, error)
	FindByIDFn func(id uint) (*entity.Audio, error)
	UpdateFn   func(id uint, req request.UpdateAudioRequest, audio *multipart.FileHeader, thumb *multipart.FileHeader) (*entity.Audio, error)
	DeleteFn   func(id uint) error
	SearchFn   func(query string) ([]entity.Audio, error)
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

var _ port.AudioService = (*MockAudioService)(nil)
