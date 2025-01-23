package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/Ashutowwsh/dns-server-go/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient() *RedisClient {

	appConfig := config.LoadConfig()

	redisAddr := appConfig.RedisURL
	redisPassword := appConfig.RedisPassword

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	return &RedisClient{Client: rdb}
}

func (r *RedisClient) GetCache(ctx context.Context, key string) (string, error) {
	result, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Miss
	} else if err != nil {
		return "", fmt.Errorf("redis got error : %w", err)
	}

	return result, nil
}

func (r *RedisClient) SetCache(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.Client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisClient) RateLimit(ctx context.Context, clientIP string, limit int, window time.Duration) (bool, error) {
	key := fmt.Sprintf("rate : %s", clientIP)

	count, err := r.Client.Incr(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis incr error : %w", err)
	}

	if count == 1 {
		err = r.Client.Expire(ctx, key, window).Err()
		if err != nil {
			return false, fmt.Errorf("redis expire error : %w", err)
		}
	}

	return count <= int64(limit), nil
}
