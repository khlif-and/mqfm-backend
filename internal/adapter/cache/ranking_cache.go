package cache

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"go.uber.org/zap"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/logger"
)

const (
	rankingCacheKey  = "cache:ranking:monthly"
	popularCacheKey  = "cache:popular:audios"
	rankingCacheTTL  = 5 * time.Hour
	popularCacheTTL  = 7 * time.Hour
	rankingLimit     = 20
	popularLimit     = 15
	rankingMaxLikes  = 35000
	popularMaxLikes  = 10000
)

type RankingCache struct {
	cache    port.CacheRepository
	scoreRepo port.AudioScoreRepository
}

func NewRankingCache(cache port.CacheRepository, scoreRepo port.AudioScoreRepository) *RankingCache {
	return &RankingCache{cache: cache, scoreRepo: scoreRepo}
}

func (rc *RankingCache) GetRanking(ctx context.Context) ([]entity.AudioScore, error) {
	cached, err := rc.cache.Get(ctx, rankingCacheKey)
	if err == nil && cached != "" {
		var scores []entity.AudioScore
		if json.Unmarshal([]byte(cached), &scores) == nil {
			return scores, nil
		}
	}

	return rc.RefreshRanking(ctx)
}

func (rc *RankingCache) RefreshRanking(ctx context.Context) ([]entity.AudioScore, error) {
	scores, err := rc.scoreRepo.FindTopByMonthlyLikes(rankingLimit)
	if err != nil {
		logger.Error("ranking cache refresh failed", zap.Error(err))
		return nil, err
	}

	if len(scores) > 0 {
		_ = rc.cache.Set(ctx, rankingCacheKey, scores, rankingCacheTTL)
	}

	logger.Info("ranking cache refreshed", zap.Int("count", len(scores)))
	return scores, nil
}

func (rc *RankingCache) GetPopular(ctx context.Context) ([]entity.AudioScore, error) {
	cached, err := rc.cache.Get(ctx, popularCacheKey)
	if err == nil && cached != "" {
		var scores []entity.AudioScore
		if json.Unmarshal([]byte(cached), &scores) == nil {
			return scores, nil
		}
	}

	return rc.RefreshPopular(ctx)
}

func (rc *RankingCache) RefreshPopular(ctx context.Context) ([]entity.AudioScore, error) {
	scores, err := rc.scoreRepo.FindTopByLikes(popularLimit*3, popularMaxLikes)
	if err != nil {
		logger.Error("popular cache refresh failed", zap.Error(err))
		return nil, err
	}

	if len(scores) > popularLimit {
		rand.Shuffle(len(scores), func(i, j int) {
			scores[i], scores[j] = scores[j], scores[i]
		})
		scores = scores[:popularLimit]
	}

	if len(scores) > 0 {
		_ = rc.cache.Set(ctx, popularCacheKey, scores, popularCacheTTL)
	}

	logger.Info("popular cache refreshed", zap.Int("count", len(scores)))
	return scores, nil
}
