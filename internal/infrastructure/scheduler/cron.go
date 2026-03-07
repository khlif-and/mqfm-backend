package scheduler

import (
	"time"

	"go.uber.org/zap"

	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/logger"
)

type ScoreRecalculator struct {
	recommendationSvc port.RecommendationService
	targetHour        int
}

func NewScoreRecalculator(svc port.RecommendationService, targetHour int) *ScoreRecalculator {
	return &ScoreRecalculator{
		recommendationSvc: svc,
		targetHour:        targetHour,
	}
}

func (s *ScoreRecalculator) Start() {
	go func() {
		logger.Info("scheduler started: daily score recalculation", zap.Int("target_hour", s.targetHour))
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), s.targetHour, 0, 0, 0, now.Location())
			if now.After(next) {
				next = next.Add(24 * time.Hour)
			}

			sleepDuration := time.Until(next)
			logger.Info("scheduler: next recalculation", zap.String("at", next.Format(time.RFC3339)))
			time.Sleep(sleepDuration)

			logger.Info("scheduler: starting score recalculation")
			start := time.Now()

			if err := s.recommendationSvc.RecalculateScores(); err != nil {
				logger.Error("scheduler: recalculation failed", zap.Error(err))
			} else {
				elapsed := time.Since(start)
				logger.Info("scheduler: recalculation completed", zap.Duration("elapsed", elapsed))
			}
		}
	}()
}
