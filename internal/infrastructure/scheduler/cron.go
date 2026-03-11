package scheduler

import (
	"context"
	"math/rand"
	"time"

	"go.uber.org/zap"

	"mqfm-backend/internal/adapter/cache"
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

type RankingScheduler struct {
	rankingCache *cache.RankingCache
}

func NewRankingScheduler(rc *cache.RankingCache) *RankingScheduler {
	return &RankingScheduler{rankingCache: rc}
}

func (s *RankingScheduler) Start() {
	go s.rankingRefresh()
	go s.popularRefresh()
}

func (s *RankingScheduler) rankingRefresh() {
	logger.Info("scheduler started: ranking cache refresh every 5 hours")
	ctx := context.Background()
	_, _ = s.rankingCache.RefreshRanking(ctx)
	for {
		time.Sleep(5 * time.Hour)
		logger.Info("scheduler: refreshing ranking cache")
		if _, err := s.rankingCache.RefreshRanking(ctx); err != nil {
			logger.Error("scheduler: ranking refresh failed", zap.Error(err))
		}
	}
}

func (s *RankingScheduler) popularRefresh() {
	logger.Info("scheduler started: popular cache refresh every 7 hours")
	ctx := context.Background()
	_, _ = s.rankingCache.RefreshPopular(ctx)
	for {
		time.Sleep(7 * time.Hour)
		logger.Info("scheduler: refreshing popular cache")
		if _, err := s.rankingCache.RefreshPopular(ctx); err != nil {
			logger.Error("scheduler: popular refresh failed", zap.Error(err))
		}
	}
}

type VotingScheduler struct {
	voteSvc port.AudioVoteService
}

func NewVotingScheduler(svc port.AudioVoteService) *VotingScheduler {
	return &VotingScheduler{voteSvc: svc}
}

func (v *VotingScheduler) Start() {
	go v.dailyRankingRecalculation()
	go v.weeklyReset()
	go v.monthlyReset()
}

func (v *VotingScheduler) dailyRankingRecalculation() {
	logger.Info("scheduler started: daily voting ranking recalculation at 02:00")
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
		if now.After(next) {
			next = next.Add(24 * time.Hour)
		}
		time.Sleep(time.Until(next))

		logger.Info("scheduler: starting voting ranking recalculation")
		start := time.Now()
		if err := v.voteSvc.RecalculateRankings(); err != nil {
			logger.Error("scheduler: voting ranking recalculation failed", zap.Error(err))
		} else {
			logger.Info("scheduler: voting ranking recalculation completed", zap.Duration("elapsed", time.Since(start)))
		}
	}
}

func (v *VotingScheduler) weeklyReset() {
	logger.Info("scheduler started: weekly vote reset on Monday 00:00")
	for {
		now := time.Now()
		daysUntilMonday := (8 - int(now.Weekday())) % 7
		if daysUntilMonday == 0 && now.Hour() >= 0 {
			daysUntilMonday = 7
		}
		next := time.Date(now.Year(), now.Month(), now.Day()+daysUntilMonday, 0, 0, 0, 0, now.Location())
		time.Sleep(time.Until(next))

		logger.Info("scheduler: starting weekly vote reset")
		if err := v.voteSvc.ResetWeeklyVotes(); err != nil {
			logger.Error("scheduler: weekly vote reset failed", zap.Error(err))
		} else {
			logger.Info("scheduler: weekly vote reset completed")
		}
	}
}

func (v *VotingScheduler) monthlyReset() {
	logger.Info("scheduler started: monthly vote reset on 1st 00:00")
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
		time.Sleep(time.Until(next))

		logger.Info("scheduler: starting monthly vote reset")
		if err := v.voteSvc.ResetMonthlyVotes(); err != nil {
			logger.Error("scheduler: monthly vote reset failed", zap.Error(err))
		} else {
			logger.Info("scheduler: monthly vote reset completed")
		}
	}
}

type NotificationScheduler struct {
	notificationSvc port.NotificationService
}

func NewNotificationScheduler(svc port.NotificationService) *NotificationScheduler {
	return &NotificationScheduler{notificationSvc: svc}
}

func (n *NotificationScheduler) Start() {
	go func() {
		logger.Info("scheduler started: daily notification reminder at 08:00")
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location())
			if now.After(next) {
				next = next.Add(24 * time.Hour)
			}
			time.Sleep(time.Until(next))

			logger.Info("scheduler: sending daily reminder notifications")
			if err := n.notificationSvc.NotifyDailyReminder(); err != nil {
				logger.Error("scheduler: daily reminder failed", zap.Error(err))
			} else {
				logger.Info("scheduler: daily reminder completed")
			}
		}
	}()
}

type DownloadCleanupScheduler struct {
	downloadSvc port.DownloadService
}

func NewDownloadCleanupScheduler(svc port.DownloadService) *DownloadCleanupScheduler {
	return &DownloadCleanupScheduler{downloadSvc: svc}
}

func (d *DownloadCleanupScheduler) Start() {
	go func() {
		logger.Info("scheduler started: daily download cleanup at 04:00")
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), 4, 0, 0, 0, now.Location())
			if now.After(next) {
				next = next.Add(24 * time.Hour)
			}
			time.Sleep(time.Until(next))

			logger.Info("scheduler: cleaning up expired downloads")
			if count, err := d.downloadSvc.CleanupExpired(); err != nil {
				logger.Error("scheduler: download cleanup failed", zap.Error(err))
			} else {
				logger.Info("scheduler: download cleanup completed", zap.Int64("removed", count))
			}
		}
	}()
}

// WeeklyLikeScheduler refreshes weekly_likes every Monday with real likes +
// a deterministic robot boost (seeded by week number + audio_id) so the ranking
// rotates each week without being fully random.
type WeeklyLikeScheduler struct {
	scoreRepo port.AudioScoreRepository
	likeRepo  port.LikeRepository
}

func NewWeeklyLikeScheduler(sr port.AudioScoreRepository, lr port.LikeRepository) *WeeklyLikeScheduler {
	return &WeeklyLikeScheduler{scoreRepo: sr, likeRepo: lr}
}

func (w *WeeklyLikeScheduler) Start() {
	go func() {
		logger.Info("scheduler started: weekly like ranking refresh")
		// Run immediately so ranking is populated on startup.
		w.refresh()
		for {
			now := time.Now()
			// Next Monday 00:01
			daysUntilMonday := (int(time.Monday) - int(now.Weekday()) + 7) % 7
			if daysUntilMonday == 0 {
				daysUntilMonday = 7
			}
			next := time.Date(now.Year(), now.Month(), now.Day()+daysUntilMonday, 0, 1, 0, 0, now.Location())
			time.Sleep(time.Until(next))
			w.refresh()
		}
	}()
}

func (w *WeeklyLikeScheduler) refresh() {
	logger.Info("scheduler: refreshing weekly likes ranking")
	start := time.Now()

	since := time.Now().AddDate(0, 0, -7)
	realLikes, err := w.likeRepo.AggregateWeeklyLikeCounts(since)
	if err != nil {
		logger.Error("weekly scheduler: failed to aggregate weekly likes", zap.Error(err))
		return
	}

	allScores, err := w.scoreRepo.FindTopByScore(1000)
	if err != nil {
		logger.Error("weekly scheduler: failed to fetch audio scores", zap.Error(err))
		return
	}

	_, weekNum := time.Now().ISOWeek()
	weeklyData := make(map[uint]int64, len(allScores))
	for _, s := range allScores {
		realCount := realLikes[s.AudioID]
		// Deterministic boost per (week, audio): same within a week, differs each week.
		r := rand.New(rand.NewSource(int64(weekNum)*100000 + int64(s.AudioID)))
		boost := r.Int63n(10) // 0-9 robot likes
		weeklyData[s.AudioID] = realCount + boost
	}

	if err := w.scoreRepo.BulkUpdateWeeklyLikes(weeklyData); err != nil {
		logger.Error("weekly scheduler: failed to update weekly likes", zap.Error(err))
		return
	}

	logger.Info("weekly scheduler: weekly likes updated",
		zap.Int("audios", len(weeklyData)),
		zap.Duration("elapsed", time.Since(start)),
	)
}
