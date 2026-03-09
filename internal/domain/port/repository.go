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
	FindBySeriesID(seriesID uint) ([]entity.Audio, error)
	FindNewByArtists(artists []string, since time.Time) ([]entity.Audio, error)
	CountAll() (int64, error)
}

type PlaylistRepository interface {
	Create(playlist *entity.Playlist) error
	FindByUserID(userID uint) ([]entity.Playlist, error)
	FindByID(id uint, userID uint) (*entity.Playlist, error)
	Search(userID uint, query string) ([]entity.Playlist, error)
	AddAudio(playlist *entity.Playlist, audio *entity.Audio) error
	FindAudioByID(id uint) (*entity.Audio, error)
	RemoveAudio(playlist *entity.Playlist, audio *entity.Audio) error
	FindByShareToken(token string) (*entity.Playlist, error)
	Update(id uint, updates map[string]interface{}) error
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

type BookmarkRepository interface {
	Create(bookmark *entity.Bookmark) error
	FindByUser(userID uint) ([]entity.Bookmark, error)
	FindByUserAndAudio(userID, audioID uint) ([]entity.Bookmark, error)
	Delete(id, userID uint) error
	DeleteAllByUserAndAudio(userID, audioID uint) error
}

type NotificationRepository interface {
	Create(notification *entity.Notification) error
	BulkCreate(notifications []entity.Notification) error
	FindByUser(userID uint, limit, offset int) ([]entity.Notification, error)
	MarkAsRead(id, userID uint) error
	MarkAllAsRead(userID uint) error
	CountUnread(userID uint) (int64, error)
	GetSetting(userID uint) (*entity.NotificationSetting, error)
	UpsertSetting(setting *entity.NotificationSetting) error
	FindUserIDsWithSetting(field string, value bool) ([]uint, error)
}

type AudioProgressRepository interface {
	Upsert(progress *entity.AudioProgress) error
	FindByUserAndAudio(userID, audioID uint) (*entity.AudioProgress, error)
	FindByUser(userID uint) ([]entity.AudioProgress, error)
	FindCompletedByUser(userID uint) ([]entity.AudioProgress, error)
	Delete(userID, audioID uint) error
}

type DownloadRepository interface {
	Create(download *entity.Download) error
	FindByUser(userID uint) ([]entity.Download, error)
	Delete(id, userID uint) error
	Exists(userID, audioID uint) (bool, error)
	SumSizeByUser(userID uint) (int64, error)
}

type ListeningStatRepository interface {
	Create(stat *entity.ListeningStat) error
	GetWeeklySummary(userID uint) (int, error)
	GetMonthlySummary(userID uint) (int, error)
	GetTopCategories(userID uint, limit int) ([]CategoryStat, error)
	GetTopArtists(userID uint, limit int) ([]ArtistStat, error)
	GetDailySummary(userID uint, days int) ([]DailyStat, error)
}

type CategoryStat struct {
	CategoryID uint
	Name       string
	TotalTime  int
}

type ArtistStat struct {
	Artist    string
	TotalTime int
}

type DailyStat struct {
	Date      string
	TotalTime int
}

type AudioClipRepository interface {
	Create(clip *entity.AudioClip) error
	FindByUser(userID uint) ([]entity.AudioClip, error)
	FindByShareToken(token string) (*entity.AudioClip, error)
	Delete(id, userID uint) error
}

type EventRepository interface {
	Create(event *entity.Event) error
	FindAll() ([]entity.Event, error)
	FindByID(id uint) (*entity.Event, error)
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	FindUpcoming(limit int) ([]entity.Event, error)
	CreateRSVP(rsvp *entity.EventRSVP) error
	DeleteRSVP(userID, eventID uint) error
	FindRSVPsByEvent(eventID uint) ([]entity.EventRSVP, error)
	FindRSVPsByUser(userID uint) ([]entity.EventRSVP, error)
	ExistsRSVP(userID, eventID uint) (bool, error)
	CountRSVP(eventID uint) (int64, error)
}

type UserPreferenceRepository interface {
	Upsert(pref *entity.UserPreference) error
	FindByUser(userID uint) (*entity.UserPreference, error)
}

type AudioSeriesRepository interface {
	Create(series *entity.AudioSeries) error
	FindAll() ([]entity.AudioSeries, error)
	FindByID(id uint) (*entity.AudioSeries, error)
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	AddItem(item *entity.AudioSeriesItem) error
	RemoveItem(seriesID, audioID uint) error
	FindItems(seriesID uint) ([]entity.AudioSeriesItem, error)
	Search(query string) ([]entity.AudioSeries, error)
}

type AudioVoteRepository interface {
	Create(vote *entity.AudioVote) error
	Delete(userID, audioID uint) error
	Exists(userID, audioID uint) (bool, error)
	CountByAudio(audioID uint) (int64, error)
	CountWeeklyByAudio(audioID uint) (int64, error)
	CountMonthlyByAudio(audioID uint) (int64, error)
	FindVotedAudioIDs(userID uint) ([]uint, error)
}

type AudioRankingRepository interface {
	Upsert(ranking *entity.AudioRanking) error
	BulkUpsert(rankings []entity.AudioRanking) error
	FindTopWeekly(limit int) ([]entity.AudioRanking, error)
	FindTopMonthly(limit int) ([]entity.AudioRanking, error)
	FindByAudioID(audioID uint) (*entity.AudioRanking, error)
	CountAll() (int64, error)
	FindAll(limit, offset int) ([]entity.AudioRanking, error)
}

type FavoriteArtistRepository interface {
	Create(fav *entity.FavoriteArtist) error
	Delete(userID uint, artistName string) error
	FindByUser(userID uint) ([]entity.FavoriteArtist, error)
	Exists(userID uint, artistName string) (bool, error)
}

type SmartResumeRepository interface {
	Upsert(resume *entity.SmartResume) error
	FindByUser(userID uint) (*entity.SmartResume, error)
}

type UserLocationRepository interface {
	Upsert(location *entity.UserLocation) error
	FindByUser(userID uint) (*entity.UserLocation, error)
}

type PlaylistCollaboratorRepository interface {
	Create(collab *entity.PlaylistCollaborator) error
	Delete(playlistID, userID uint) error
	FindByPlaylist(playlistID uint) ([]entity.PlaylistCollaborator, error)
	Exists(playlistID, userID uint) (bool, error)
	IsOwnerOrCollaborator(playlistID, userID uint) (bool, error)
}


