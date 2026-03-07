package port

import (
	"mqfm-backend/internal/domain/entity"
	"time"
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
	FindByIDs(ids []uint) ([]entity.Audio, error)
	FindByArtist(artist string, limit int) ([]entity.Audio, error)
	FindByCategoryID(categoryID uint, limit int) ([]entity.Audio, error)
	FindAllActive() ([]entity.Audio, error)
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
	CountByAudio(audioID uint) (int64, error)
	AggregateLikeCounts() (map[uint]int64, error)
}

type HistoryRepository interface {
	Upsert(history *entity.History) error
	FindByUser(userID uint) ([]entity.History, error)
	DeleteByUserAndAudio(userID, audioID uint) error
	DeleteAllByUser(userID uint) error
	CountByAudio(audioID uint) (int64, error)
	FindFrequentByUser(userID uint, minPlays int, limit int) ([]entity.History, error)
	AggregatePlayCounts() (map[uint]int64, error)
}

type OTPRepository interface {
	Create(otp *entity.OTP) error
	FindLatestByEmail(email string) (*entity.OTP, error)
	MarkVerified(id uint) error
	CountRecentByEmail(email string, since time.Time) (int64, error)
	DeleteExpired() error
}

type AudioScoreRepository interface {
	Upsert(score *entity.AudioScore) error
	FindTopByScore(limit int) ([]entity.AudioScore, error)
	FindByAudioID(audioID uint) (*entity.AudioScore, error)
	FindByAudioIDs(audioIDs []uint) ([]entity.AudioScore, error)
	DeleteAll() error
	BulkUpsert(scores []entity.AudioScore) error
}


