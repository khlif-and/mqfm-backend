package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"mqfm-backend/internal/app/service"
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/tests/mocks"
)

func TestBookmarkCreate_Success(t *testing.T) {
	repo := &mocks.MockBookmarkRepository{
		CreateFn: func(bookmark *entity.Bookmark) error {
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
	assert.Equal(t, uint(1), bm.UserID)
	assert.Equal(t, uint(10), bm.AudioID)
	assert.Equal(t, 120, bm.PositionSeconds)
	assert.Equal(t, "Important part", bm.Label)
}

func TestBookmarkCreate_Error(t *testing.T) {
	repo := &mocks.MockBookmarkRepository{
		CreateFn: func(bookmark *entity.Bookmark) error {
			return errors.New("db error")
		},
	}

	svc := service.NewBookmarkService(repo)
	bm, err := svc.Create(1, request.CreateBookmarkRequest{AudioID: 10, PositionSeconds: 120})

	assert.Error(t, err)
	assert.Nil(t, bm)
}

func TestBookmarkGetByUser_Success(t *testing.T) {
	repo := &mocks.MockBookmarkRepository{
		FindByUserFn: func(userID uint) ([]entity.Bookmark, error) {
			return []entity.Bookmark{
				{ID: 1, UserID: userID, AudioID: 10},
				{ID: 2, UserID: userID, AudioID: 20},
			}, nil
		},
	}

	svc := service.NewBookmarkService(repo)
	bookmarks, err := svc.GetByUser(1)

	assert.NoError(t, err)
	assert.Len(t, bookmarks, 2)
}

func TestBookmarkGetByUserAndAudio_Success(t *testing.T) {
	repo := &mocks.MockBookmarkRepository{
		FindByUserAndAudioFn: func(userID, audioID uint) ([]entity.Bookmark, error) {
			return []entity.Bookmark{
				{ID: 1, UserID: userID, AudioID: audioID, PositionSeconds: 60},
			}, nil
		},
	}

	svc := service.NewBookmarkService(repo)
	bookmarks, err := svc.GetByUserAndAudio(1, 10)

	assert.NoError(t, err)
	assert.Len(t, bookmarks, 1)
}

func TestBookmarkDelete_Success(t *testing.T) {
	repo := &mocks.MockBookmarkRepository{
		DeleteFn: func(id, userID uint) error { return nil },
	}

	svc := service.NewBookmarkService(repo)
	err := svc.Delete(1, 1)

	assert.NoError(t, err)
}
