package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"mqfm-backend/internal/app/service"
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/shared/dto/request"
	resumemock "mqfm-backend/tests/mocks/resume"
)

// ──────────────────────── Update ────────────────────────

func TestSmartResume_Update_Success(t *testing.T) {
	var upserted *entity.SmartResume
	repo := &resumemock.MockSmartResumeRepository{
		UpsertFn: func(resume *entity.SmartResume) error {
			upserted = resume
			assert.Equal(t, uint(1), resume.UserID)
			assert.Equal(t, uint(10), resume.AudioID)
			assert.Equal(t, uint(5), resume.PlaylistID)
			assert.Equal(t, 120, resume.PositionSeconds)
			return nil
		},
		FindByUserFn: func(userID uint) (*entity.SmartResume, error) {
			return &entity.SmartResume{
				ID: 1, UserID: 1, AudioID: 10, PlaylistID: 5, PositionSeconds: 120,
			}, nil
		},
	}
	svc := service.NewSmartResumeService(repo)

	result, err := svc.Update(1, request.UpdateResumeRequest{
		AudioID: 10, PlaylistID: 5, PositionSeconds: 120,
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, upserted)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, 120, result.PositionSeconds)
}

func TestSmartResume_Update_UpsertError(t *testing.T) {
	repo := &resumemock.MockSmartResumeRepository{
		UpsertFn: func(resume *entity.SmartResume) error {
			return errors.New("upsert failed")
		},
	}
	svc := service.NewSmartResumeService(repo)

	result, err := svc.Update(1, request.UpdateResumeRequest{AudioID: 10, PositionSeconds: 60})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "upsert failed", err.Error())
}

func TestSmartResume_Update_FindAfterUpsertError(t *testing.T) {
	repo := &resumemock.MockSmartResumeRepository{
		UpsertFn: func(resume *entity.SmartResume) error { return nil },
		FindByUserFn: func(userID uint) (*entity.SmartResume, error) {
			return nil, errors.New("find failed")
		},
	}
	svc := service.NewSmartResumeService(repo)

	result, err := svc.Update(1, request.UpdateResumeRequest{AudioID: 10, PositionSeconds: 60})
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ──────────────────────── Get ────────────────────────

func TestSmartResume_Get_Success(t *testing.T) {
	expected := &entity.SmartResume{ID: 1, UserID: 1, AudioID: 10, PositionSeconds: 300}
	repo := &resumemock.MockSmartResumeRepository{
		FindByUserFn: func(userID uint) (*entity.SmartResume, error) {
			assert.Equal(t, uint(1), userID)
			return expected, nil
		},
	}
	svc := service.NewSmartResumeService(repo)

	result, err := svc.Get(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.Equal(t, 300, result.PositionSeconds)
}

func TestSmartResume_Get_NotFound(t *testing.T) {
	repo := &resumemock.MockSmartResumeRepository{
		FindByUserFn: func(userID uint) (*entity.SmartResume, error) {
			return nil, errors.New("not found")
		},
	}
	svc := service.NewSmartResumeService(repo)

	result, err := svc.Get(99)
	assert.Error(t, err)
	assert.Nil(t, result)
}
