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
	LinkGoogle(userID uint, req request.GoogleLoginRequest) (*entity.User, error)
	UnlinkGoogle(userID uint) (*entity.User, error)
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
	Create(req request.CreateAudioRequest, audioFile *multipart.FileHeader, thumbnailFile *multipart.FileHeader) (*entity.Audio, error)
	FindAll() ([]entity.Audio, error)
	FindByID(id uint) (*entity.Audio, error)
	Update(id uint, req request.UpdateAudioRequest, audioFile *multipart.FileHeader, thumbnailFile *multipart.FileHeader) (*entity.Audio, error)
	Delete(id uint) error
	Search(query string) ([]entity.Audio, error)
}

type PlaylistService interface {
	Create(playlist *entity.Playlist, file interface{}) error
	GetByUserID(userID uint) ([]entity.Playlist, error)
	GetByID(id uint) (*entity.Playlist, error)
	Update(id uint, updates map[string]interface{}) (*entity.Playlist, error)
	Delete(id, userID uint) error
	Search(query string) ([]entity.Playlist, error)
	AddAudioToPlaylist(playlistID, audioID uint) error
	RemoveAudioFromPlaylist(playlistID, audioID uint) error
	SharePlaylist(playlistID uint) (string, error)
	GetByShareToken(token string) (*entity.Playlist, error)
}

type LikeService interface {
	Like(userID uint, req request.LikeRequest) (*entity.Like, error)
	Unlike(userID uint, req request.UnlikeRequest) error
	GetLikes(userID string, targetType string) ([]entity.Like, error)
	CountByTarget(targetType string, targetID uint) (int64, error)
}

type HistoryService interface {
	RecordPlay(userID uint, req request.HistoryRequest) error
	GetHistory(userID uint) ([]entity.History, error)
	DeleteHistory(userID, audioID uint) error
	ClearHistory(userID uint) error
}

type OTPService interface {
	SendOTP(email string) error
	VerifyOTP(email string, code string) (*entity.User, error)
}

type EmailService interface {
	SendAsync(to string, subject string, body string)
}

type RecommendationService interface {
	GetPopular(limit int) ([]entity.Audio, error)
	GetMostListened(limit int) ([]entity.Audio, error)
	GetByArtist(artist string, limit int) ([]entity.Audio, error)
	GetSimilar(audioID uint, limit int) ([]entity.Audio, error)
	GetQuickPick(userID uint, limit int) ([]entity.Audio, error)
	GetOnboarding(limit int) ([]entity.Audio, error)
	GetPersonalized(userID uint, limit int) ([]entity.Audio, error)
	GetLocationBased(userID uint, limit int) ([]entity.Audio, error)
	GetTimeBasedPersonalized(userID uint, hour int, limit int) ([]entity.Audio, error)
	RecalculateScores() error
}

type ColorExtractorService interface {
	ExtractDominantColor(imagePath string) (string, error)
}

type BookmarkService interface {
	Create(userID uint, req request.CreateBookmarkRequest) (*entity.Bookmark, error)
	GetByUser(userID uint) ([]entity.Bookmark, error)
	GetByUserAndAudio(userID, audioID uint) ([]entity.Bookmark, error)
	Delete(id, userID uint) error
}

type NotificationService interface {
	GetByUser(userID uint, page, limit int) ([]entity.Notification, error)
	MarkAsRead(id, userID uint) error
	MarkAllAsRead(userID uint) error
	CountUnread(userID uint) (int64, error)
	GetSetting(userID uint) (*entity.NotificationSetting, error)
	UpdateSetting(userID uint, req request.UpdateNotificationSettingRequest) (*entity.NotificationSetting, error)
	NotifyNewAudio(audio *entity.Audio) error
	NotifyDailyReminder() error
	NotifyEvent(event *entity.Event) error
}

type AudioProgressService interface {
	UpdateProgress(userID uint, req request.UpdateProgressRequest) (*entity.AudioProgress, error)
	GetProgress(userID, audioID uint) (*entity.AudioProgress, error)
	GetAllProgress(userID uint) ([]entity.AudioProgress, error)
	GetCompleted(userID uint) ([]entity.AudioProgress, error)
}

type DownloadService interface {
	RecordDownload(userID uint, req request.DownloadRequest) (*entity.Download, error)
	GetDownloads(userID uint) ([]entity.Download, error)
	DeleteDownload(id, userID uint) error
	GetStorageUsage(userID uint) (int64, error)
	GetNewFromFavorites(userID uint) ([]entity.Audio, error)
	CleanupExpired() (int64, error)
}

type ListeningStatService interface {
	RecordStat(userID uint, req request.RecordStatRequest) error
	GetWeeklySummary(userID uint) (int, error)
	GetMonthlySummary(userID uint) (int, error)
	GetTopCategories(userID uint, limit int) ([]CategoryStat, error)
	GetTopArtists(userID uint, limit int) ([]ArtistStat, error)
	GetDailySummary(userID uint, days int) ([]DailyStat, error)
	GetRecap(userID uint) (*ListeningRecap, error)
}

type ListeningRecap struct {
	WeeklyMinutes  int
	MonthlyMinutes int
	TopCategories  []CategoryStat
	TopArtists     []ArtistStat
	DailyStats     []DailyStat
}

type AudioClipService interface {
	CreateClip(userID uint, req request.CreateClipRequest) (*entity.AudioClip, error)
	GetByUser(userID uint) ([]entity.AudioClip, error)
	GetByShareToken(token string) (*entity.AudioClip, error)
	Delete(id, userID uint) error
}

type EventService interface {
	Create(req request.CreateEventRequest, file *multipart.FileHeader) (*entity.Event, error)
	FindAll() ([]entity.Event, error)
	FindByID(id uint) (*entity.Event, error)
	Update(id uint, req request.UpdateEventRequest, file *multipart.FileHeader) (*entity.Event, error)
	Delete(id uint) error
	GetUpcoming(limit int) ([]entity.Event, error)
	RSVP(userID, eventID uint) error
	CancelRSVP(userID, eventID uint) error
	GetUserRSVPs(userID uint) ([]entity.EventRSVP, error)
	GetRSVPCount(eventID uint) (int64, error)
}

type UserPreferenceService interface {
	GetOrCreate(userID uint) (*entity.UserPreference, error)
	Update(userID uint, req request.UpdatePreferenceRequest) (*entity.UserPreference, error)
}

type AudioSeriesService interface {
	Create(req request.CreateSeriesRequest, file *multipart.FileHeader) (*entity.AudioSeries, error)
	FindAll() ([]entity.AudioSeries, error)
	FindByID(id uint) (*entity.AudioSeries, error)
	Update(id uint, req request.UpdateSeriesRequest, file *multipart.FileHeader) (*entity.AudioSeries, error)
	Delete(id uint) error
	AddItem(req request.AddSeriesItemRequest) error
	RemoveItem(seriesID, audioID uint) error
	GetItems(seriesID uint) ([]entity.AudioSeriesItem, error)
	Search(query string) ([]entity.AudioSeries, error)
}

type AudioVoteService interface {
	Vote(userID, audioID uint) error
	Unvote(userID, audioID uint) error
	GetWeeklyRanking(limit int) ([]entity.AudioRanking, error)
	GetMonthlyRanking(limit int) ([]entity.AudioRanking, error)
	HasVoted(userID, audioID uint) (bool, error)
	RecalculateRankings() error
	ResetWeeklyVotes() error
	ResetMonthlyVotes() error
}

type SmartResumeService interface {
	Update(userID uint, req request.UpdateResumeRequest) (*entity.SmartResume, error)
	Get(userID uint) (*entity.SmartResume, error)
}

type ShareService interface {
	GenerateAudioShareLink(audioID uint) string
	GenerateClipShareLink(shareToken string) string
	GeneratePlaylistShareLink(shareToken string) string
}

type FavoriteArtistService interface {
	Add(userID uint, artistName string) error
	Remove(userID uint, artistName string) error
	GetByUser(userID uint) ([]entity.FavoriteArtist, error)
}

type UserLocationService interface {
	Update(userID uint, req request.UpdateLocationRequest) (*entity.UserLocation, error)
	Get(userID uint) (*entity.UserLocation, error)
}

type PlaylistCollabService interface {
	AddCollaborator(ownerID, playlistID, collaboratorID uint) error
	RemoveCollaborator(ownerID, playlistID, collaboratorID uint) error
	GetCollaborators(playlistID uint) ([]entity.PlaylistCollaborator, error)
	ContributeAudio(userID, playlistID, audioID uint) error
	JoinByShareToken(userID uint, token string) error
}

type AudioConverterService interface {
	ConvertToOGG(inputPath string) (string, error)
	CreateClip(inputPath string, startSec, endSec int) (string, error)
}


