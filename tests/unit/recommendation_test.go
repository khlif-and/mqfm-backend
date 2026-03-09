package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mqfm-backend/internal/app/service"
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/tests/mocks"
)

func mkRecSvc(ar *mocks.MockAudioRepository, sr *mocks.MockAudioScoreRepository, hr *mocks.MockHistoryRepository, lr *mocks.MockLikeRepository, loc *mocks.MockUserLocationRepository) interface{ GetPopular(int) ([]entity.Audio, error) } {
	if ar == nil { ar = &mocks.MockAudioRepository{} }
	if sr == nil { sr = &mocks.MockAudioScoreRepository{} }
	if hr == nil { hr = &mocks.MockHistoryRepository{} }
	if lr == nil { lr = &mocks.MockLikeRepository{} }
	if loc == nil { loc = &mocks.MockUserLocationRepository{} }
	return service.NewRecommendationService(ar, sr, hr, lr, loc)
}

func TestRecommendationGetPopular_Success(t *testing.T) {
	audio1 := entity.Audio{ID: 1, Title: "Kajian 1"}
	sr := &mocks.MockAudioScoreRepository{
		FindTopFn: func(limit int) ([]entity.AudioScore, error) {
			return []entity.AudioScore{{AudioID: 1, Audio: &audio1, WeightScore: 10.0}}, nil
		},
	}
	svc := service.NewRecommendationService(&mocks.MockAudioRepository{}, sr, &mocks.MockHistoryRepository{}, &mocks.MockLikeRepository{}, &mocks.MockUserLocationRepository{})
	audios, err := svc.GetPopular(10)
	assert.NoError(t, err)
	assert.Len(t, audios, 1)
	assert.Equal(t, "Kajian 1", audios[0].Title)
}

func TestRecommendationGetMostListened_Success(t *testing.T) {
	audio1 := entity.Audio{ID: 1, Title: "Most Played"}
	sr := &mocks.MockAudioScoreRepository{
		FindTopFn: func(limit int) ([]entity.AudioScore, error) {
			return []entity.AudioScore{{AudioID: 1, Audio: &audio1, TotalPlays: 100}}, nil
		},
	}
	svc := service.NewRecommendationService(&mocks.MockAudioRepository{}, sr, &mocks.MockHistoryRepository{}, &mocks.MockLikeRepository{}, &mocks.MockUserLocationRepository{})
	audios, err := svc.GetMostListened(10)
	assert.NoError(t, err)
	assert.Len(t, audios, 1)
}

func TestRecommendationGetByArtist_Success(t *testing.T) {
	ar := &mocks.MockAudioRepository{
		FindByArtistFn: func(artist string, limit int) ([]entity.Audio, error) {
			return []entity.Audio{{ID: 1, Artist: artist, Title: "Kajian A"}}, nil
		},
	}
	svc := service.NewRecommendationService(ar, &mocks.MockAudioScoreRepository{}, &mocks.MockHistoryRepository{}, &mocks.MockLikeRepository{}, &mocks.MockUserLocationRepository{})
	audios, err := svc.GetByArtist("Ustadz A", 10)
	assert.NoError(t, err)
	assert.Len(t, audios, 1)
}

func TestRecommendationGetQuickPick_Success(t *testing.T) {
	audio := entity.Audio{ID: 1, Title: "Frequent"}
	hr := &mocks.MockHistoryRepository{
		FindFrequentFn: func(userID uint, minPlays int, limit int) ([]entity.History, error) {
			return []entity.History{{ID: 1, UserID: userID, Audio: &audio, PlayCount: 5}}, nil
		},
	}
	svc := service.NewRecommendationService(&mocks.MockAudioRepository{}, &mocks.MockAudioScoreRepository{}, hr, &mocks.MockLikeRepository{}, &mocks.MockUserLocationRepository{})
	audios, err := svc.GetQuickPick(1, 10)
	assert.NoError(t, err)
	assert.Len(t, audios, 1)
}

func TestRecommendationGetSimilar_Success(t *testing.T) {
	target := entity.Audio{ID: 1, Title: "Target", CategoryID: 1, Artist: "A"}
	similar := entity.Audio{ID: 2, Title: "Similar", CategoryID: 1, Artist: "A"}
	ar := &mocks.MockAudioRepository{
		FindByIDFn:      func(id uint) (*entity.Audio, error) { return &target, nil },
		FindAllActiveFn: func() ([]entity.Audio, error) { return []entity.Audio{target, similar}, nil },
	}
	sr := &mocks.MockAudioScoreRepository{
		FindByAudiosFn: func(audioIDs []uint) ([]entity.AudioScore, error) { return []entity.AudioScore{}, nil },
	}
	svc := service.NewRecommendationService(ar, sr, &mocks.MockHistoryRepository{}, &mocks.MockLikeRepository{}, &mocks.MockUserLocationRepository{})
	audios, err := svc.GetSimilar(1, 5)
	assert.NoError(t, err)
	assert.Len(t, audios, 1)
	assert.Equal(t, "Similar", audios[0].Title)
}

func TestRecommendationRecalculateScores_Success(t *testing.T) {
	hr := &mocks.MockHistoryRepository{
		AggregateFn: func() (map[uint]int64, error) { return map[uint]int64{1: 100, 2: 50}, nil },
	}
	lr := &mocks.MockLikeRepository{
		AggregateFn: func() (map[uint]int64, error) { return map[uint]int64{1: 10, 2: 5}, nil },
	}
	sr := &mocks.MockAudioScoreRepository{
		BulkUpsertFn: func(scores []entity.AudioScore) error { return nil },
	}
	svc := service.NewRecommendationService(&mocks.MockAudioRepository{}, sr, hr, lr, &mocks.MockUserLocationRepository{})
	err := svc.RecalculateScores()
	assert.NoError(t, err)
}

func TestRecommendationGetLocationBased_GPSNearby(t *testing.T) {
	loc := &mocks.MockUserLocationRepository{
		FindByUserFn: func(userID uint) (*entity.UserLocation, error) {
			return &entity.UserLocation{UserID: 1, Latitude: -6.2, Longitude: 106.8, City: "Jakarta"}, nil
		},
		FindAllFn: func() ([]entity.UserLocation, error) {
			return []entity.UserLocation{
				{UserID: 1, Latitude: -6.2, Longitude: 106.8},
				{UserID: 2, Latitude: -6.21, Longitude: 106.81},
				{UserID: 3, Latitude: 3.5, Longitude: 98.7},
			}, nil
		},
	}
	hr := &mocks.MockHistoryRepository{
		FindByUsersFn: func(userIDs []uint) ([]entity.History, error) {
			a1 := entity.Audio{ID: 10, Title: "Jakarta Audio"}
			return []entity.History{{UserID: 2, AudioID: 10, Audio: &a1, PlayCount: 5}}, nil
		},
	}
	sr := &mocks.MockAudioScoreRepository{
		FindTopFn: func(limit int) ([]entity.AudioScore, error) { return []entity.AudioScore{}, nil },
	}
	ar := &mocks.MockAudioRepository{
		FindByIDsFn: func(ids []uint) ([]entity.Audio, error) {
			return []entity.Audio{{ID: 10, Title: "Jakarta Audio"}}, nil
		},
	}
	svc := service.NewRecommendationService(ar, sr, hr, &mocks.MockLikeRepository{}, loc)
	audios, err := svc.GetLocationBased(1, 10)
	assert.NoError(t, err)
	assert.NotEmpty(t, audios)
	assert.Equal(t, "Jakarta Audio", audios[0].Title)
}

func TestRecommendationGetLocationBased_NoGPS_FallbackCity(t *testing.T) {
	loc := &mocks.MockUserLocationRepository{
		FindByUserFn: func(userID uint) (*entity.UserLocation, error) {
			return &entity.UserLocation{UserID: 1, Latitude: 0, Longitude: 0, City: "Jakarta"}, nil
		},
		FindAllFn: func() ([]entity.UserLocation, error) { return []entity.UserLocation{}, nil },
	}
	ar := &mocks.MockAudioRepository{
		SearchFn: func(query string) ([]entity.Audio, error) {
			return []entity.Audio{{ID: 1, Title: "Jakarta Kajian"}}, nil
		},
	}
	sr := &mocks.MockAudioScoreRepository{
		FindTopFn: func(limit int) ([]entity.AudioScore, error) { return []entity.AudioScore{}, nil },
	}
	svc := service.NewRecommendationService(ar, sr, &mocks.MockHistoryRepository{}, &mocks.MockLikeRepository{}, loc)
	audios, err := svc.GetLocationBased(1, 10)
	assert.NoError(t, err)
	assert.Len(t, audios, 1)
}

func TestRecommendationGetLocationBased_NoLocation_FallbackPopular(t *testing.T) {
	loc := &mocks.MockUserLocationRepository{
		FindByUserFn: func(userID uint) (*entity.UserLocation, error) {
			return &entity.UserLocation{City: ""}, nil
		},
		FindAllFn: func() ([]entity.UserLocation, error) { return []entity.UserLocation{}, nil },
	}
	sr := &mocks.MockAudioScoreRepository{
		FindTopFn: func(limit int) ([]entity.AudioScore, error) { return []entity.AudioScore{}, nil },
	}
	svc := service.NewRecommendationService(&mocks.MockAudioRepository{}, sr, &mocks.MockHistoryRepository{}, &mocks.MockLikeRepository{}, loc)
	audios, err := svc.GetLocationBased(1, 10)
	assert.NoError(t, err)
	assert.Empty(t, audios)
}

func TestRecommendationGetTimeBasedPersonalized(t *testing.T) {
	sr := &mocks.MockAudioScoreRepository{
		FindTopFn: func(limit int) ([]entity.AudioScore, error) { return []entity.AudioScore{}, nil },
	}
	hr := &mocks.MockHistoryRepository{
		FindFrequentFn: func(userID uint, minPlays int, limit int) ([]entity.History, error) { return []entity.History{}, nil },
		FindByUserFn:   func(userID uint) ([]entity.History, error) { return []entity.History{}, nil },
	}
	ar := &mocks.MockAudioRepository{
		FindAllActiveFn: func() ([]entity.Audio, error) { return []entity.Audio{}, nil },
	}
	lr := &mocks.MockLikeRepository{
		FindByUserFn: func(userID uint, targetType string) ([]entity.Like, error) { return []entity.Like{}, nil },
	}
	svc := service.NewRecommendationService(ar, sr, hr, lr, &mocks.MockUserLocationRepository{})
	for _, h := range []int{8, 14, 19, 23} {
		_, err := svc.GetTimeBasedPersonalized(1, h, 10)
		assert.NoError(t, err)
	}
}
