package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"mqfm-backend/internal/app/service"
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/shared/dto/request"
	preferencemock "mqfm-backend/tests/mocks/preference"
)

// ──────────────────────── GetOrCreate ────────────────────────

func TestUserPreference_GetOrCreate_ExistingReturned(t *testing.T) {
	existing := &entity.UserPreference{ID: 1, UserID: 5, PlaybackSpeed: 1.5, SleepTimerMinutes: 30}
	repo := &preferencemock.MockUserPreferenceRepository{
		FindByUserFn: func(userID uint) (*entity.UserPreference, error) {
			assert.Equal(t, uint(5), userID)
			return existing, nil
		},
	}
	svc := service.NewUserPreferenceService(repo)
	result, err := svc.GetOrCreate(5)
	assert.NoError(t, err)
	assert.Equal(t, existing, result)
	assert.Equal(t, 1.5, result.PlaybackSpeed)
	assert.Equal(t, 30, result.SleepTimerMinutes)
}

func TestUserPreference_GetOrCreate_CreatesDefault(t *testing.T) {
	var upserted *entity.UserPreference
	repo := &preferencemock.MockUserPreferenceRepository{
		FindByUserFn: func(userID uint) (*entity.UserPreference, error) {
			return nil, errors.New("not found")
		},
		UpsertFn: func(pref *entity.UserPreference) error {
			upserted = pref
			return nil
		},
	}
	svc := service.NewUserPreferenceService(repo)
	result, err := svc.GetOrCreate(5)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(5), result.UserID)
	assert.Equal(t, 1.0, result.PlaybackSpeed, "default playback speed must be 1.0")
	assert.Equal(t, 0, result.SleepTimerMinutes, "default sleep timer must be 0")
	assert.False(t, result.AutoDownloadWifi, "default auto_download_wifi must be false")
	assert.NotNil(t, upserted, "Upsert must be called")
}

func TestUserPreference_GetOrCreate_UpsertError(t *testing.T) {
	repo := &preferencemock.MockUserPreferenceRepository{
		FindByUserFn: func(userID uint) (*entity.UserPreference, error) {
			return nil, errors.New("not found")
		},
		UpsertFn: func(pref *entity.UserPreference) error {
			return errors.New("db error")
		},
	}
	svc := service.NewUserPreferenceService(repo)
	result, err := svc.GetOrCreate(5)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "db error", err.Error())
}

// ──────────────────────── Update ────────────────────────

func TestUserPreference_Update_AllFields(t *testing.T) {
	speed := 2.0
	timer := 15
	wifi := true
	onLike := false

	var upsertedPref *entity.UserPreference
	repo := &preferencemock.MockUserPreferenceRepository{
		FindByUserFn: func(userID uint) (*entity.UserPreference, error) {
			return &entity.UserPreference{ID: 1, UserID: 5, PlaybackSpeed: 1.0}, nil
		},
		UpsertFn: func(pref *entity.UserPreference) error {
			upsertedPref = pref
			return nil
		},
	}
	svc := service.NewUserPreferenceService(repo)

	result, err := svc.Update(5, request.UpdatePreferenceRequest{
		PlaybackSpeed:      &speed,
		SleepTimerMinutes:  &timer,
		AutoDownloadWifi:   &wifi,
		AutoDownloadOnLike: &onLike,
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2.0, result.PlaybackSpeed)
	assert.Equal(t, 15, result.SleepTimerMinutes)
	assert.True(t, result.AutoDownloadWifi)
	assert.False(t, result.AutoDownloadOnLike)
	assert.Equal(t, upsertedPref, result, "returned value must be the upserted one")
}

func TestUserPreference_Update_PartialFields(t *testing.T) {
	speed := 1.75
	repo := &preferencemock.MockUserPreferenceRepository{
		FindByUserFn: func(userID uint) (*entity.UserPreference, error) {
			return &entity.UserPreference{
				ID: 1, UserID: 5, PlaybackSpeed: 1.0,
				SleepTimerMinutes: 10, AutoDownloadWifi: true, AutoDownloadOnLike: true,
			}, nil
		},
		UpsertFn: func(pref *entity.UserPreference) error { return nil },
	}
	svc := service.NewUserPreferenceService(repo)

	result, err := svc.Update(5, request.UpdatePreferenceRequest{PlaybackSpeed: &speed})
	assert.NoError(t, err)
	assert.Equal(t, 1.75, result.PlaybackSpeed, "playback speed should be updated")
	assert.Equal(t, 10, result.SleepTimerMinutes, "unset field should stay unchanged")
	assert.True(t, result.AutoDownloadWifi, "unset field should stay unchanged")
	assert.True(t, result.AutoDownloadOnLike, "unset field should stay unchanged")
}

func TestUserPreference_Update_NoFields(t *testing.T) {
	repo := &preferencemock.MockUserPreferenceRepository{
		FindByUserFn: func(userID uint) (*entity.UserPreference, error) {
			return &entity.UserPreference{ID: 1, UserID: 5, PlaybackSpeed: 1.0}, nil
		},
		UpsertFn: func(pref *entity.UserPreference) error { return nil },
	}
	svc := service.NewUserPreferenceService(repo)

	result, err := svc.Update(5, request.UpdatePreferenceRequest{})
	assert.NoError(t, err)
	assert.Equal(t, 1.0, result.PlaybackSpeed, "unchanged when no fields provided")
}

func TestUserPreference_Update_UpsertError(t *testing.T) {
	speed := 1.5
	repo := &preferencemock.MockUserPreferenceRepository{
		FindByUserFn: func(userID uint) (*entity.UserPreference, error) {
			return &entity.UserPreference{ID: 1, UserID: 5, PlaybackSpeed: 1.0}, nil
		},
		UpsertFn: func(pref *entity.UserPreference) error {
			return errors.New("upsert failed")
		},
	}
	svc := service.NewUserPreferenceService(repo)

	result, err := svc.Update(5, request.UpdatePreferenceRequest{PlaybackSpeed: &speed})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "upsert failed", err.Error())
}
