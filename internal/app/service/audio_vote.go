package service

import (
	"errors"
	"math/rand"
	"sort"

	"go.uber.org/zap"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/logger"
)

const (
	maxRankedAudios = 15000
	displayTopN     = 20
)

type audioVoteService struct {
	voteRepo    port.AudioVoteRepository
	rankingRepo port.AudioRankingRepository
	audioRepo   port.AudioRepository
}

func NewAudioVoteService(
	voteRepo port.AudioVoteRepository,
	rankingRepo port.AudioRankingRepository,
	audioRepo port.AudioRepository,
) port.AudioVoteService {
	return &audioVoteService{
		voteRepo:    voteRepo,
		rankingRepo: rankingRepo,
		audioRepo:   audioRepo,
	}
}

func (s *audioVoteService) Vote(userID, audioID uint) error {
	exists, err := s.voteRepo.Exists(userID, audioID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already voted")
	}

	if err := s.voteRepo.Create(&entity.AudioVote{
		UserID:  userID,
		AudioID: audioID,
	}); err != nil {
		return err
	}

	ranking, err := s.rankingRepo.FindByAudioID(audioID)
	if err != nil {
		ranking = &entity.AudioRanking{
			AudioID:     audioID,
			RandomBoost: rand.Float64() * 100,
		}
	}

	total, _ := s.voteRepo.CountByAudio(audioID)
	weekly, _ := s.voteRepo.CountWeeklyByAudio(audioID)
	monthly, _ := s.voteRepo.CountMonthlyByAudio(audioID)

	ranking.TotalVotes = total
	ranking.WeeklyVotes = weekly
	ranking.MonthlyVotes = monthly

	return s.rankingRepo.Upsert(ranking)
}

func (s *audioVoteService) Unvote(userID, audioID uint) error {
	if err := s.voteRepo.Delete(userID, audioID); err != nil {
		return err
	}

	ranking, err := s.rankingRepo.FindByAudioID(audioID)
	if err != nil {
		return nil
	}

	total, _ := s.voteRepo.CountByAudio(audioID)
	weekly, _ := s.voteRepo.CountWeeklyByAudio(audioID)
	monthly, _ := s.voteRepo.CountMonthlyByAudio(audioID)

	ranking.TotalVotes = total
	ranking.WeeklyVotes = weekly
	ranking.MonthlyVotes = monthly

	return s.rankingRepo.Upsert(ranking)
}

func (s *audioVoteService) GetWeeklyRanking(limit int) ([]entity.AudioRanking, error) {
	if limit > displayTopN {
		limit = displayTopN
	}
	return s.rankingRepo.FindTopWeekly(limit)
}

func (s *audioVoteService) GetMonthlyRanking(limit int) ([]entity.AudioRanking, error) {
	if limit > displayTopN {
		limit = displayTopN
	}
	return s.rankingRepo.FindTopMonthly(limit)
}

func (s *audioVoteService) HasVoted(userID, audioID uint) (bool, error) {
	return s.voteRepo.Exists(userID, audioID)
}

func (s *audioVoteService) RecalculateRankings() error {
	totalAudios, err := s.audioRepo.CountAll()
	if err != nil {
		return err
	}

	batchSize := 500
	limit := int(totalAudios)
	if limit > maxRankedAudios {
		limit = maxRankedAudios
	}

	allRankings, err := s.rankingRepo.FindAll(limit, 0)
	if err != nil {
		return err
	}

	for i := range allRankings {
		r := &allRankings[i]
		weekly, _ := s.voteRepo.CountWeeklyByAudio(r.AudioID)
		monthly, _ := s.voteRepo.CountMonthlyByAudio(r.AudioID)
		total, _ := s.voteRepo.CountByAudio(r.AudioID)

		r.WeeklyVotes = weekly
		r.MonthlyVotes = monthly
		r.TotalVotes = total
	}

	sort.Slice(allRankings, func(i, j int) bool {
		scoreI := float64(allRankings[i].WeeklyVotes) + allRankings[i].RandomBoost
		scoreJ := float64(allRankings[j].WeeklyVotes) + allRankings[j].RandomBoost
		return scoreI > scoreJ
	})
	for i := range allRankings {
		allRankings[i].WeeklyRank = i + 1
	}

	sort.Slice(allRankings, func(i, j int) bool {
		scoreI := float64(allRankings[i].MonthlyVotes) + allRankings[i].RandomBoost
		scoreJ := float64(allRankings[j].MonthlyVotes) + allRankings[j].RandomBoost
		return scoreI > scoreJ
	})
	for i := range allRankings {
		allRankings[i].MonthlyRank = i + 1
	}

	for i := 0; i < len(allRankings); i += batchSize {
		end := i + batchSize
		if end > len(allRankings) {
			end = len(allRankings)
		}
		if err := s.rankingRepo.BulkUpsert(allRankings[i:end]); err != nil {
			logger.Error("ranking batch upsert failed", zap.Error(err))
		}
	}

	logger.Info("rankings recalculated", zap.Int("total", len(allRankings)))
	return nil
}

func (s *audioVoteService) ResetWeeklyVotes() error {
	allRankings, err := s.rankingRepo.FindAll(maxRankedAudios, 0)
	if err != nil {
		return err
	}

	for i := range allRankings {
		allRankings[i].WeeklyVotes = 0
		allRankings[i].WeeklyRank = 0
		allRankings[i].RandomBoost = rand.Float64() * 100
	}

	return s.rankingRepo.BulkUpsert(allRankings)
}

func (s *audioVoteService) ResetMonthlyVotes() error {
	allRankings, err := s.rankingRepo.FindAll(maxRankedAudios, 0)
	if err != nil {
		return err
	}

	for i := range allRankings {
		allRankings[i].MonthlyVotes = 0
		allRankings[i].MonthlyRank = 0
	}

	return s.rankingRepo.BulkUpsert(allRankings)
}
