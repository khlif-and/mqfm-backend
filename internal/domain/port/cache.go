package port

import (
	"context"
	"time"
)

type CacheTier int

const (
	CacheTierLight  CacheTier = iota
	CacheTierMedium
	CacheTierHeavy
)

func (t CacheTier) TTL() time.Duration {
	switch t {
	case CacheTierLight:
		return 3 * time.Minute
	case CacheTierMedium:
		return 15 * time.Minute
	case CacheTierHeavy:
		return 60 * time.Minute
	default:
		return 5 * time.Minute
	}
}

type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	DeleteByPattern(ctx context.Context, pattern string) error
}

type CacheManager interface {
	SetWithTier(ctx context.Context, key string, value interface{}, tier CacheTier) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	DeleteByPattern(ctx context.Context, pattern string) error
	PrefetchAudioBatch(ctx context.Context, audioIDs []uint) error
	InvalidateAudio(ctx context.Context, audioID uint) error
	WarmUp(ctx context.Context, key string, value interface{}, tier CacheTier) error
}

type TokenStore interface {
	StoreToken(ctx context.Context, userID uint, role string, token string, ttl time.Duration) error
	GetToken(ctx context.Context, userID uint, role string) (string, error)
	DeleteToken(ctx context.Context, userID uint, role string) error
	RefreshToken(ctx context.Context, userID uint, role string) (string, error)
}
