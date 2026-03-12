package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"mqfm-backend/internal/app/service"
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/shared/dto/request"
	bookmarkmock "mqfm-backend/tests/mocks/bookmark"
)

func TestBookmarkCreate_Success(t *testing.T) {
	var saved *entity.Bookmark
	repo := &bookmarkmock.MockBookmarkRepository{
		CreateFn: func(bookmark *entity.Bookmark) error {
			saved = bookmark
			bookmark.ID = 1
			return nil
		},
	}

	svc := service.NewBookmarkService(repo)
	bm, err := svc.Create(1, request.CreateBookmarkRequest{
		AudioID:         10,
		PositionSeconds: 120,
		Label:           "Important part",
	})

	assert.NoError(t, err)
	assert.NotNil(t, bm)
	assert.Equal(t, uint(1), bm.ID)
	assert.Equal(t, uint(1), bm.UserID)
	assert.Equal(t, uint(10), bm.AudioID)
	assert.Equal(t, 120, bm.PositionSeconds)
	assert.Equal(t, "Important part", bm.Label)
	assert.Same(t, saved, bm, "returned object must be the same pointer passed to repo")
}

func TestBookmarkCreate_RepoError(t *testing.T) {
	dbErr := errors.New("db error")
	repo := &bookmarkmock.MockBookmarkRepository{
		CreateFn: func(bookmark *entity.Bookmark) error { return dbErr },
	}

	svc := service.NewBookmarkService(repo)
	bm, err := svc.Create(1, request.CreateBookmarkRequest{AudioID: 10, PositionSeconds: 120})

	assert.ErrorIs(t, err, dbErr, "must propagate exact db error")
	assert.Nil(t, bm)
}

func TestBookmarkGetByUser_Success(t *testing.T) {
	repo := &bookmarkmock.MockBookmarkRepository{
		FindByUserFn: func(userID uint) ([]entity.Bookmark, error) {
			assert.Equal(t, uint(1), userID)
			return []entity.Bookmark{
				{ID: 1, UserID: userID, AudioID: 10, Label: "A"},
				{ID: 2, UserID: userID, AudioID: 20, Label: "B"},
			}, nil
		},
	}

	svc := service.NewBookmarkService(repo)
	bookmarks, err := svc.GetByUser(1)

	assert.NoError(t, err)
	assert.Len(t, bookmarks, 2)
	assert.Equal(t, uint(10), bookmarks[0].AudioID)
	assert.Equal(t, uint(20), bookmarks[1].AudioID)
}

func TestBookmarkGetByUser_RepoError(t *testing.T) {
	repoErr := errors.New("timeout")
	repo := &bookmarkmock.MockBookmarkRepository{
		FindByUserFn: func(userID uint) ([]entity.Bookmark, error) { return nil, repoErr },
	}

	svc := service.NewBookmarkService(repo)
	bookmarks, err := svc.GetByUser(1)

	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, bookmarks)
}

func TestBookmarkGetByUserAndAudio_Success(t *testing.T) {
	repo := &bookmarkmock.MockBookmarkRepository{
		FindByUserAndAudioFn: func(userID, audioID uint) ([]entity.Bookmark, error) {
			assert.Equal(t, uint(1), userID)
			assert.Equal(t, uint(10), audioID)
			return []entity.Bookmark{
				{ID: 1, UserID: userID, AudioID: audioID, PositionSeconds: 60},
			}, nil
		},
	}

	svc := service.NewBookmarkService(repo)
	bookmarks, err := svc.GetByUserAndAudio(1, 10)

	assert.NoError(t, err)
	assert.Len(t, bookmarks, 1)
	assert.Equal(t, 60, bookmarks[0].PositionSeconds)
}

func TestBookmarkGetByUserAndAudio_RepoError(t *testing.T) {
	repoErr := errors.New("db down")
	repo := &bookmarkmock.MockBookmarkRepository{
		FindByUserAndAudioFn: func(userID, audioID uint) ([]entity.Bookmark, error) { return nil, repoErr },
	}

	svc := service.NewBookmarkService(repo)
	bookmarks, err := svc.GetByUserAndAudio(1, 10)

	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, bookmarks)
}

func TestBookmarkDelete_Success(t *testing.T) {
	var deletedID, deletedUserID uint
	repo := &bookmarkmock.MockBookmarkRepository{
		DeleteFn: func(id, userID uint) error {
			deletedID = id
			deletedUserID = userID
			return nil
		},
	}

	svc := service.NewBookmarkService(repo)
	err := svc.Delete(5, 1)

	assert.NoError(t, err)
	assert.Equal(t, uint(5), deletedID, "must pass correct bookmark ID")
	assert.Equal(t, uint(1), deletedUserID, "must pass correct user ID")
}

func TestBookmarkDelete_RepoError(t *testing.T) {
	repoErr := errors.New("not found")
	repo := &bookmarkmock.MockBookmarkRepository{
		DeleteFn: func(id, userID uint) error { return repoErr },
	}

	svc := service.NewBookmarkService(repo)
	err := svc.Delete(1, 1)

	assert.ErrorIs(t, err, repoErr)
}
