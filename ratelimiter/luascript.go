package ratelimiter

// Lua script for atomic token bucket
// Returns: {allowed, tokensRemaining * 100, tokensBefore * 100, refillAmount * 100}
// (multiplied by 100 to preserve 2 decimal places as integers)
// Now uses milliseconds for accurate sub-second timing
const tokenBucketLua = `
local capacity = tonumber(ARGV[1])
local refillRate = tonumber(ARGV[2])
local nowMs = tonumber(ARGV[3])

local data = redis.call("HMGET", KEYS[1], "tokens", "last_refill")
local tokens = tonumber(data[1])
local lastRefillMs = tonumber(data[2])

local refillAmount = 0
if tokens == nil then
	tokens = capacity
	lastRefillMs = nowMs
else
	local elapsedMs = nowMs - lastRefillMs
	local elapsedSec = elapsedMs / 1000
	refillAmount = elapsedSec * refillRate
	tokens = math.min(capacity, tokens + refillAmount)
end

local tokensBefore = tokens

local allowed = 0
if tokens >= 1 then
	tokens = tokens - 1
	allowed = 1
end

redis.call("HMSET", KEYS[1], "tokens", tokens, "last_refill", nowMs)
redis.call("EXPIRE", KEYS[1], math.ceil(capacity / refillRate * 2))

return {allowed, math.floor(tokens * 100), math.floor(tokensBefore * 100), math.floor(refillAmount * 100)}
`
