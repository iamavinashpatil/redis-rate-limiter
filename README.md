# Redis Rate Limiter

A distributed rate limiter implementation in Go using Redis and the **Token Bucket** algorithm. Designed for high-throughput, low-latency scenarios where you need to control request rates across multiple instances.

## Features

- **Token Bucket Algorithm** — Smooth rate limiting with burst support
- **Atomic Operations** — Lua scripting ensures race-condition-free updates
- **Distributed** — Works across multiple application instances via Redis
- **Auto-expiring Keys** — Prevents memory bloat for inactive clients
- **Millisecond Precision** — Accurate sub-second token refill timing
- **Detailed Results** — Returns token counts for debugging and monitoring

## How It Works

The Token Bucket algorithm maintains a bucket of tokens for each client:

1. Bucket starts at full capacity
2. Tokens refill at a constant rate over time
3. Each request consumes one token
4. Requests are denied when the bucket is empty

This allows short bursts while maintaining an average rate limit.

## Installation

```bash
go get github.com/iamavinashpatil/redis-rate-limiter
```

## Prerequisites

- Go 1.23+
- Redis server running (default: `localhost:6379`)

## Usage

```go
package main

import (
    "context"
    "fmt"

    "github.com/iamavinashpatil/redis-rate-limiter/ratelimiter"
    "github.com/redis/go-redis/v9"
)

func main() {
    ctx := context.Background()

    // Connect to Redis
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    // Create rate limiter: 5 tokens capacity, 1 token/second refill rate
    limiter := ratelimiter.NewRedisRateLimiter(redisClient, 5, 1)

    // Check if request is allowed
    clientID := "user-123"
    result, err := limiter.Allow(ctx, clientID)
    if err != nil {
        panic(err)
    }

    if result.Allowed {
        fmt.Printf("Request permitted (%.2f tokens remaining)\n", result.TokensLeft)
    } else {
        fmt.Printf("Rate limit exceeded (only %.2f tokens, need 1)\n", result.TokensBefore)
    }
}
```

## API Reference

### `NewRedisRateLimiter(client, capacity, refillRate)`

Creates a new rate limiter instance.

| Parameter | Type | Description |
|-----------|------|-------------|
| `client` | `*redis.Client` | Redis client connection |
| `capacity` | `int` | Maximum tokens in the bucket (burst size) |
| `refillRate` | `float64` | Tokens added per second |

### `Allow(ctx, clientID) (Result, error)`

Checks if a request should be allowed.

| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | Context for the Redis call |
| `clientID` | `string` | Unique identifier for the client |

### `Result` struct

| Field | Type | Description |
|-------|------|-------------|
| `Allowed` | `bool` | Whether the request was allowed |
| `TokensLeft` | `float64` | Tokens remaining after this request |
| `TokensBefore` | `float64` | Tokens available before consuming |
| `RefillAmount` | `float64` | Tokens refilled since last request |

## Configuration Examples

| Use Case | Capacity | Refill Rate | Effect |
|----------|----------|-------------|--------|
| API endpoint | `100` | `10` | 10 req/sec sustained, 100 req burst |
| Login attempts | `5` | `0.1` | 1 attempt per 10 sec, 5 max burst |
| Message sending | `20` | `1` | 1 msg/sec sustained, 20 msg burst |

## Running the Example

```bash
# Start Redis (if not running)
redis-server

# Run the demo
go run main.go
```

Output:
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Rate Limiter Demo: capacity=5, refill=1 token/sec
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✓ Request  1: ALLOWED
  └─ Had 5.00 tokens → used 1 → 4.00 tokens left
✓ Request  2: ALLOWED
  └─ Had 4.30 tokens → used 1 → 3.30 tokens left
✓ Request  3: ALLOWED
  └─ Had 3.60 tokens → used 1 → 2.60 tokens left
✓ Request  4: ALLOWED
  └─ Had 2.90 tokens → used 1 → 1.90 tokens left
✓ Request  5: ALLOWED
  └─ Had 2.20 tokens → used 1 → 1.20 tokens left
✓ Request  6: ALLOWED
  └─ Had 1.50 tokens → used 1 → 0.50 tokens left
✗ Request  7: REJECTED
  └─ Had 0.80 tokens → need 1 token but only have 0.80!
✓ Request  8: ALLOWED
  └─ Had 1.10 tokens → used 1 → 0.10 tokens left
✗ Request  9: REJECTED
  └─ Had 0.40 tokens → need 1 token but only have 0.40!
✗ Request 10: REJECTED
  └─ Had 0.70 tokens → need 1 token but only have 0.70!
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

> **Note:** With 300ms sleep and 1 token/sec refill rate, ~0.30 tokens are added between each request. Output may vary slightly based on system timing.

## License

MIT
