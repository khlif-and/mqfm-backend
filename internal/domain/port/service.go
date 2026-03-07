package port

import (
	"mime/multipart"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/shared/dto/request"
)

type AdminAuthService interface {
	Register(req request.AdminRegisterRequest) (*entity.Admin, error)
	Login(req request.AdminLoginRequest) (string, *entity.Admin, error)
	UpdateAdmin(id uint, updates map[string]interface{}) (*entity.Admin, error)
	GetAdminByID(id uint) (*entity.Admin, error)
}

type UserAuthService interface {
	Register(req request.UserRegisterRequest, file *multipart.FileHeader) (*entity.User, error)
	Login(req request.UserLoginRequest) (string, *entity.User, error)
	GoogleLogin(req request.GoogleLoginRequest) (string, *entity.User, error)
	UpdateUser(id uint, req request.UpdateUserRequest, file *multipart.FileHeader) (*entity.User, error)
	GetUserByID(id uint) (*entity.User, error)
}

type CategoryService interface {
	Create(req request.CreateCategoryRequest, file *multipart.FileHeader) (*entity.Category, error)
	FindAll() ([]entity.Category, error)
	FindByID(id uint) (*entity.Category, error)
	Update(id uint, req request.UpdateCategoryRequest, file *multipart.FileHeader) (*entity.Category, error)
	Delete(id uint) error
	Search(query string) ([]entity.Category, error)
}

type AudioService interface {
	Create(req request.CreateAudioRequest, file *multipart.FileHeader) (*entity.Audio, error)
	FindAll() ([]entity.Audio, error)
	FindByID(id uint) (*entity.Audio, error)
	Update(id uint, req request.UpdateAudioRequest, file *multipart.FileHeader) (*entity.Audio, error)
	Delete(id uint) error
	Search(query string) ([]entity.Audio, error)
}

type PlaylistService interface {
	Create(playlist *entity.Playlist) error
	GetByUserID(userID uint) ([]entity.Playlist, error)
	GetByID(id uint, userID uint) (*entity.Playlist, error)
	Search(userID uint, query string) ([]entity.Playlist, error)
	AddAudioToPlaylist(userID, playlistID, audioID uint) error
}

type LikeService interface {
	LikeAudio(userID uint, req request.LikeRequest) (*entity.Like, error)
	UnlikeAudio(userID, audioID uint) error
	GetLikedAudios(userID uint) ([]entity.Like, error)
}

type HistoryService interface {
	RecordPlay(userID uint, req request.HistoryRequest) error
	GetHistory(userID uint) ([]entity.History, error)
	DeleteHistory(userID, audioID uint) error
	ClearHistory(userID uint) error
}

type LivestreamService interface {
	UpdateLiveStatus(channelID string) error
	GetStatus() (*entity.LiveStream, error)
}
