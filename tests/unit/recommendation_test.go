package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"mqfm-backend/internal/app/service"
	"mqfm-backend/internal/domain/entity"
	audiomock "mqfm-backend/tests/mocks/audio"
	historymock "mqfm-backend/tests/mocks/history"
	likemock "mqfm-backend/tests/mocks/like"
	locationmock "mqfm-backend/tests/mocks/location"
)

func mkRecSvc(ar *audiomock.MockAudioRepository, sr *audiomock.MockAudioScoreRepository, hr *historymock.MockHistoryRepository, lr *likemock.MockLikeRepository, loc *locationmock.MockUserLocationRepository) interface{ GetPopular(int) ([]entity.Audio, error) } {
	if ar == nil { ar = &audiomock.MockAudioRepository{} }
	if sr == nil { sr = &audiomock.MockAudioScoreRepository{} }
	if hr == nil { hr = &historymock.MockHistoryRepository{} }
	if lr == nil { lr = &likemock.MockLikeRepository{} }
	if loc == nil { loc = &locationmock.MockUserLocationRepository{} }
	return service.NewRecommendationService(ar, sr, hr, lr, loc)
}

func TestRecommendationGetPopular_Success(t *testing.T) {
	audio1 := entity.Audio{ID: 1, Title: "Kajian 1"}
	sr := &audiomock.MockAudioScoreRepository{
		FindTopFn: func(limit int) ([]entity.AudioScore, error) {
			return []entity.AudioScore{{AudioID: 1, Audio: &audio1, WeightScore: 10.0}}, nil
		},
	}
	svc := service.NewRecommendationService(&audiomock.MockAudioRepository{}, sr, &historymock.MockHistoryRepository{}, &likemock.MockLikeRepository{}, &locationmock.MockUserLocationRepository{})
	audios, err := svc.GetPopular(10)
	assert.NoError(t, err)
	assert.Len(t, audios, 1)
	assert.Equal(t, "Kajian 1", audios[0].Title)
}

func TestRecommendationGetMostListened_Success(t *testing.T) {
	audio1 := entity.Audio{ID: 1, Title: "Most Played"}
	sr := &audiomock.MockAudioScoreRepository{
		FindTopFn: func(limit int) ([]entity.AudioScore, error) {
			return []entity.AudioScore{{AudioID: 1, Audio: &audio1, TotalPlays: 100}}, nil
		},
	}
	svc := service.NewRecommendationService(&audiomock.MockAudioRepository{}, sr, &historymock.MockHistoryRepository{}, &likemock.MockLikeRepository{}, &locationmock.MockUserLocationRepository{})
	audios, err := svc.GetMostListened(10)
	assert.NoError(t, err)
	assert.Len(t, audios, 1)
}

func TestRecommendationGetByArtist_Success(t *testing.T) {
	ar := &audiomock.MockAudioRepository{
		FindByArtistFn: func(artist string, limit int) ([]entity.Audio, error) {
			return []entity.Audio{{ID: 1, Artist: artist, Title: "Kajian A"}}, nil
		},
	}
	svc := service.NewRecommendationService(ar, &audiomock.MockAudioScoreRepository{}, &historymock.MockHistoryRepository{}, &likemock.MockLikeRepository{}, &locationmock.MockUserLocationRepository{})
	audios, err := svc.GetByArtist("Ustadz A", 10)
	assert.NoError(t, err)
	assert.Len(t, audios, 1)
}

func TestRecommendationGetQuickPick_Success(t *testing.T) {
	audio := entity.Audio{ID: 1, Title: "Frequent"}
	hr := &historymock.MockHistoryRepository{
		FindFrequentFn: func(userID uint, minPlays int, limit int) ([]entity.History, error) {
			return []entity.History{{ID: 1, UserID: userID, Audio: &audio, PlayCount: 5}}, nil
		},
	}
	svc := service.NewRecommendationService(&audiomock.MockAudioRepository{}, &audiomock.MockAudioScoreRepository{}, hr, &likemock.MockLikeRepository{}, &locationmock.MockUserLocationRepository{})
	audios, err := svc.GetQuickPick(1, 10)
	assert.NoError(t, err)
	assert.Len(t, audios, 1)
}

func TestRecommendationGetSimilar_Success(t *testing.T) {
	target := entity.Audio{ID: 1, Title: "Target", CategoryID: 1, Artist: "A"}
	similar := entity.Audio{ID: 2, Title: "Similar", CategoryID: 1, Artist: "A"}
	ar := &audiomock.MockAudioRepository{
		FindByIDFn:      func(id uint) (*entity.Audio, error) { return &target, nil },
		FindAllActiveFn: func() ([]entity.Audio, error) { return []entity.Audio{target, similar}, nil },
	}
	sr := &audiomock.MockAudioScoreRepository{
		FindByAudiosFn: func(audioIDs []uint) ([]entity.AudioScore, error) { return []entity.AudioScore{}, nil },
	}
	svc := service.NewRecommendationService(ar, sr, &historymock.MockHistoryRepository{}, &likemock.MockLikeRepository{}, &locationmock.MockUserLocationRepository{})
	audios, err := svc.GetSimilar(1, 5)
	assert.NoError(t, err)
	assert.Len(t, audios, 1)
	assert.Equal(t, "Similar", audios[0].Title)
}

func TestRecommendationRecalculateScores_Success(t *testing.T) {
	var savedScores []entity.AudioScore
	hr := &historymock.MockHistoryRepository{
		AggregateFn: func() (map[uint]int64, error) { return map[uint]int64{1: 100, 2: 50}, nil },
	}
	lr := &likemock.MockLikeRepository{
		AggregateFn: func() (map[uint]int64, error) { return map[uint]int64{1: 10, 2: 5}, nil },
	}
	sr := &audiomock.MockAudioScoreRepository{
		BulkUpsertFn: func(scores []entity.AudioScore) error {
			savedScores = scores
			return nil
		},
	}
	svc := service.NewRecommendationService(&audiomock.MockAudioRepository{}, sr, hr, lr, &locationmock.MockUserLocationRepository{})
	err := svc.RecalculateScores()
	assert.NoError(t, err)
	assert.Len(t, savedScores, 2)

	// Verify weight formula: plays*1.0 + likes*2.0
	scoreMap := make(map[uint]entity.AudioScore)
	for _, s := range savedScores {
		scoreMap[s.AudioID] = s
	}
	s1 := scoreMap[1]
	assert.Equal(t, int64(100), s1.TotalPlays)
	assert.Equal(t, int64(10), s1.TotalLikes)
	assert.Equal(t, float64(100)*1.0+float64(10)*2.0, s1.WeightScore, "weight = plays*1.0 + likes*2.0")

	s2 := scoreMap[2]
	assert.Equal(t, int64(50), s2.TotalPlays)
	assert.Equal(t, int64(5), s2.TotalLikes)
	assert.Equal(t, float64(50)*1.0+float64(5)*2.0, s2.WeightScore)
}

func TestRecommendationRecalculateScores_AggregatePlaysFails(t *testing.T) {
	playErr := errors.New("play aggregate error")
	hr := &historymock.MockHistoryRepository{
		AggregateFn: func() (map[uint]int64, error) { return nil, playErr },
	}
	svc := service.NewRecommendationService(&audiomock.MockAudioRepository{}, &audiomock.MockAudioScoreRepository{}, hr, &likemock.MockLikeRepository{}, &locationmock.MockUserLocationRepository{})
	err := svc.RecalculateScores()
	assert.ErrorIs(t, err, playErr)
}

func TestRecommendationRecalculateScores_AggregateLikesFails(t *testing.T) {
	likeErr := errors.New("like aggregate error")
	hr := &historymock.MockHistoryRepository{
		AggregateFn: func() (map[uint]int64, error) { return map[uint]int64{}, nil },
	}
	lr := &likemock.MockLikeRepository{
		AggregateFn: func() (map[uint]int64, error) { return nil, likeErr },
	}
	svc := service.NewRecommendationService(&audiomock.MockAudioRepository{}, &audiomock.MockAudioScoreRepository{}, hr, lr, &locationmock.MockUserLocationRepository{})
	err := svc.RecalculateScores()
	assert.ErrorIs(t, err, likeErr)
}

func TestRecommendationRecalculateScores_BulkUpsertFails(t *testing.T) {
	bulkErr := errors.New("bulk upsert failed")
	hr := &historymock.MockHistoryRepository{
		AggregateFn: func() (map[uint]int64, error) { return map[uint]int64{1: 10}, nil },
	}
	lr := &likemock.MockLikeRepository{
		AggregateFn: func() (map[uint]int64, error) { return map[uint]int64{1: 5}, nil },
	}
	sr := &audiomock.MockAudioScoreRepository{
		BulkUpsertFn: func(scores []entity.AudioScore) error { return bulkErr },
	}
	svc := service.NewRecommendationService(&audiomock.MockAudioRepository{}, sr, hr, lr, &locationmock.MockUserLocationRepository{})
	err := svc.RecalculateScores()
	assert.ErrorIs(t, err, bulkErr)
}

func TestRecommendationGetLocationBased_GPSNearby(t *testing.T) {
	loc := &locationmock.MockUserLocationRepository{
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
	hr := &historymock.MockHistoryRepository{
		FindByUsersFn: func(userIDs []uint) ([]entity.History, error) {
			a1 := entity.Audio{ID: 10, Title: "Jakarta Audio"}
			return []entity.History{{UserID: 2, AudioID: 10, Audio: &a1, PlayCount: 5}}, nil
		},
	}
	sr := &audiomock.MockAudioScoreRepository{
		FindTopFn: func(limit int) ([]entity.AudioScore, error) { return []entity.AudioScore{}, nil },
	}
	ar := &audiomock.MockAudioRepository{
		FindByIDsFn: func(ids []uint) ([]entity.Audio, error) {
			return []entity.Audio{{ID: 10, Title: "Jakarta Audio"}}, nil
		},
	}
	svc := service.NewRecommendationService(ar, sr, hr, &likemock.MockLikeRepository{}, loc)
	audios, err := svc.GetLocationBased(1, 10)
	assert.NoError(t, err)
	assert.NotEmpty(t, audios)
	assert.Equal(t, "Jakarta Audio", audios[0].Title)
}

func TestRecommendationGetLocationBased_NoGPS_FallbackCity(t *testing.T) {
	loc := &locationmock.MockUserLocationRepository{
		FindByUserFn: func(userID uint) (*entity.UserLocation, error) {
			return &entity.UserLocation{UserID: 1, Latitude: 0, Longitude: 0, City: "Jakarta"}, nil
		},
		FindAllFn: func() ([]entity.UserLocation, error) { return []entity.UserLocation{}, nil },
	}
	ar := &audiomock.MockAudioRepository{
		SearchFn: func(query string) ([]entity.Audio, error) {
			return []entity.Audio{{ID: 1, Title: "Jakarta Kajian"}}, nil
		},
	}
	sr := &audiomock.MockAudioScoreRepository{
		FindTopFn: func(limit int) ([]entity.AudioScore, error) { return []entity.AudioScore{}, nil },
	}
	svc := service.NewRecommendationService(ar, sr, &historymock.MockHistoryRepository{}, &likemock.MockLikeRepository{}, loc)
	audios, err := svc.GetLocationBased(1, 10)
	assert.NoError(t, err)
	assert.Len(t, audios, 1)
}

func TestRecommendationGetLocationBased_NoLocation_FallbackPopular(t *testing.T) {
	loc := &locationmock.MockUserLocationRepository{
		FindByUserFn: func(userID uint) (*entity.UserLocation, error) {
			return &entity.UserLocation{City: ""}, nil
		},
		FindAllFn: func() ([]entity.UserLocation, error) { return []entity.UserLocation{}, nil },
	}
	sr := &audiomock.MockAudioScoreRepository{
		FindTopFn: func(limit int) ([]entity.AudioScore, error) { return []entity.AudioScore{}, nil },
	}
	svc := service.NewRecommendationService(&audiomock.MockAudioRepository{}, sr, &historymock.MockHistoryRepository{}, &likemock.MockLikeRepository{}, loc)
	audios, err := svc.GetLocationBased(1, 10)
	assert.NoError(t, err)
	assert.Empty(t, audios)
}

func TestRecommendationGetTimeBasedPersonalized(t *testing.T) {
	sr := &audiomock.MockAudioScoreRepository{
		FindTopFn: func(limit int) ([]entity.AudioScore, error) { return []entity.AudioScore{}, nil },
	}
	hr := &historymock.MockHistoryRepository{
		FindFrequentFn: func(userID uint, minPlays int, limit int) ([]entity.History, error) { return []entity.History{}, nil },
		FindByUserFn:   func(userID uint) ([]entity.History, error) { return []entity.History{}, nil },
	}
	ar := &audiomock.MockAudioRepository{
		FindAllActiveFn: func() ([]entity.Audio, error) { return []entity.Audio{}, nil },
	}
	lr := &likemock.MockLikeRepository{
		FindByUserFn: func(userID uint, targetType string) ([]entity.Like, error) { return []entity.Like{}, nil },
	}
	svc := service.NewRecommendationService(ar, sr, hr, lr, &locationmock.MockUserLocationRepository{})
	for _, h := range []int{8, 14, 19, 23} {
		_, err := svc.GetTimeBasedPersonalized(1, h, 10)
		assert.NoError(t, err)
	}
}
