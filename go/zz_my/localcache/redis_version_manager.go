package localcache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var _ VersionManager = (*RedisVersionManager)(nil)

type RedisVersionManager struct {
	client     *redis.Client
	prefixKey  string
	expiration time.Duration
}

func NewRedisVersionManager(client *redis.Client, prefixKey string) *RedisVersionManager {
	return &RedisVersionManager{client: client, prefixKey: prefixKey, expiration: 1 * time.Hour}
}

func (r *RedisVersionManager) key(k string) string {
	return fmt.Sprintf("LocalCacheRedisKey:%s:%s", r.prefixKey, k)
}
func (r *RedisVersionManager) Version(ctx context.Context, key string) (string, error) {
	if r.client.Expire(ctx, r.key(key), r.expiration).Val() {
		v := r.client.Get(ctx, r.key(key))
		return v.String(), v.Err()
	} else {
		return r.IncrVersion(ctx, key)
	}
}
func (r *RedisVersionManager) IncrVersion(ctx context.Context, key string) (string, error) {
	var result string
	if _, err := r.client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		intCmd := pipe.Incr(ctx, r.key(key))
		result = fmt.Sprint(intCmd.Val())
		pipe.Expire(ctx, r.key(key), r.expiration)
		return nil
	}); err != nil {
		return "", err
	}

	return result, nil
}
