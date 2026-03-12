package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"mqfm-backend/internal/app/service"
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/shared/dto/request"
	locationmock "mqfm-backend/tests/mocks/location"
)

// ──────────────────────── Update ────────────────────────

func TestUserLocation_Update_Success(t *testing.T) {
	var upserted *entity.UserLocation
	repo := &locationmock.MockUserLocationRepository{
		UpsertFn: func(location *entity.UserLocation) error {
			upserted = location
			assert.Equal(t, uint(3), location.UserID)
			assert.Equal(t, -6.2, location.Latitude)
			assert.Equal(t, 106.8, location.Longitude)
			assert.Equal(t, "Jakarta", location.City)
			return nil
		},
		FindByUserFn: func(userID uint) (*entity.UserLocation, error) {
			assert.Equal(t, uint(3), userID)
			return &entity.UserLocation{
				ID: 1, UserID: 3, Latitude: -6.2, Longitude: 106.8, City: "Jakarta",
			}, nil
		},
	}
	svc := service.NewUserLocationService(repo)

	result, err := svc.Update(3, request.UpdateLocationRequest{
		Latitude: -6.2, Longitude: 106.8, City: "Jakarta",
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, upserted)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Jakarta", result.City)
}

func TestUserLocation_Update_UpsertError(t *testing.T) {
	repo := &locationmock.MockUserLocationRepository{
		UpsertFn: func(location *entity.UserLocation) error {
			return errors.New("upsert failed")
		},
	}
	svc := service.NewUserLocationService(repo)

	result, err := svc.Update(3, request.UpdateLocationRequest{Latitude: 1, Longitude: 2})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "upsert failed", err.Error())
}

func TestUserLocation_Update_FindAfterUpsertError(t *testing.T) {
	repo := &locationmock.MockUserLocationRepository{
		UpsertFn: func(location *entity.UserLocation) error { return nil },
		FindByUserFn: func(userID uint) (*entity.UserLocation, error) {
			return nil, errors.New("find failed")
		},
	}
	svc := service.NewUserLocationService(repo)

	result, err := svc.Update(3, request.UpdateLocationRequest{Latitude: 1, Longitude: 2})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "find failed", err.Error())
}

// ──────────────────────── Get ────────────────────────

func TestUserLocation_Get_Success(t *testing.T) {
	expected := &entity.UserLocation{ID: 1, UserID: 3, Latitude: -6.2, Longitude: 106.8, City: "Bandung"}
	repo := &locationmock.MockUserLocationRepository{
		FindByUserFn: func(userID uint) (*entity.UserLocation, error) {
			assert.Equal(t, uint(3), userID)
			return expected, nil
		},
	}
	svc := service.NewUserLocationService(repo)

	result, err := svc.Get(3)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestUserLocation_Get_NotFound(t *testing.T) {
	repo := &locationmock.MockUserLocationRepository{
		FindByUserFn: func(userID uint) (*entity.UserLocation, error) {
			return nil, errors.New("not found")
		},
	}
	svc := service.NewUserLocationService(repo)

	result, err := svc.Get(99)
	assert.Error(t, err)
	assert.Nil(t, result)
}
