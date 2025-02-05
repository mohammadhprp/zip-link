package services

import (
	"context"
	"time"

	"github.com/mohammadhprp/zip-link/configs"
)

type CacheService struct {
	redisClient *configs.RedisClient
}

func NewCacheService(cache *configs.RedisClient) *CacheService {
	return &CacheService{
		redisClient: cache,
	}
}

// Set stores a key-value pair with optional expiration
func (s *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return s.redisClient.Client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key
func (s *CacheService) Get(ctx context.Context, key string) (string, error) {
	return s.redisClient.Client.Get(ctx, key).Result()
}

// Delete removes a key
func (s *CacheService) Delete(ctx context.Context, key string) error {
	return s.redisClient.Client.Del(ctx, key).Err()
}

// SetHash stores a hash map
func (s *CacheService) SetHash(ctx context.Context, key string, values map[string]interface{}) error {
	return s.redisClient.Client.HSet(ctx, key, values).Err()
}

// GetHash retrieves all fields from a hash
func (s *CacheService) GetHash(ctx context.Context, key string) (map[string]string, error) {
	return s.redisClient.Client.HGetAll(ctx, key).Result()
}
