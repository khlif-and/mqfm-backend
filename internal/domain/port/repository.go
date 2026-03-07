package port

import (
	"mqfm-backend/internal/domain/entity"
)

type AdminRepository interface {
	Create(admin *entity.Admin) error
	FindByEmail(email string) (*entity.Admin, error)
	FindByID(id uint) (*entity.Admin, error)
	Update(id uint, updates map[string]interface{}) error
}

type UserRepository interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	FindByID(id uint) (*entity.User, error)
	FindByProviderID(provider string, providerID string) (*entity.User, error)
	Update(id uint, updates map[string]interface{}) error
}

type CategoryRepository interface {
	FindAll() ([]entity.Category, error)
	FindByID(id uint) (*entity.Category, error)
	Create(category *entity.Category) error
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	Search(query string) ([]entity.Category, error)
}

type AudioRepository interface {
	FindAll() ([]entity.Audio, error)
	FindByID(id uint) (*entity.Audio, error)
	Create(audio *entity.Audio) error
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	Search(query string) ([]entity.Audio, error)
}

type PlaylistRepository interface {
	Create(playlist *entity.Playlist) error
	FindByUserID(userID uint) ([]entity.Playlist, error)
	FindByID(id uint, userID uint) (*entity.Playlist, error)
	Search(userID uint, query string) ([]entity.Playlist, error)
	AddAudio(playlist *entity.Playlist, audio *entity.Audio) error
	FindAudioByID(id uint) (*entity.Audio, error)
}

type LikeRepository interface {
	Create(like *entity.Like) error
	Delete(userID, audioID uint) error
	FindByUser(userID uint) ([]entity.Like, error)
	Exists(userID, audioID uint) (bool, error)
}

type HistoryRepository interface {
	Upsert(history *entity.History) error
	FindByUser(userID uint) ([]entity.History, error)
	DeleteByUserAndAudio(userID, audioID uint) error
	DeleteAllByUser(userID uint) error
}

type LivestreamRepository interface {
	FindFirst() (*entity.LiveStream, error)
	Create(ls *entity.LiveStream) error
	Save(ls *entity.LiveStream) error
}
