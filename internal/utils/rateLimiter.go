package utils

import (
	"sync"
	"time"
)

type RateLimitConfig struct {
	rateLimitKey string
	quota int 
	refillDuration int // in minutes
	createdAt time.Time
	updatedAt time.Time
}

func NewRateLimitConfig(rateLimitKey string, quota int, refillDuration int) *RateLimitConfig {
	return &RateLimitConfig{
		rateLimitKey: rateLimitKey,
		quota: quota,
		refillDuration: refillDuration,
	}
}

type RateLimiter struct {
	rateLimitKey string
	count int
	lastRequestTime time.Time
	lock sync.Mutex
}

// Concurrent access handling
func (r *RateLimiter) acquire() bool {
	r.lock.Lock()
	defer r.lock.Unlock()

	rateLimitConfig := NewRateLimitConfig(r.rateLimitKey, 10, 1)

	// Fill the bucket
	currentTime := time.Now() // 1
	timeElapsed := currentTime.Sub(r.lastRequestTime).Seconds() // 1ms
	tokensToAdd := (time.Duration(rateLimitConfig.quota)/time.Duration(rateLimitConfig.refillDuration)) * time.Duration(timeElapsed)
	r.count += int(tokensToAdd)
	if r.count >= rateLimitConfig.quota {
		r.count = rateLimitConfig.quota
	}

	// 10/1 * 0.001 -> 0 
	// 8 
	// 


	if r.count <= 0 {
		return false // Rate limit exceeded
	}

	r.count--
	r.lastRequestTime = currentTime
	return true;
}

// 1st call 1 ms
// 2nd call 900 ms 

// Config = Shiva, 10, 1
// RateLimiter = Shiva, 0, 
// lastTime = 1
// 900 - 1 ms = 890ms = 0s
// 


// Rate limiter
// FunctionalRequirements
// 1. Unique key - UserID / API Path
// 2. limit is configurable
// 3. 429 - Too many requests

// Entities:
// RateLimitConfig: id(), rateLimitKey, quota, duration, createdAt, updatedAt

// Refill -> Token Bucket
// RateLimit: rateLimitKey, count, lastRequestTime

// Hashmap -> rateLimitKey -> RateLimit