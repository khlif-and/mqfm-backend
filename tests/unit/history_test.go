package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"mqfm-backend/internal/app/service"
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/shared/dto/request"
	historymock "mqfm-backend/tests/mocks/history"
)

func TestRecordPlay_Success(t *testing.T) {
	var upserted *entity.History
	repo := &historymock.MockHistoryRepository{
		UpsertFn: func(history *entity.History) error {
			upserted = history
			return nil
		},
	}

	svc := service.NewHistoryService(repo)
	err := svc.RecordPlay(1, request.HistoryRequest{AudioID: 10})

	assert.NoError(t, err)
	assert.NotNil(t, upserted)
	assert.Equal(t, uint(1), upserted.UserID, "must pass correct userID")
	assert.Equal(t, uint(10), upserted.AudioID, "must pass correct audioID")
}

func TestRecordPlay_RepoError(t *testing.T) {
	repoErr := errors.New("deadlock")
	repo := &historymock.MockHistoryRepository{
		UpsertFn: func(history *entity.History) error { return repoErr },
	}

	svc := service.NewHistoryService(repo)
	err := svc.RecordPlay(1, request.HistoryRequest{AudioID: 10})

	assert.ErrorIs(t, err, repoErr, "must propagate exact repo error")
}

func TestGetHistory_Success(t *testing.T) {
	repo := &historymock.MockHistoryRepository{
		FindByUserFn: func(userID uint) ([]entity.History, error) {
			assert.Equal(t, uint(1), userID)
			return []entity.History{
				{ID: 1, UserID: userID, AudioID: 10, PlayCount: 5},
				{ID: 2, UserID: userID, AudioID: 20, PlayCount: 3},
			}, nil
		},
	}

	svc := service.NewHistoryService(repo)
	histories, err := svc.GetHistory(1)

	assert.NoError(t, err)
	assert.Len(t, histories, 2)
	assert.Equal(t, 5, histories[0].PlayCount)
	assert.Equal(t, 3, histories[1].PlayCount)
}

func TestGetHistory_Empty(t *testing.T) {
	repo := &historymock.MockHistoryRepository{
		FindByUserFn: func(userID uint) ([]entity.History, error) {
			return []entity.History{}, nil
		},
	}

	svc := service.NewHistoryService(repo)
	histories, err := svc.GetHistory(1)

	assert.NoError(t, err)
	assert.Empty(t, histories)
}

func TestGetHistory_RepoError(t *testing.T) {
	repoErr := errors.New("timeout")
	repo := &historymock.MockHistoryRepository{
		FindByUserFn: func(userID uint) ([]entity.History, error) { return nil, repoErr },
	}

	svc := service.NewHistoryService(repo)
	histories, err := svc.GetHistory(1)

	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, histories)
}

func TestDeleteHistory_Success(t *testing.T) {
	var delUserID, delAudioID uint
	repo := &historymock.MockHistoryRepository{
		DeleteByUserAudioFn: func(userID, audioID uint) error {
			delUserID = userID
			delAudioID = audioID
			return nil
		},
	}

	svc := service.NewHistoryService(repo)
	err := svc.DeleteHistory(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), delUserID)
	assert.Equal(t, uint(10), delAudioID)
}

func TestDeleteHistory_RepoError(t *testing.T) {
	repoErr := errors.New("not found")
	repo := &historymock.MockHistoryRepository{
		DeleteByUserAudioFn: func(userID, audioID uint) error { return repoErr },
	}

	svc := service.NewHistoryService(repo)
	err := svc.DeleteHistory(1, 10)

	assert.ErrorIs(t, err, repoErr)
}

func TestClearHistory_Success(t *testing.T) {
	var clearedUserID uint
	repo := &historymock.MockHistoryRepository{
		DeleteAllFn: func(userID uint) error {
			clearedUserID = userID
			return nil
		},
	}

	svc := service.NewHistoryService(repo)
	err := svc.ClearHistory(1)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), clearedUserID)
}

func TestClearHistory_RepoError(t *testing.T) {
	repoErr := errors.New("permission denied")
	repo := &historymock.MockHistoryRepository{
		DeleteAllFn: func(userID uint) error { return repoErr },
	}

	svc := service.NewHistoryService(repo)
	err := svc.ClearHistory(1)

	assert.ErrorIs(t, err, repoErr)
}
