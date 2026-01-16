package ratelimiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisRateLimiter struct {
	capacity int
	refRate  float64
	client   *redis.Client
	script   *redis.Script
}

// Result contains detailed information about a rate limit check
type Result struct {
	Allowed      bool    // Whether the request was allowed
	TokensLeft   float64 // Tokens remaining after this request
	TokensBefore float64 // Tokens available before consuming
	RefillAmount float64 // Tokens refilled since last request
}

func NewRedisRateLimiter(client *redis.Client, cap int, refRate float64) *redisRateLimiter {
	return &redisRateLimiter{
		capacity: cap,
		refRate:  refRate,
		client:   client,
		script:   redis.NewScript(tokenBucketLua),
	}
}

func (r *redisRateLimiter) Allow(ctx context.Context, clientId string) (Result, error) {
	key := "rate_limit:" + clientId
	nowMs := time.Now().UnixMilli() // Use milliseconds for accurate sub-second timing
	res, err := r.script.Run(
		ctx, r.client, []string{key}, r.capacity, r.refRate, nowMs,
	).Int64Slice()
	if err != nil {
		return Result{}, err
	}

	return Result{
		Allowed:      res[0] == 1,
		TokensLeft:   float64(res[1]) / 100,
		TokensBefore: float64(res[2]) / 100,
		RefillAmount: float64(res[3]) / 100,
	}, nil
}
