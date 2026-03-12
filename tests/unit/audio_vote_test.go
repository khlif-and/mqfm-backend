package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"mqfm-backend/internal/app/service"
	"mqfm-backend/internal/domain/entity"
	audiomock "mqfm-backend/tests/mocks/audio"
)

// ──────────────────────── Vote ────────────────────────

func TestVote_Success_NewRanking(t *testing.T) {
	var upsertedRanking *entity.AudioRanking
	voteRepo := &audiomock.MockAudioVoteRepository{
		ExistsFn:       func(userID, audioID uint) (bool, error) { return false, nil },
		CreateFn:       func(vote *entity.AudioVote) error { return nil },
		CountByAudioFn: func(audioID uint) (int64, error) { return 5, nil },
		CountWeeklyFn:  func(audioID uint) (int64, error) { return 3, nil },
		CountMonthlyFn: func(audioID uint) (int64, error) { return 4, nil },
	}
	rankingRepo := &audiomock.MockAudioRankingRepository{
		FindByAudioFn: func(audioID uint) (*entity.AudioRanking, error) {
			return nil, errors.New("not found") // forces new ranking creation
		},
		UpsertFn: func(ranking *entity.AudioRanking) error {
			upsertedRanking = ranking
			return nil
		},
	}
	audioRepo := &audiomock.MockAudioRepository{} // unused for Vote

	svc := service.NewAudioVoteService(voteRepo, rankingRepo, audioRepo)
	err := svc.Vote(1, 100)

	assert.NoError(t, err)
	assert.NotNil(t, upsertedRanking)
	assert.Equal(t, uint(100), upsertedRanking.AudioID)
	assert.Equal(t, int64(5), upsertedRanking.TotalVotes)
	assert.Equal(t, int64(3), upsertedRanking.WeeklyVotes)
	assert.Equal(t, int64(4), upsertedRanking.MonthlyVotes)
	assert.True(t, upsertedRanking.RandomBoost >= 0 && upsertedRanking.RandomBoost <= 100)
}

func TestVote_Success_ExistingRanking(t *testing.T) {
	existing := &entity.AudioRanking{
		ID:          1,
		AudioID:     100,
		RandomBoost: 42.5,
	}
	var upsertedRanking *entity.AudioRanking
	voteRepo := &audiomock.MockAudioVoteRepository{
		ExistsFn:       func(userID, audioID uint) (bool, error) { return false, nil },
		CreateFn:       func(vote *entity.AudioVote) error { return nil },
		CountByAudioFn: func(audioID uint) (int64, error) { return 10, nil },
		CountWeeklyFn:  func(audioID uint) (int64, error) { return 7, nil },
		CountMonthlyFn: func(audioID uint) (int64, error) { return 8, nil },
	}
	rankingRepo := &audiomock.MockAudioRankingRepository{
		FindByAudioFn: func(audioID uint) (*entity.AudioRanking, error) {
			return existing, nil
		},
		UpsertFn: func(ranking *entity.AudioRanking) error {
			upsertedRanking = ranking
			return nil
		},
	}

	svc := service.NewAudioVoteService(voteRepo, rankingRepo, &audiomock.MockAudioRepository{})
	err := svc.Vote(1, 100)

	assert.NoError(t, err)
	assert.NotNil(t, upsertedRanking)
	// RandomBoost preserved from existing
	assert.Equal(t, 42.5, upsertedRanking.RandomBoost)
	assert.Equal(t, int64(10), upsertedRanking.TotalVotes)
	assert.Equal(t, int64(7), upsertedRanking.WeeklyVotes)
	assert.Equal(t, int64(8), upsertedRanking.MonthlyVotes)
}

func TestVote_AlreadyVoted(t *testing.T) {
	voteRepo := &audiomock.MockAudioVoteRepository{
		ExistsFn: func(userID, audioID uint) (bool, error) { return true, nil },
	}
	svc := service.NewAudioVoteService(voteRepo, &audiomock.MockAudioRankingRepository{}, &audiomock.MockAudioRepository{})
	err := svc.Vote(1, 100)

	assert.Error(t, err)
	assert.Equal(t, "already voted", err.Error())
}

func TestVote_ExistsCheckFails(t *testing.T) {
	dbErr := errors.New("connection lost")
	voteRepo := &audiomock.MockAudioVoteRepository{
		ExistsFn: func(userID, audioID uint) (bool, error) { return false, dbErr },
	}
	svc := service.NewAudioVoteService(voteRepo, &audiomock.MockAudioRankingRepository{}, &audiomock.MockAudioRepository{})
	err := svc.Vote(1, 100)

	assert.ErrorIs(t, err, dbErr)
}

func TestVote_CreateFails(t *testing.T) {
	createErr := errors.New("insert failed")
	voteRepo := &audiomock.MockAudioVoteRepository{
		ExistsFn: func(userID, audioID uint) (bool, error) { return false, nil },
		CreateFn: func(vote *entity.AudioVote) error { return createErr },
	}
	svc := service.NewAudioVoteService(voteRepo, &audiomock.MockAudioRankingRepository{}, &audiomock.MockAudioRepository{})
	err := svc.Vote(1, 100)

	assert.ErrorIs(t, err, createErr)
}

func TestVote_UpsertFails(t *testing.T) {
	upsertErr := errors.New("upsert failed")
	voteRepo := &audiomock.MockAudioVoteRepository{
		ExistsFn:       func(userID, audioID uint) (bool, error) { return false, nil },
		CreateFn:       func(vote *entity.AudioVote) error { return nil },
		CountByAudioFn: func(audioID uint) (int64, error) { return 1, nil },
		CountWeeklyFn:  func(audioID uint) (int64, error) { return 1, nil },
		CountMonthlyFn: func(audioID uint) (int64, error) { return 1, nil },
	}
	rankingRepo := &audiomock.MockAudioRankingRepository{
		FindByAudioFn: func(audioID uint) (*entity.AudioRanking, error) {
			return &entity.AudioRanking{AudioID: audioID}, nil
		},
		UpsertFn: func(ranking *entity.AudioRanking) error { return upsertErr },
	}
	svc := service.NewAudioVoteService(voteRepo, rankingRepo, &audiomock.MockAudioRepository{})
	err := svc.Vote(1, 100)

	assert.ErrorIs(t, err, upsertErr)
}

// ──────────────────────── Unvote ────────────────────────

func TestUnvote_Success(t *testing.T) {
	var upsertedRanking *entity.AudioRanking
	voteRepo := &audiomock.MockAudioVoteRepository{
		DeleteFn:       func(userID, audioID uint) error { return nil },
		CountByAudioFn: func(audioID uint) (int64, error) { return 4, nil },
		CountWeeklyFn:  func(audioID uint) (int64, error) { return 2, nil },
		CountMonthlyFn: func(audioID uint) (int64, error) { return 3, nil },
	}
	rankingRepo := &audiomock.MockAudioRankingRepository{
		FindByAudioFn: func(audioID uint) (*entity.AudioRanking, error) {
			return &entity.AudioRanking{AudioID: audioID, TotalVotes: 5}, nil
		},
		UpsertFn: func(ranking *entity.AudioRanking) error {
			upsertedRanking = ranking
			return nil
		},
	}
	svc := service.NewAudioVoteService(voteRepo, rankingRepo, &audiomock.MockAudioRepository{})
	err := svc.Unvote(1, 100)

	assert.NoError(t, err)
	assert.NotNil(t, upsertedRanking)
	assert.Equal(t, int64(4), upsertedRanking.TotalVotes, "votes must be decremented")
	assert.Equal(t, int64(2), upsertedRanking.WeeklyVotes)
	assert.Equal(t, int64(3), upsertedRanking.MonthlyVotes)
}

func TestUnvote_DeleteFails(t *testing.T) {
	delErr := errors.New("delete failed")
	voteRepo := &audiomock.MockAudioVoteRepository{
		DeleteFn: func(userID, audioID uint) error { return delErr },
	}
	svc := service.NewAudioVoteService(voteRepo, &audiomock.MockAudioRankingRepository{}, &audiomock.MockAudioRepository{})
	err := svc.Unvote(1, 100)

	assert.ErrorIs(t, err, delErr)
}

func TestUnvote_RankingNotFound_NoError(t *testing.T) {
	voteRepo := &audiomock.MockAudioVoteRepository{
		DeleteFn: func(userID, audioID uint) error { return nil },
	}
	rankingRepo := &audiomock.MockAudioRankingRepository{
		FindByAudioFn: func(audioID uint) (*entity.AudioRanking, error) {
			return nil, errors.New("not found")
		},
	}
	svc := service.NewAudioVoteService(voteRepo, rankingRepo, &audiomock.MockAudioRepository{})
	err := svc.Unvote(1, 100)

	assert.NoError(t, err, "unvote should succeed even if ranking not found")
}

// ──────────────────────── GetWeeklyRanking ────────────────────────

func TestGetWeeklyRanking_Success(t *testing.T) {
	expected := []entity.AudioRanking{
		{AudioID: 1, WeeklyVotes: 100, WeeklyRank: 1},
		{AudioID: 2, WeeklyVotes: 80, WeeklyRank: 2},
	}
	rankingRepo := &audiomock.MockAudioRankingRepository{
		FindTopWeeklyFn: func(limit int) ([]entity.AudioRanking, error) {
			assert.Equal(t, 10, limit)
			return expected, nil
		},
	}
	svc := service.NewAudioVoteService(&audiomock.MockAudioVoteRepository{}, rankingRepo, &audiomock.MockAudioRepository{})
	rankings, err := svc.GetWeeklyRanking(10)

	assert.NoError(t, err)
	assert.Equal(t, expected, rankings)
}

func TestGetWeeklyRanking_LimitCapped(t *testing.T) {
	rankingRepo := &audiomock.MockAudioRankingRepository{
		FindTopWeeklyFn: func(limit int) ([]entity.AudioRanking, error) {
			assert.Equal(t, 20, limit, "limit must be capped at displayTopN=20")
			return []entity.AudioRanking{}, nil
		},
	}
	svc := service.NewAudioVoteService(&audiomock.MockAudioVoteRepository{}, rankingRepo, &audiomock.MockAudioRepository{})
	_, err := svc.GetWeeklyRanking(100)
	assert.NoError(t, err)
}

// ──────────────────────── GetMonthlyRanking ────────────────────────

func TestGetMonthlyRanking_Success(t *testing.T) {
	expected := []entity.AudioRanking{
		{AudioID: 1, MonthlyVotes: 500, MonthlyRank: 1},
		{AudioID: 2, MonthlyVotes: 400, MonthlyRank: 2},
	}
	rankingRepo := &audiomock.MockAudioRankingRepository{
		FindTopMonthlyFn: func(limit int) ([]entity.AudioRanking, error) {
			assert.Equal(t, 5, limit)
			return expected, nil
		},
	}
	svc := service.NewAudioVoteService(&audiomock.MockAudioVoteRepository{}, rankingRepo, &audiomock.MockAudioRepository{})
	rankings, err := svc.GetMonthlyRanking(5)

	assert.NoError(t, err)
	assert.Equal(t, expected, rankings)
}

func TestGetMonthlyRanking_LimitCapped(t *testing.T) {
	rankingRepo := &audiomock.MockAudioRankingRepository{
		FindTopMonthlyFn: func(limit int) ([]entity.AudioRanking, error) {
			assert.Equal(t, 20, limit, "limit must be capped at displayTopN=20")
			return []entity.AudioRanking{}, nil
		},
	}
	svc := service.NewAudioVoteService(&audiomock.MockAudioVoteRepository{}, rankingRepo, &audiomock.MockAudioRepository{})
	_, err := svc.GetMonthlyRanking(50)
	assert.NoError(t, err)
}

// ──────────────────────── HasVoted ────────────────────────

func TestHasVoted_True(t *testing.T) {
	voteRepo := &audiomock.MockAudioVoteRepository{
		ExistsFn: func(userID, audioID uint) (bool, error) {
			assert.Equal(t, uint(1), userID)
			assert.Equal(t, uint(100), audioID)
			return true, nil
		},
	}
	svc := service.NewAudioVoteService(voteRepo, &audiomock.MockAudioRankingRepository{}, &audiomock.MockAudioRepository{})
	voted, err := svc.HasVoted(1, 100)

	assert.NoError(t, err)
	assert.True(t, voted)
}

func TestHasVoted_False(t *testing.T) {
	voteRepo := &audiomock.MockAudioVoteRepository{
		ExistsFn: func(userID, audioID uint) (bool, error) { return false, nil },
	}
	svc := service.NewAudioVoteService(voteRepo, &audiomock.MockAudioRankingRepository{}, &audiomock.MockAudioRepository{})
	voted, err := svc.HasVoted(1, 100)

	assert.NoError(t, err)
	assert.False(t, voted)
}

func TestHasVoted_Error(t *testing.T) {
	dbErr := errors.New("db err")
	voteRepo := &audiomock.MockAudioVoteRepository{
		ExistsFn: func(userID, audioID uint) (bool, error) { return false, dbErr },
	}
	svc := service.NewAudioVoteService(voteRepo, &audiomock.MockAudioRankingRepository{}, &audiomock.MockAudioRepository{})
	_, err := svc.HasVoted(1, 100)

	assert.ErrorIs(t, err, dbErr)
}

// ──────────────────────── RecalculateRankings ────────────────────────

func TestRecalculateRankings_Success(t *testing.T) {
	var bulkUpsertCalled int
	audioRepo := &audiomock.MockAudioRepository{
		CountAllFn: func() (int64, error) { return 2, nil },
	}
	rankingRepo := &audiomock.MockAudioRankingRepository{
		FindAllFn: func(limit, offset int) ([]entity.AudioRanking, error) {
			return []entity.AudioRanking{
				{AudioID: 1, RandomBoost: 10.0},
				{AudioID: 2, RandomBoost: 5.0},
			}, nil
		},
		BulkUpsertFn: func(rankings []entity.AudioRanking) error {
			bulkUpsertCalled++
			// Verify ranks assigned
			for _, r := range rankings {
				assert.True(t, r.WeeklyRank > 0, "weekly rank must be assigned")
				assert.True(t, r.MonthlyRank > 0, "monthly rank must be assigned")
			}
			return nil
		},
	}
	voteRepo := &audiomock.MockAudioVoteRepository{
		CountByAudioFn: func(audioID uint) (int64, error) {
			if audioID == 1 {
				return 50, nil
			}
			return 20, nil
		},
		CountWeeklyFn: func(audioID uint) (int64, error) {
			if audioID == 1 {
				return 10, nil
			}
			return 5, nil
		},
		CountMonthlyFn: func(audioID uint) (int64, error) {
			if audioID == 1 {
				return 30, nil
			}
			return 15, nil
		},
	}

	svc := service.NewAudioVoteService(voteRepo, rankingRepo, audioRepo)
	err := svc.RecalculateRankings()

	assert.NoError(t, err)
	assert.Equal(t, 1, bulkUpsertCalled, "batch upsert must be called once for 2 items")
}

func TestRecalculateRankings_CountAllFails(t *testing.T) {
	countErr := errors.New("count failed")
	audioRepo := &audiomock.MockAudioRepository{
		CountAllFn: func() (int64, error) { return 0, countErr },
	}
	svc := service.NewAudioVoteService(&audiomock.MockAudioVoteRepository{}, &audiomock.MockAudioRankingRepository{}, audioRepo)
	err := svc.RecalculateRankings()

	assert.ErrorIs(t, err, countErr)
}

// ──────────────────────── ResetWeeklyVotes ────────────────────────

func TestResetWeeklyVotes_Success(t *testing.T) {
	var savedRankings []entity.AudioRanking
	rankingRepo := &audiomock.MockAudioRankingRepository{
		FindAllFn: func(limit, offset int) ([]entity.AudioRanking, error) {
			return []entity.AudioRanking{
				{AudioID: 1, WeeklyVotes: 100, WeeklyRank: 1, RandomBoost: 50.0},
				{AudioID: 2, WeeklyVotes: 80, WeeklyRank: 2, RandomBoost: 30.0},
			}, nil
		},
		BulkUpsertFn: func(rankings []entity.AudioRanking) error {
			savedRankings = rankings
			return nil
		},
	}
	svc := service.NewAudioVoteService(&audiomock.MockAudioVoteRepository{}, rankingRepo, &audiomock.MockAudioRepository{})
	err := svc.ResetWeeklyVotes()

	assert.NoError(t, err)
	assert.Len(t, savedRankings, 2)
	for _, r := range savedRankings {
		assert.Equal(t, int64(0), r.WeeklyVotes, "weekly votes must be reset to 0")
		assert.Equal(t, 0, r.WeeklyRank, "weekly rank must be reset to 0")
		// RandomBoost is re-rolled, just check it is in valid range
		assert.True(t, r.RandomBoost >= 0 && r.RandomBoost <= 100, "random boost re-rolled in [0,100)")
	}
}

func TestResetWeeklyVotes_FindAllFails(t *testing.T) {
	findErr := errors.New("find failed")
	rankingRepo := &audiomock.MockAudioRankingRepository{
		FindAllFn: func(limit, offset int) ([]entity.AudioRanking, error) { return nil, findErr },
	}
	svc := service.NewAudioVoteService(&audiomock.MockAudioVoteRepository{}, rankingRepo, &audiomock.MockAudioRepository{})
	err := svc.ResetWeeklyVotes()

	assert.ErrorIs(t, err, findErr)
}

// ──────────────────────── ResetMonthlyVotes ────────────────────────

func TestResetMonthlyVotes_Success(t *testing.T) {
	var savedRankings []entity.AudioRanking
	rankingRepo := &audiomock.MockAudioRankingRepository{
		FindAllFn: func(limit, offset int) ([]entity.AudioRanking, error) {
			return []entity.AudioRanking{
				{AudioID: 1, MonthlyVotes: 500, MonthlyRank: 1},
				{AudioID: 2, MonthlyVotes: 400, MonthlyRank: 2},
			}, nil
		},
		BulkUpsertFn: func(rankings []entity.AudioRanking) error {
			savedRankings = rankings
			return nil
		},
	}
	svc := service.NewAudioVoteService(&audiomock.MockAudioVoteRepository{}, rankingRepo, &audiomock.MockAudioRepository{})
	err := svc.ResetMonthlyVotes()

	assert.NoError(t, err)
	assert.Len(t, savedRankings, 2)
	for _, r := range savedRankings {
		assert.Equal(t, int64(0), r.MonthlyVotes, "monthly votes must be reset to 0")
		assert.Equal(t, 0, r.MonthlyRank, "monthly rank must be reset to 0")
	}
}

func TestResetMonthlyVotes_FindAllFails(t *testing.T) {
	findErr := errors.New("find failed")
	rankingRepo := &audiomock.MockAudioRankingRepository{
		FindAllFn: func(limit, offset int) ([]entity.AudioRanking, error) { return nil, findErr },
	}
	svc := service.NewAudioVoteService(&audiomock.MockAudioVoteRepository{}, rankingRepo, &audiomock.MockAudioRepository{})
	err := svc.ResetMonthlyVotes()

	assert.ErrorIs(t, err, findErr)
}
