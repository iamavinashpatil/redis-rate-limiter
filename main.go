package main

import (
	"context"
	"fmt"
	"time"

	"github.com/iamavinashpatil/redis-rate-limiter/ratelimiter"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	limiter := ratelimiter.NewRedisRateLimiter(redisClient, 5, 1)
	clientID := "user-123"

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Rate Limiter Demo: capacity=5, refill=1 token/sec")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for i := 1; i <= 10; i++ {
		result, _ := limiter.Allow(ctx, clientID)

		if result.Allowed {
			fmt.Printf("✓ Request %2d: ALLOWED\n", i)
			fmt.Printf("  └─ Had %.2f tokens → used 1 → %.2f tokens left\n",
				result.TokensBefore, result.TokensLeft)
		} else {
			fmt.Printf("✗ Request %2d: REJECTED\n", i)
			fmt.Printf("  └─ Had %.2f tokens → need 1 token but only have %.2f!\n",
				result.TokensBefore, result.TokensBefore)
		}

		time.Sleep(300 * time.Millisecond)
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
