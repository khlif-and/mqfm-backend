package mocks

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"time"
)

type MockAdminRepository struct {
	CreateFn      func(admin *entity.Admin) error
	FindByEmailFn func(email string) (*entity.Admin, error)
	FindByIDFn    func(id uint) (*entity.Admin, error)
	UpdateFn      func(id uint, updates map[string]interface{}) error
}

func (m *MockAdminRepository) Create(admin *entity.Admin) error      { return m.CreateFn(admin) }
func (m *MockAdminRepository) FindByEmail(email string) (*entity.Admin, error) {
	return m.FindByEmailFn(email)
}
func (m *MockAdminRepository) FindByID(id uint) (*entity.Admin, error) { return m.FindByIDFn(id) }
func (m *MockAdminRepository) Update(id uint, updates map[string]interface{}) error {
	return m.UpdateFn(id, updates)
}

type MockUserRepository struct {
	CreateFn         func(user *entity.User) error
	FindByEmailFn    func(email string) (*entity.User, error)
	FindByIDFn       func(id uint) (*entity.User, error)
	FindByProviderFn func(provider string, providerID string) (*entity.User, error)
	UpdateFn         func(id uint, updates map[string]interface{}) error
}

func (m *MockUserRepository) Create(user *entity.User) error { return m.CreateFn(user) }
func (m *MockUserRepository) FindByEmail(email string) (*entity.User, error) {
	return m.FindByEmailFn(email)
}
func (m *MockUserRepository) FindByID(id uint) (*entity.User, error) { return m.FindByIDFn(id) }
func (m *MockUserRepository) FindByProviderID(provider string, providerID string) (*entity.User, error) {
	return m.FindByProviderFn(provider, providerID)
}
func (m *MockUserRepository) Update(id uint, updates map[string]interface{}) error {
	return m.UpdateFn(id, updates)
}

type MockCategoryRepository struct {
	FindAllFn func() ([]entity.Category, error)
	FindByIDFn func(id uint) (*entity.Category, error)
	CreateFn   func(category *entity.Category) error
	UpdateFn   func(id uint, updates map[string]interface{}) error
	DeleteFn   func(id uint) error
	SearchFn   func(query string) ([]entity.Category, error)
}

func (m *MockCategoryRepository) FindAll() ([]entity.Category, error) { return m.FindAllFn() }
func (m *MockCategoryRepository) FindByID(id uint) (*entity.Category, error) {
	return m.FindByIDFn(id)
}
func (m *MockCategoryRepository) Create(category *entity.Category) error {
	return m.CreateFn(category)
}
func (m *MockCategoryRepository) Update(id uint, updates map[string]interface{}) error {
	return m.UpdateFn(id, updates)
}
func (m *MockCategoryRepository) Delete(id uint) error            { return m.DeleteFn(id) }
func (m *MockCategoryRepository) Search(query string) ([]entity.Category, error) {
	return m.SearchFn(query)
}

type MockAudioRepository struct {
	FindAllFn        func() ([]entity.Audio, error)
	FindByIDFn       func(id uint) (*entity.Audio, error)
	CreateFn         func(audio *entity.Audio) error
	UpdateFn         func(id uint, updates map[string]interface{}) error
	DeleteFn         func(id uint) error
	SearchFn         func(query string) ([]entity.Audio, error)
	FindByIDsFn      func(ids []uint) ([]entity.Audio, error)
	FindByArtistFn   func(artist string, limit int) ([]entity.Audio, error)
	FindByCategoryFn func(categoryID uint, limit int) ([]entity.Audio, error)
	FindAllActiveFn  func() ([]entity.Audio, error)
	FindBySeriesFn   func(seriesID uint) ([]entity.Audio, error)
	FindNewByArtistsFn func(artists []string, since time.Time) ([]entity.Audio, error)
	CountAllFn       func() (int64, error)
}

func (m *MockAudioRepository) FindAll() ([]entity.Audio, error)   { return m.FindAllFn() }
func (m *MockAudioRepository) FindByID(id uint) (*entity.Audio, error) { return m.FindByIDFn(id) }
func (m *MockAudioRepository) Create(audio *entity.Audio) error    { return m.CreateFn(audio) }
func (m *MockAudioRepository) Update(id uint, updates map[string]interface{}) error {
	return m.UpdateFn(id, updates)
}
func (m *MockAudioRepository) Delete(id uint) error               { return m.DeleteFn(id) }
func (m *MockAudioRepository) Search(query string) ([]entity.Audio, error) {
	return m.SearchFn(query)
}
func (m *MockAudioRepository) FindByIDs(ids []uint) ([]entity.Audio, error) {
	return m.FindByIDsFn(ids)
}
func (m *MockAudioRepository) FindByArtist(artist string, limit int) ([]entity.Audio, error) {
	return m.FindByArtistFn(artist, limit)
}
func (m *MockAudioRepository) FindByCategoryID(categoryID uint, limit int) ([]entity.Audio, error) {
	return m.FindByCategoryFn(categoryID, limit)
}
func (m *MockAudioRepository) FindAllActive() ([]entity.Audio, error) { return m.FindAllActiveFn() }
func (m *MockAudioRepository) FindBySeriesID(seriesID uint) ([]entity.Audio, error) {
	return m.FindBySeriesFn(seriesID)
}
func (m *MockAudioRepository) FindNewByArtists(artists []string, since time.Time) ([]entity.Audio, error) {
	return m.FindNewByArtistsFn(artists, since)
}
func (m *MockAudioRepository) CountAll() (int64, error) { return m.CountAllFn() }

type MockLikeRepository struct {
	CreateFn        func(like *entity.Like) error
	DeleteFn        func(userID uint, targetType string, targetID uint) error
	FindByUserFn    func(userID uint, targetType string) ([]entity.Like, error)
	ExistsFn        func(userID uint, targetType string, targetID uint) (bool, error)
	CountByTargetFn func(targetType string, targetID uint) (int64, error)
	AggregateFn     func() (map[uint]int64, error)
}

func (m *MockLikeRepository) Create(like *entity.Like) error { return m.CreateFn(like) }
func (m *MockLikeRepository) Delete(userID uint, targetType string, targetID uint) error {
	return m.DeleteFn(userID, targetType, targetID)
}
func (m *MockLikeRepository) FindByUser(userID uint, targetType string) ([]entity.Like, error) {
	return m.FindByUserFn(userID, targetType)
}
func (m *MockLikeRepository) Exists(userID uint, targetType string, targetID uint) (bool, error) {
	return m.ExistsFn(userID, targetType, targetID)
}
func (m *MockLikeRepository) CountByTarget(targetType string, targetID uint) (int64, error) {
	return m.CountByTargetFn(targetType, targetID)
}
func (m *MockLikeRepository) AggregateLikeCounts() (map[uint]int64, error) { return m.AggregateFn() }
func (m *MockLikeRepository) AggregateWeeklyLikeCounts(since time.Time) (map[uint]int64, error) {
	return map[uint]int64{}, nil
}

type MockHistoryRepository struct {
	UpsertFn            func(history *entity.History) error
	FindByUserFn        func(userID uint) ([]entity.History, error)
	FindByUsersFn       func(userIDs []uint) ([]entity.History, error)
	DeleteByUserAudioFn func(userID, audioID uint) error
	DeleteAllFn         func(userID uint) error
	CountByAudioFn      func(audioID uint) (int64, error)
	FindFrequentFn      func(userID uint, minPlays int, limit int) ([]entity.History, error)
	AggregateFn         func() (map[uint]int64, error)
}

func (m *MockHistoryRepository) Upsert(history *entity.History) error { return m.UpsertFn(history) }
func (m *MockHistoryRepository) FindByUser(userID uint) ([]entity.History, error) {
	return m.FindByUserFn(userID)
}
func (m *MockHistoryRepository) FindByUsers(userIDs []uint) ([]entity.History, error) {
	return m.FindByUsersFn(userIDs)
}
func (m *MockHistoryRepository) DeleteByUserAndAudio(userID, audioID uint) error {
	return m.DeleteByUserAudioFn(userID, audioID)
}
func (m *MockHistoryRepository) DeleteAllByUser(userID uint) error { return m.DeleteAllFn(userID) }
func (m *MockHistoryRepository) CountByAudio(audioID uint) (int64, error) {
	return m.CountByAudioFn(audioID)
}
func (m *MockHistoryRepository) FindFrequentByUser(userID uint, minPlays int, limit int) ([]entity.History, error) {
	return m.FindFrequentFn(userID, minPlays, limit)
}
func (m *MockHistoryRepository) AggregatePlayCounts() (map[uint]int64, error) {
	return m.AggregateFn()
}

type MockOTPRepository struct {
	CreateFn           func(otp *entity.OTP) error
	FindLatestFn       func(email string) (*entity.OTP, error)
	MarkVerifiedFn     func(id uint) error
	CountRecentFn      func(email string, since time.Time) (int64, error)
	DeleteExpiredFn    func() error
}

func (m *MockOTPRepository) Create(otp *entity.OTP) error { return m.CreateFn(otp) }
func (m *MockOTPRepository) FindLatestByEmail(email string) (*entity.OTP, error) {
	return m.FindLatestFn(email)
}
func (m *MockOTPRepository) MarkVerified(id uint) error { return m.MarkVerifiedFn(id) }
func (m *MockOTPRepository) CountRecentByEmail(email string, since time.Time) (int64, error) {
	return m.CountRecentFn(email, since)
}
func (m *MockOTPRepository) DeleteExpired() error { return m.DeleteExpiredFn() }

type MockPlaylistRepository struct {
	CreateFn         func(playlist *entity.Playlist) error
	FindByUserIDFn   func(userID uint) ([]entity.Playlist, error)
	FindByIDFn       func(id uint) (*entity.Playlist, error)
	FindAllFn        func() ([]entity.Playlist, error)
	SearchFn         func(query string) ([]entity.Playlist, error)
	AddAudioFn       func(playlist *entity.Playlist, audio *entity.Audio) error
	FindAudioByIDFn  func(id uint) (*entity.Audio, error)
	RemoveAudioFn    func(playlist *entity.Playlist, audio *entity.Audio) error
	FindByShareFn    func(token string) (*entity.Playlist, error)
	UpdateFn         func(id uint, updates map[string]interface{}) error
	DeleteFn         func(id uint) error
	CountAudiosFn    func(playlistID uint) (int, error)
}

func (m *MockPlaylistRepository) Create(playlist *entity.Playlist) error {
	return m.CreateFn(playlist)
}
func (m *MockPlaylistRepository) FindByUserID(userID uint) ([]entity.Playlist, error) {
	return m.FindByUserIDFn(userID)
}
func (m *MockPlaylistRepository) FindByID(id uint) (*entity.Playlist, error) {
	return m.FindByIDFn(id)
}
func (m *MockPlaylistRepository) FindAll() ([]entity.Playlist, error) {
	return m.FindAllFn()
}
func (m *MockPlaylistRepository) Search(query string) ([]entity.Playlist, error) {
	return m.SearchFn(query)
}
func (m *MockPlaylistRepository) AddAudio(playlist *entity.Playlist, audio *entity.Audio) error {
	return m.AddAudioFn(playlist, audio)
}
func (m *MockPlaylistRepository) FindAudioByID(id uint) (*entity.Audio, error) {
	return m.FindAudioByIDFn(id)
}
func (m *MockPlaylistRepository) RemoveAudio(playlist *entity.Playlist, audio *entity.Audio) error {
	return m.RemoveAudioFn(playlist, audio)
}
func (m *MockPlaylistRepository) FindByShareToken(token string) (*entity.Playlist, error) {
	return m.FindByShareFn(token)
}
func (m *MockPlaylistRepository) Update(id uint, updates map[string]interface{}) error {
	return m.UpdateFn(id, updates)
}
func (m *MockPlaylistRepository) Delete(id uint) error {
	return m.DeleteFn(id)
}
func (m *MockPlaylistRepository) CountAudios(playlistID uint) (int, error) {
	return m.CountAudiosFn(playlistID)
}

type MockBookmarkRepository struct {
	CreateFn              func(bookmark *entity.Bookmark) error
	FindByUserFn          func(userID uint) ([]entity.Bookmark, error)
	FindByUserAndAudioFn  func(userID, audioID uint) ([]entity.Bookmark, error)
	DeleteFn              func(id, userID uint) error
	DeleteAllByUserAudioFn func(userID, audioID uint) error
}

func (m *MockBookmarkRepository) Create(bookmark *entity.Bookmark) error {
	return m.CreateFn(bookmark)
}
func (m *MockBookmarkRepository) FindByUser(userID uint) ([]entity.Bookmark, error) {
	return m.FindByUserFn(userID)
}
func (m *MockBookmarkRepository) FindByUserAndAudio(userID, audioID uint) ([]entity.Bookmark, error) {
	return m.FindByUserAndAudioFn(userID, audioID)
}
func (m *MockBookmarkRepository) Delete(id, userID uint) error { return m.DeleteFn(id, userID) }
func (m *MockBookmarkRepository) DeleteAllByUserAndAudio(userID, audioID uint) error {
	return m.DeleteAllByUserAudioFn(userID, audioID)
}

type MockAudioScoreRepository struct {
	UpsertFn        func(score *entity.AudioScore) error
	FindTopFn       func(limit int) ([]entity.AudioScore, error)
	FindTopByLikesFn func(limit int, maxLikes int64) ([]entity.AudioScore, error)
	FindByAudioFn   func(audioID uint) (*entity.AudioScore, error)
	FindByAudiosFn  func(audioIDs []uint) ([]entity.AudioScore, error)
	DeleteAllFn     func() error
	BulkUpsertFn    func(scores []entity.AudioScore) error
}

func (m *MockAudioScoreRepository) Upsert(score *entity.AudioScore) error {
	return m.UpsertFn(score)
}
func (m *MockAudioScoreRepository) FindTopByScore(limit int) ([]entity.AudioScore, error) {
	return m.FindTopFn(limit)
}
func (m *MockAudioScoreRepository) FindTopByLikes(limit int, maxLikes int64) ([]entity.AudioScore, error) {
	return m.FindTopByLikesFn(limit, maxLikes)
}
func (m *MockAudioScoreRepository) FindByAudioID(audioID uint) (*entity.AudioScore, error) {
	return m.FindByAudioFn(audioID)
}
func (m *MockAudioScoreRepository) FindByAudioIDs(audioIDs []uint) ([]entity.AudioScore, error) {
	return m.FindByAudiosFn(audioIDs)
}
func (m *MockAudioScoreRepository) DeleteAll() error { return m.DeleteAllFn() }
func (m *MockAudioScoreRepository) BulkUpsert(scores []entity.AudioScore) error {
	return m.BulkUpsertFn(scores)
}
func (m *MockAudioScoreRepository) FindTopByWeeklyLikes(limit int) ([]entity.AudioScore, error) {
	return []entity.AudioScore{}, nil
}
func (m *MockAudioScoreRepository) BulkUpdateWeeklyLikes(data map[uint]int64) error { return nil }

type MockNotificationRepository struct {
	CreateFn             func(notification *entity.Notification) error
	BulkCreateFn         func(notifications []entity.Notification) error
	FindByUserFn         func(userID uint, limit, offset int) ([]entity.Notification, error)
	MarkAsReadFn         func(id, userID uint) error
	MarkAllAsReadFn      func(userID uint) error
	CountUnreadFn        func(userID uint) (int64, error)
	GetSettingFn         func(userID uint) (*entity.NotificationSetting, error)
	UpsertSettingFn      func(setting *entity.NotificationSetting) error
	FindUserIDsSettingFn func(field string, value bool) ([]uint, error)
}

func (m *MockNotificationRepository) Create(notification *entity.Notification) error {
	return m.CreateFn(notification)
}
func (m *MockNotificationRepository) BulkCreate(notifications []entity.Notification) error {
	return m.BulkCreateFn(notifications)
}
func (m *MockNotificationRepository) FindByUser(userID uint, limit, offset int) ([]entity.Notification, error) {
	return m.FindByUserFn(userID, limit, offset)
}
func (m *MockNotificationRepository) MarkAsRead(id, userID uint) error {
	return m.MarkAsReadFn(id, userID)
}
func (m *MockNotificationRepository) MarkAllAsRead(userID uint) error {
	return m.MarkAllAsReadFn(userID)
}
func (m *MockNotificationRepository) CountUnread(userID uint) (int64, error) {
	return m.CountUnreadFn(userID)
}
func (m *MockNotificationRepository) GetSetting(userID uint) (*entity.NotificationSetting, error) {
	return m.GetSettingFn(userID)
}
func (m *MockNotificationRepository) UpsertSetting(setting *entity.NotificationSetting) error {
	return m.UpsertSettingFn(setting)
}
func (m *MockNotificationRepository) FindUserIDsWithSetting(field string, value bool) ([]uint, error) {
	return m.FindUserIDsSettingFn(field, value)
}

type MockUserLocationRepository struct {
	UpsertFn     func(location *entity.UserLocation) error
	FindByUserFn func(userID uint) (*entity.UserLocation, error)
	FindAllFn    func() ([]entity.UserLocation, error)
}

func (m *MockUserLocationRepository) Upsert(location *entity.UserLocation) error {
	return m.UpsertFn(location)
}
func (m *MockUserLocationRepository) FindByUser(userID uint) (*entity.UserLocation, error) {
	return m.FindByUserFn(userID)
}
func (m *MockUserLocationRepository) FindAll() ([]entity.UserLocation, error) {
	return m.FindAllFn()
}

type MockEventRepository struct {
	CreateFn       func(event *entity.Event) error
	FindAllFn      func() ([]entity.Event, error)
	FindByIDFn     func(id uint) (*entity.Event, error)
	UpdateFn       func(id uint, updates map[string]interface{}) error
	DeleteFn       func(id uint) error
	FindUpcomingFn func(limit int) ([]entity.Event, error)
	CreateRSVPFn   func(rsvp *entity.EventRSVP) error
	DeleteRSVPFn   func(userID, eventID uint) error
	FindRSVPsByEventFn func(eventID uint) ([]entity.EventRSVP, error)
	FindRSVPsByUserFn  func(userID uint) ([]entity.EventRSVP, error)
	ExistsRSVPFn      func(userID, eventID uint) (bool, error)
	CountRSVPFn       func(eventID uint) (int64, error)
}

func (m *MockEventRepository) Create(event *entity.Event) error     { return m.CreateFn(event) }
func (m *MockEventRepository) FindAll() ([]entity.Event, error)     { return m.FindAllFn() }
func (m *MockEventRepository) FindByID(id uint) (*entity.Event, error) { return m.FindByIDFn(id) }
func (m *MockEventRepository) Update(id uint, updates map[string]interface{}) error {
	return m.UpdateFn(id, updates)
}
func (m *MockEventRepository) Delete(id uint) error                  { return m.DeleteFn(id) }
func (m *MockEventRepository) FindUpcoming(limit int) ([]entity.Event, error) {
	return m.FindUpcomingFn(limit)
}
func (m *MockEventRepository) CreateRSVP(rsvp *entity.EventRSVP) error { return m.CreateRSVPFn(rsvp) }
func (m *MockEventRepository) DeleteRSVP(userID, eventID uint) error {
	return m.DeleteRSVPFn(userID, eventID)
}
func (m *MockEventRepository) FindRSVPsByEvent(eventID uint) ([]entity.EventRSVP, error) {
	return m.FindRSVPsByEventFn(eventID)
}
func (m *MockEventRepository) FindRSVPsByUser(userID uint) ([]entity.EventRSVP, error) {
	return m.FindRSVPsByUserFn(userID)
}
func (m *MockEventRepository) ExistsRSVP(userID, eventID uint) (bool, error) {
	return m.ExistsRSVPFn(userID, eventID)
}
func (m *MockEventRepository) CountRSVP(eventID uint) (int64, error) {
	return m.CountRSVPFn(eventID)
}

type MockAudioVoteRepository struct {
	CreateFn            func(vote *entity.AudioVote) error
	DeleteFn            func(userID, audioID uint) error
	ExistsFn            func(userID, audioID uint) (bool, error)
	CountByAudioFn      func(audioID uint) (int64, error)
	CountWeeklyFn       func(audioID uint) (int64, error)
	CountMonthlyFn      func(audioID uint) (int64, error)
	FindVotedAudioIDsFn func(userID uint) ([]uint, error)
}

func (m *MockAudioVoteRepository) Create(vote *entity.AudioVote) error { return m.CreateFn(vote) }
func (m *MockAudioVoteRepository) Delete(userID, audioID uint) error {
	return m.DeleteFn(userID, audioID)
}
func (m *MockAudioVoteRepository) Exists(userID, audioID uint) (bool, error) {
	return m.ExistsFn(userID, audioID)
}
func (m *MockAudioVoteRepository) CountByAudio(audioID uint) (int64, error) {
	return m.CountByAudioFn(audioID)
}
func (m *MockAudioVoteRepository) CountWeeklyByAudio(audioID uint) (int64, error) {
	return m.CountWeeklyFn(audioID)
}
func (m *MockAudioVoteRepository) CountMonthlyByAudio(audioID uint) (int64, error) {
	return m.CountMonthlyFn(audioID)
}
func (m *MockAudioVoteRepository) FindVotedAudioIDs(userID uint) ([]uint, error) {
	return m.FindVotedAudioIDsFn(userID)
}

type MockAudioRankingRepository struct {
	UpsertFn         func(ranking *entity.AudioRanking) error
	BulkUpsertFn     func(rankings []entity.AudioRanking) error
	FindTopWeeklyFn  func(limit int) ([]entity.AudioRanking, error)
	FindTopMonthlyFn func(limit int) ([]entity.AudioRanking, error)
	FindByAudioFn    func(audioID uint) (*entity.AudioRanking, error)
	CountAllFn       func() (int64, error)
	FindAllFn        func(limit, offset int) ([]entity.AudioRanking, error)
}

func (m *MockAudioRankingRepository) Upsert(ranking *entity.AudioRanking) error {
	return m.UpsertFn(ranking)
}
func (m *MockAudioRankingRepository) BulkUpsert(rankings []entity.AudioRanking) error {
	return m.BulkUpsertFn(rankings)
}
func (m *MockAudioRankingRepository) FindTopWeekly(limit int) ([]entity.AudioRanking, error) {
	return m.FindTopWeeklyFn(limit)
}
func (m *MockAudioRankingRepository) FindTopMonthly(limit int) ([]entity.AudioRanking, error) {
	return m.FindTopMonthlyFn(limit)
}
func (m *MockAudioRankingRepository) FindByAudioID(audioID uint) (*entity.AudioRanking, error) {
	return m.FindByAudioFn(audioID)
}
func (m *MockAudioRankingRepository) CountAll() (int64, error) { return m.CountAllFn() }
func (m *MockAudioRankingRepository) FindAll(limit, offset int) ([]entity.AudioRanking, error) {
	return m.FindAllFn(limit, offset)
}

var _ port.AdminRepository = (*MockAdminRepository)(nil)
var _ port.UserRepository = (*MockUserRepository)(nil)
var _ port.CategoryRepository = (*MockCategoryRepository)(nil)
var _ port.AudioRepository = (*MockAudioRepository)(nil)
var _ port.LikeRepository = (*MockLikeRepository)(nil)
var _ port.HistoryRepository = (*MockHistoryRepository)(nil)
var _ port.OTPRepository = (*MockOTPRepository)(nil)
var _ port.PlaylistRepository = (*MockPlaylistRepository)(nil)
var _ port.BookmarkRepository = (*MockBookmarkRepository)(nil)
var _ port.AudioScoreRepository = (*MockAudioScoreRepository)(nil)
var _ port.NotificationRepository = (*MockNotificationRepository)(nil)
var _ port.UserLocationRepository = (*MockUserLocationRepository)(nil)
var _ port.EventRepository = (*MockEventRepository)(nil)
var _ port.AudioVoteRepository = (*MockAudioVoteRepository)(nil)
var _ port.AudioRankingRepository = (*MockAudioRankingRepository)(nil)

type MockDownloadRepository struct {
	CreateFn        func(download *entity.Download) error
	FindByUserFn    func(userID uint) ([]entity.Download, error)
	DeleteFn        func(id, userID uint) error
	ExistsFn        func(userID, audioID uint) (bool, error)
	SumSizeByUserFn func(userID uint) (int64, error)
	DeleteExpiredFn func() (int64, error)
}

func (m *MockDownloadRepository) Create(download *entity.Download) error {
	return m.CreateFn(download)
}
func (m *MockDownloadRepository) FindByUser(userID uint) ([]entity.Download, error) {
	return m.FindByUserFn(userID)
}
func (m *MockDownloadRepository) Delete(id, userID uint) error { return m.DeleteFn(id, userID) }
func (m *MockDownloadRepository) Exists(userID, audioID uint) (bool, error) {
	return m.ExistsFn(userID, audioID)
}
func (m *MockDownloadRepository) SumSizeByUser(userID uint) (int64, error) {
	return m.SumSizeByUserFn(userID)
}
func (m *MockDownloadRepository) DeleteExpired() (int64, error) { return m.DeleteExpiredFn() }

var _ port.DownloadRepository = (*MockDownloadRepository)(nil)
