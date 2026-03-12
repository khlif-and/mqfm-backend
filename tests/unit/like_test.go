package service_test

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"mqfm-backend/internal/app/service"
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/dto/request"
	downloadmock "mqfm-backend/tests/mocks/download"
	likemock "mqfm-backend/tests/mocks/like"
	preferencemock "mqfm-backend/tests/mocks/preference"
)

// newLikeSvc is a helper that wires all three dependencies for the LikeService.
func newLikeSvc(repo *likemock.MockLikeRepository, dl *downloadmock.MockDownloadService, pref *preferencemock.MockUserPreferenceRepository) interface {
	Like(uint, request.LikeRequest) (*entity.Like, error)
	Unlike(uint, request.UnlikeRequest) error
	GetLikes(string, string) ([]entity.Like, error)
	CountByTarget(string, uint) (int64, error)
} {
	if dl == nil {
		dl = &downloadmock.MockDownloadService{}
	}
	if pref == nil {
		pref = &preferencemock.MockUserPreferenceRepository{}
	}
	return service.NewLikeService(repo, dl, pref)
}

// ──────────────────────── Like (create) ────────────────────────

func TestLike_Success_Audio(t *testing.T) {
	var createdLike *entity.Like
	repo := &likemock.MockLikeRepository{
		ExistsFn: func(userID uint, targetType string, targetID uint) (bool, error) {
			assert.Equal(t, uint(1), userID)
			assert.Equal(t, "audio", targetType)
			assert.Equal(t, uint(10), targetID)
			return false, nil
		},
		CreateFn: func(like *entity.Like) error {
			createdLike = like
			like.ID = 1
			return nil
		},
	}
	// pref returns AutoDownloadOnLike = true so the goroutine fires
	pref := &preferencemock.MockUserPreferenceRepository{
		FindByUserFn: func(userID uint) (*entity.UserPreference, error) {
			return &entity.UserPreference{AutoDownloadOnLike: true}, nil
		},
	}
	var downloadCalled int32
	dl := &downloadmock.MockDownloadService{
		RecordDownloadFn: func(userID uint, req request.DownloadRequest) (*entity.Download, error) {
			atomic.AddInt32(&downloadCalled, 1)
			assert.Equal(t, uint(1), userID)
			assert.Equal(t, uint(10), req.AudioID)
			return &entity.Download{}, nil
		},
	}

	svc := newLikeSvc(repo, dl, pref)
	like, err := svc.Like(1, request.LikeRequest{TargetType: "audio", TargetID: 10})

	assert.NoError(t, err)
	assert.NotNil(t, like)
	assert.Equal(t, uint(1), like.ID)
	assert.Equal(t, uint(1), like.UserID)
	assert.Equal(t, "audio", like.TargetType)
	assert.Equal(t, uint(10), like.TargetID)
	assert.Equal(t, createdLike, like)

	// Wait briefly for the goroutine to trigger auto-download
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, int32(1), atomic.LoadInt32(&downloadCalled), "auto-download should fire for audio like")
}

func TestLike_Success_Playlist_NoAutoDownload(t *testing.T) {
	repo := &likemock.MockLikeRepository{
		ExistsFn: func(userID uint, targetType string, targetID uint) (bool, error) {
			return false, nil
		},
		CreateFn: func(like *entity.Like) error { like.ID = 2; return nil },
	}
	var downloadCalled int32
	dl := &downloadmock.MockDownloadService{
		RecordDownloadFn: func(userID uint, req request.DownloadRequest) (*entity.Download, error) {
			atomic.AddInt32(&downloadCalled, 1)
			return &entity.Download{}, nil
		},
	}

	svc := newLikeSvc(repo, dl, nil)
	like, err := svc.Like(5, request.LikeRequest{TargetType: "playlist", TargetID: 20})

	assert.NoError(t, err)
	assert.NotNil(t, like)
	assert.Equal(t, "playlist", like.TargetType)

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, int32(0), atomic.LoadInt32(&downloadCalled), "auto-download should NOT fire for playlist like")
}

func TestLike_AlreadyLiked_ReturnsError(t *testing.T) {
	repo := &likemock.MockLikeRepository{
		ExistsFn: func(userID uint, targetType string, targetID uint) (bool, error) { return true, nil },
	}
	svc := newLikeSvc(repo, nil, nil)
	like, err := svc.Like(1, request.LikeRequest{TargetType: "audio", TargetID: 10})

	assert.Error(t, err)
	assert.Nil(t, like)
	assert.Equal(t, constant.MsgAlreadyLiked, err.Error(), "error message must be the constant")
}

func TestLike_ExistsCheckFails(t *testing.T) {
	dbErr := errors.New("connection refused")
	repo := &likemock.MockLikeRepository{
		ExistsFn: func(userID uint, targetType string, targetID uint) (bool, error) { return false, dbErr },
	}
	svc := newLikeSvc(repo, nil, nil)
	like, err := svc.Like(1, request.LikeRequest{TargetType: "audio", TargetID: 10})

	assert.ErrorIs(t, err, dbErr, "must propagate the exact db error")
	assert.Nil(t, like)
}

func TestLike_CreateFails(t *testing.T) {
	createErr := errors.New("duplicate key")
	repo := &likemock.MockLikeRepository{
		ExistsFn: func(userID uint, targetType string, targetID uint) (bool, error) { return false, nil },
		CreateFn: func(like *entity.Like) error { return createErr },
	}
	svc := newLikeSvc(repo, nil, nil)
	like, err := svc.Like(1, request.LikeRequest{TargetType: "audio", TargetID: 10})

	assert.ErrorIs(t, err, createErr)
	assert.Nil(t, like)
}

func TestLike_AutoDownload_DisabledInPreference(t *testing.T) {
	repo := &likemock.MockLikeRepository{
		ExistsFn: func(userID uint, targetType string, targetID uint) (bool, error) { return false, nil },
		CreateFn: func(like *entity.Like) error { like.ID = 1; return nil },
	}
	pref := &preferencemock.MockUserPreferenceRepository{
		FindByUserFn: func(userID uint) (*entity.UserPreference, error) {
			return &entity.UserPreference{AutoDownloadOnLike: false}, nil
		},
	}
	var downloadCalled int32
	dl := &downloadmock.MockDownloadService{
		RecordDownloadFn: func(userID uint, req request.DownloadRequest) (*entity.Download, error) {
			atomic.AddInt32(&downloadCalled, 1)
			return &entity.Download{}, nil
		},
	}
	svc := newLikeSvc(repo, dl, pref)
	_, err := svc.Like(1, request.LikeRequest{TargetType: "audio", TargetID: 10})
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, int32(0), atomic.LoadInt32(&downloadCalled), "auto-download must NOT fire when preference is off")
}

func TestLike_AutoDownload_PrefRepoError_NoDownload(t *testing.T) {
	repo := &likemock.MockLikeRepository{
		ExistsFn: func(userID uint, targetType string, targetID uint) (bool, error) { return false, nil },
		CreateFn: func(like *entity.Like) error { like.ID = 1; return nil },
	}
	pref := &preferencemock.MockUserPreferenceRepository{
		FindByUserFn: func(userID uint) (*entity.UserPreference, error) {
			return nil, errors.New("pref not found")
		},
	}
	var downloadCalled int32
	dl := &downloadmock.MockDownloadService{
		RecordDownloadFn: func(userID uint, req request.DownloadRequest) (*entity.Download, error) {
			atomic.AddInt32(&downloadCalled, 1)
			return &entity.Download{}, nil
		},
	}
	svc := newLikeSvc(repo, dl, pref)
	_, err := svc.Like(1, request.LikeRequest{TargetType: "audio", TargetID: 10})
	assert.NoError(t, err) // like itself succeeds

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, int32(0), atomic.LoadInt32(&downloadCalled), "auto-download must NOT fire when pref repo errors")
}

// ──────────────────────── Unlike ────────────────────────

func TestUnlike_Success(t *testing.T) {
	var deletedUser uint
	var deletedType string
	var deletedID uint
	repo := &likemock.MockLikeRepository{
		DeleteFn: func(userID uint, targetType string, targetID uint) error {
			deletedUser = userID
			deletedType = targetType
			deletedID = targetID
			return nil
		},
	}
	svc := newLikeSvc(repo, nil, nil)
	err := svc.Unlike(1, request.UnlikeRequest{TargetType: "audio", TargetID: 10})

	assert.NoError(t, err)
	assert.Equal(t, uint(1), deletedUser)
	assert.Equal(t, "audio", deletedType)
	assert.Equal(t, uint(10), deletedID)
}

func TestUnlike_RepoError(t *testing.T) {
	repoErr := errors.New("not found")
	repo := &likemock.MockLikeRepository{
		DeleteFn: func(userID uint, targetType string, targetID uint) error { return repoErr },
	}
	svc := newLikeSvc(repo, nil, nil)
	err := svc.Unlike(1, request.UnlikeRequest{TargetType: "audio", TargetID: 10})

	assert.ErrorIs(t, err, repoErr)
}

// ──────────────────────── GetLikes ────────────────────────

func TestGetLikes_Success(t *testing.T) {
	repo := &likemock.MockLikeRepository{
		FindByUserFn: func(userID uint, targetType string) ([]entity.Like, error) {
			assert.Equal(t, uint(1), userID, "userID string '1' must parse to uint 1")
			assert.Equal(t, "audio", targetType)
			return []entity.Like{
				{ID: 1, UserID: userID, TargetType: targetType, TargetID: 10},
				{ID: 2, UserID: userID, TargetType: targetType, TargetID: 20},
			}, nil
		},
	}
	svc := newLikeSvc(repo, nil, nil)
	likes, err := svc.GetLikes("1", "audio")

	assert.NoError(t, err)
	assert.Len(t, likes, 2)
	assert.Equal(t, uint(10), likes[0].TargetID)
	assert.Equal(t, uint(20), likes[1].TargetID)
}

func TestGetLikes_Empty(t *testing.T) {
	repo := &likemock.MockLikeRepository{
		FindByUserFn: func(userID uint, targetType string) ([]entity.Like, error) {
			return []entity.Like{}, nil
		},
	}
	svc := newLikeSvc(repo, nil, nil)
	likes, err := svc.GetLikes("1", "audio")

	assert.NoError(t, err)
	assert.Empty(t, likes)
}

func TestGetLikes_RepoError(t *testing.T) {
	repoErr := errors.New("timeout")
	repo := &likemock.MockLikeRepository{
		FindByUserFn: func(userID uint, targetType string) ([]entity.Like, error) {
			return nil, repoErr
		},
	}
	svc := newLikeSvc(repo, nil, nil)
	likes, err := svc.GetLikes("1", "audio")

	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, likes)
}

// ──────────────────────── CountByTarget ────────────────────────

func TestCountByTarget_Success(t *testing.T) {
	repo := &likemock.MockLikeRepository{
		CountByTargetFn: func(targetType string, targetID uint) (int64, error) {
			assert.Equal(t, "audio", targetType)
			assert.Equal(t, uint(10), targetID)
			return 42, nil
		},
	}
	svc := newLikeSvc(repo, nil, nil)
	count, err := svc.CountByTarget("audio", 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(42), count)
}

func TestCountByTarget_Zero(t *testing.T) {
	repo := &likemock.MockLikeRepository{
		CountByTargetFn: func(targetType string, targetID uint) (int64, error) { return 0, nil },
	}
	svc := newLikeSvc(repo, nil, nil)
	count, err := svc.CountByTarget("audio", 999)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestCountByTarget_RepoError(t *testing.T) {
	repoErr := errors.New("db down")
	repo := &likemock.MockLikeRepository{
		CountByTargetFn: func(targetType string, targetID uint) (int64, error) { return 0, repoErr },
	}
	svc := newLikeSvc(repo, nil, nil)
	count, err := svc.CountByTarget("audio", 10)

	assert.ErrorIs(t, err, repoErr)
	assert.Equal(t, int64(0), count)
}
