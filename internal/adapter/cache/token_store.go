package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/security"
)

type redisTokenStore struct {
	client *redis.Client
}

func NewRedisTokenStore(client *redis.Client) port.TokenStore {
	return &redisTokenStore{client: client}
}

func (s *redisTokenStore) tokenKey(userID uint, role string) string {
	return fmt.Sprintf("jwt:%s:%d", role, userID)
}

func (s *redisTokenStore) StoreToken(ctx context.Context, userID uint, role string, token string, ttl time.Duration) error {
	return s.client.Set(ctx, s.tokenKey(userID, role), token, ttl).Err()
}

func (s *redisTokenStore) GetToken(ctx context.Context, userID uint, role string) (string, error) {
	return s.client.Get(ctx, s.tokenKey(userID, role)).Result()
}

func (s *redisTokenStore) DeleteToken(ctx context.Context, userID uint, role string) error {
	return s.client.Del(ctx, s.tokenKey(userID, role)).Err()
}

func (s *redisTokenStore) RefreshToken(ctx context.Context, userID uint, role string) (string, error) {
	newToken, err := security.GenerateToken(userID, role)
	if err != nil {
		return "", err
	}
	if err := s.StoreToken(ctx, userID, role, newToken, security.TokenTTL); err != nil {
		return "", err
	}
	return newToken, nil
}
