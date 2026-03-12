package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"mqfm-backend/internal/app/service"
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
	statmock "mqfm-backend/tests/mocks/stat"
)

// ──────────────────────── RecordStat ────────────────────────

func TestListeningStat_RecordStat_Success(t *testing.T) {
	var created *entity.ListeningStat
	repo := &statmock.MockListeningStatRepository{
		CreateFn: func(stat *entity.ListeningStat) error {
			created = stat
			assert.Equal(t, uint(1), stat.UserID)
			assert.Equal(t, uint(10), stat.AudioID)
			assert.Equal(t, 300, stat.DurationSeconds)
			assert.False(t, stat.ListenedAt.IsZero(), "ListenedAt must be set")
			return nil
		},
	}
	svc := service.NewListeningStatService(repo)

	err := svc.RecordStat(1, request.RecordStatRequest{AudioID: 10, DurationSeconds: 300})
	assert.NoError(t, err)
	assert.NotNil(t, created)
}

func TestListeningStat_RecordStat_Error(t *testing.T) {
	repo := &statmock.MockListeningStatRepository{
		CreateFn: func(stat *entity.ListeningStat) error {
			return errors.New("create failed")
		},
	}
	svc := service.NewListeningStatService(repo)

	err := svc.RecordStat(1, request.RecordStatRequest{AudioID: 10, DurationSeconds: 300})
	assert.Error(t, err)
	assert.Equal(t, "create failed", err.Error())
}

// ──────────────────────── GetWeeklySummary ────────────────────────

func TestListeningStat_GetWeeklySummary_Success(t *testing.T) {
	repo := &statmock.MockListeningStatRepository{
		GetWeeklySummaryFn: func(userID uint) (int, error) {
			assert.Equal(t, uint(1), userID)
			return 3600, nil
		},
	}
	svc := service.NewListeningStatService(repo)

	result, err := svc.GetWeeklySummary(1)
	assert.NoError(t, err)
	assert.Equal(t, 3600, result)
}

func TestListeningStat_GetWeeklySummary_Error(t *testing.T) {
	repo := &statmock.MockListeningStatRepository{
		GetWeeklySummaryFn: func(userID uint) (int, error) {
			return 0, errors.New("db error")
		},
	}
	svc := service.NewListeningStatService(repo)

	result, err := svc.GetWeeklySummary(1)
	assert.Error(t, err)
	assert.Equal(t, 0, result)
}

// ──────────────────────── GetMonthlySummary ────────────────────────

func TestListeningStat_GetMonthlySummary_Success(t *testing.T) {
	repo := &statmock.MockListeningStatRepository{
		GetMonthlySummaryFn: func(userID uint) (int, error) {
			return 14400, nil
		},
	}
	svc := service.NewListeningStatService(repo)

	result, err := svc.GetMonthlySummary(1)
	assert.NoError(t, err)
	assert.Equal(t, 14400, result)
}

// ──────────────────────── GetTopCategories ────────────────────────

func TestListeningStat_GetTopCategories_Success(t *testing.T) {
	expected := []port.CategoryStat{
		{CategoryID: 1, Name: "Tafsir", TotalTime: 5000},
		{CategoryID: 2, Name: "Fiqh", TotalTime: 3000},
	}
	repo := &statmock.MockListeningStatRepository{
		GetTopCategoriesFn: func(userID uint, limit int) ([]port.CategoryStat, error) {
			assert.Equal(t, uint(1), userID)
			assert.Equal(t, 5, limit)
			return expected, nil
		},
	}
	svc := service.NewListeningStatService(repo)

	result, err := svc.GetTopCategories(1, 5)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Tafsir", result[0].Name)
}

// ──────────────────────── GetTopArtists ────────────────────────

func TestListeningStat_GetTopArtists_Success(t *testing.T) {
	expected := []port.ArtistStat{
		{Artist: "Ust. Adi Hidayat", TotalTime: 8000},
		{Artist: "Ust. Abdul Somad", TotalTime: 6000},
	}
	repo := &statmock.MockListeningStatRepository{
		GetTopArtistsFn: func(userID uint, limit int) ([]port.ArtistStat, error) {
			assert.Equal(t, 10, limit)
			return expected, nil
		},
	}
	svc := service.NewListeningStatService(repo)

	result, err := svc.GetTopArtists(1, 10)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Ust. Adi Hidayat", result[0].Artist)
}

// ──────────────────────── GetDailySummary ────────────────────────

func TestListeningStat_GetDailySummary_Success(t *testing.T) {
	expected := []port.DailyStat{
		{Date: "2025-01-01", TotalTime: 1800},
		{Date: "2025-01-02", TotalTime: 2400},
	}
	repo := &statmock.MockListeningStatRepository{
		GetDailySummaryFn: func(userID uint, days int) ([]port.DailyStat, error) {
			assert.Equal(t, 7, days)
			return expected, nil
		},
	}
	svc := service.NewListeningStatService(repo)

	result, err := svc.GetDailySummary(1, 7)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

// ──────────────────────── GetRecap ────────────────────────

func TestListeningStat_GetRecap_ConvertsSecondsToMinutes(t *testing.T) {
	repo := &statmock.MockListeningStatRepository{
		GetWeeklySummaryFn: func(userID uint) (int, error) { return 7200, nil },     // 120 min
		GetMonthlySummaryFn: func(userID uint) (int, error) { return 36000, nil },   // 600 min
		GetTopCategoriesFn: func(userID uint, limit int) ([]port.CategoryStat, error) {
			assert.Equal(t, 5, limit, "recap uses top 5 categories")
			return []port.CategoryStat{{Name: "Tafsir", TotalTime: 5000}}, nil
		},
		GetTopArtistsFn: func(userID uint, limit int) ([]port.ArtistStat, error) {
			assert.Equal(t, 5, limit, "recap uses top 5 artists")
			return []port.ArtistStat{{Artist: "Ust. Adi Hidayat", TotalTime: 8000}}, nil
		},
		GetDailySummaryFn: func(userID uint, days int) ([]port.DailyStat, error) {
			assert.Equal(t, 30, days, "recap uses 30 days")
			return []port.DailyStat{{Date: "2025-01-01", TotalTime: 1800}}, nil
		},
	}
	svc := service.NewListeningStatService(repo)

	recap, err := svc.GetRecap(1)
	assert.NoError(t, err)
	assert.NotNil(t, recap)
	assert.Equal(t, 120, recap.WeeklyMinutes, "7200 seconds / 60 = 120 minutes")
	assert.Equal(t, 600, recap.MonthlyMinutes, "36000 seconds / 60 = 600 minutes")
	assert.Len(t, recap.TopCategories, 1)
	assert.Equal(t, "Tafsir", recap.TopCategories[0].Name)
	assert.Len(t, recap.TopArtists, 1)
	assert.Equal(t, "Ust. Adi Hidayat", recap.TopArtists[0].Artist)
	assert.Len(t, recap.DailyStats, 1)
}

func TestListeningStat_GetRecap_ZeroValues(t *testing.T) {
	repo := &statmock.MockListeningStatRepository{
		GetWeeklySummaryFn:  func(userID uint) (int, error) { return 0, nil },
		GetMonthlySummaryFn: func(userID uint) (int, error) { return 0, nil },
		GetTopCategoriesFn:  func(userID uint, limit int) ([]port.CategoryStat, error) { return nil, nil },
		GetTopArtistsFn:     func(userID uint, limit int) ([]port.ArtistStat, error) { return nil, nil },
		GetDailySummaryFn:   func(userID uint, days int) ([]port.DailyStat, error) { return nil, nil },
	}
	svc := service.NewListeningStatService(repo)

	recap, err := svc.GetRecap(1)
	assert.NoError(t, err)
	assert.Equal(t, 0, recap.WeeklyMinutes)
	assert.Equal(t, 0, recap.MonthlyMinutes)
	assert.Nil(t, recap.TopCategories)
	assert.Nil(t, recap.TopArtists)
	assert.Nil(t, recap.DailyStats)
}

func TestListeningStat_GetRecap_IntegerDivision(t *testing.T) {
	repo := &statmock.MockListeningStatRepository{
		GetWeeklySummaryFn:  func(userID uint) (int, error) { return 125, nil }, // 2.08 min → 2
		GetMonthlySummaryFn: func(userID uint) (int, error) { return 59, nil },  // 0.98 min → 0
		GetTopCategoriesFn:  func(userID uint, limit int) ([]port.CategoryStat, error) { return nil, nil },
		GetTopArtistsFn:     func(userID uint, limit int) ([]port.ArtistStat, error) { return nil, nil },
		GetDailySummaryFn:   func(userID uint, days int) ([]port.DailyStat, error) { return nil, nil },
	}
	svc := service.NewListeningStatService(repo)

	recap, err := svc.GetRecap(1)
	assert.NoError(t, err)
	assert.Equal(t, 2, recap.WeeklyMinutes, "125/60 = 2 (integer division)")
	assert.Equal(t, 0, recap.MonthlyMinutes, "59/60 = 0 (integer division)")
}
