package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/logger"

	"go.uber.org/zap"
)

type cacheManager struct {
	cache    port.CacheRepository
	audioRepo port.AudioRepository
}

func NewCacheManager(cache port.CacheRepository, audioRepo port.AudioRepository) port.CacheManager {
	return &cacheManager{cache: cache, audioRepo: audioRepo}
}

func (m *cacheManager) SetWithTier(ctx context.Context, key string, value interface{}, tier port.CacheTier) error {
	return m.cache.Set(ctx, key, value, tier.TTL())
}

func (m *cacheManager) Get(ctx context.Context, key string) (string, error) {
	return m.cache.Get(ctx, key)
}

func (m *cacheManager) Delete(ctx context.Context, key string) error {
	return m.cache.Delete(ctx, key)
}

func (m *cacheManager) DeleteByPattern(ctx context.Context, pattern string) error {
	return m.cache.DeleteByPattern(ctx, pattern)
}

func (m *cacheManager) PrefetchAudioBatch(ctx context.Context, audioIDs []uint) error {
	if len(audioIDs) == 0 {
		return nil
	}

	limit := 5
	if len(audioIDs) < limit {
		limit = len(audioIDs)
	}
	batch := audioIDs[:limit]

	audios, err := m.audioRepo.FindByIDs(batch)
	if err != nil {
		return err
	}

	for _, audio := range audios {
		key := fmt.Sprintf("audio:%d", audio.ID)
		data, err := json.Marshal(audio)
		if err != nil {
			continue
		}
		if err := m.cache.Set(ctx, key, string(data), port.CacheTierHeavy.TTL()); err != nil {
			logger.Error("prefetch cache set failed", zap.Uint("audio_id", audio.ID), zap.Error(err))
		}
	}

	return nil
}

func (m *cacheManager) InvalidateAudio(ctx context.Context, audioID uint) error {
	key := fmt.Sprintf("audio:%d", audioID)
	return m.cache.Delete(ctx, key)
}

func (m *cacheManager) WarmUp(ctx context.Context, key string, value interface{}, tier port.CacheTier) error {
	_, err := m.cache.Get(ctx, key)
	if err != nil {
		return m.cache.Set(ctx, key, value, tier.TTL())
	}
	return nil
}
