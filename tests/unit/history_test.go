package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mqfm-backend/internal/app/service"
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/tests/mocks"
)

func TestRecordPlay_Success(t *testing.T) {
	repo := &mocks.MockHistoryRepository{
		UpsertFn: func(history *entity.History) error { return nil },
	}

	svc := service.NewHistoryService(repo)
	err := svc.RecordPlay(1, request.HistoryRequest{AudioID: 10})

	assert.NoError(t, err)
}

func TestGetHistory_Success(t *testing.T) {
	repo := &mocks.MockHistoryRepository{
		FindByUserFn: func(userID uint) ([]entity.History, error) {
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
}

func TestDeleteHistory_Success(t *testing.T) {
	repo := &mocks.MockHistoryRepository{
		DeleteByUserAudioFn: func(userID, audioID uint) error { return nil },
	}

	svc := service.NewHistoryService(repo)
	err := svc.DeleteHistory(1, 10)

	assert.NoError(t, err)
}

func TestClearHistory_Success(t *testing.T) {
	repo := &mocks.MockHistoryRepository{
		DeleteAllFn: func(userID uint) error { return nil },
	}

	svc := service.NewHistoryService(repo)
	err := svc.ClearHistory(1)

	assert.NoError(t, err)
}
