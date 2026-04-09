package ratelimit

import (
	"log"
	"sync"
	"time"
)

type AttemptRecord struct {
	Count        int
	FirstAttempt time.Time
	LastAttempt  time.Time
	LockedUntil  *time.Time
}

type RateLimiter struct {
	mu       sync.RWMutex
	attempts map[string]*AttemptRecord

	maxAttemptsPerMinute int
	lockoutThreshold     int
	lockoutDuration      time.Duration
	cleanupInterval      time.Duration
}

func NewRateLimiter(maxAttemptsPerMinute, lockoutThreshold int, lockoutDuration, cleanupInterval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		attempts:             make(map[string]*AttemptRecord),
		maxAttemptsPerMinute: maxAttemptsPerMinute,
		lockoutThreshold:     lockoutThreshold,
		lockoutDuration:      lockoutDuration,
		cleanupInterval:      cleanupInterval,
	}
	rl.startCleanupWorker()
	return rl
}

func (rl *RateLimiter) CheckRequest(ip, key string) (bool, time.Duration, string) {
	compositeKey := ip + ":" + key

	rl.mu.RLock()
	record, exists := rl.attempts[compositeKey]
	rl.mu.RUnlock()

	now := time.Now()

	if exists && record.LockedUntil != nil && now.Before(*record.LockedUntil) {
		return false, record.LockedUntil.Sub(now), "locked_out"
	}

	if exists && now.Sub(record.FirstAttempt) < time.Minute {
		if record.Count >= rl.maxAttemptsPerMinute {
			return false, time.Minute - now.Sub(record.FirstAttempt), "rate_limit_exceeded"
		}
	}

	return true, 0, ""
}

func (rl *RateLimiter) RecordFailure(ip, key string) (int, bool) {
	compositeKey := ip + ":" + key

	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	record, exists := rl.attempts[compositeKey]

	if !exists || now.Sub(record.FirstAttempt) > time.Minute {
		record = &AttemptRecord{
			Count:        1,
			FirstAttempt: now,
			LastAttempt:  now,
		}
		rl.attempts[compositeKey] = record
		return 1, false
	}

	record.Count++
	record.LastAttempt = now

	if record.Count >= rl.lockoutThreshold {
		lockUntil := now.Add(rl.lockoutDuration)
		record.LockedUntil = &lockUntil
		log.Printf("[RateLimiter] LOCKOUT: key=%s, attempts=%d", compositeKey, record.Count)
		return record.Count, true
	}

	return record.Count, false
}

func (rl *RateLimiter) RecordSuccess(ip, key string) {
	compositeKey := ip + ":" + key
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.attempts, compositeKey)
}

func (rl *RateLimiter) startCleanupWorker() {
	ticker := time.NewTicker(rl.cleanupInterval)
	go func() {
		for range ticker.C {
			rl.cleanup()
		}
	}()
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for key, record := range rl.attempts {
		if now.Sub(record.LastAttempt) > time.Hour {
			delete(rl.attempts, key)
			continue
		}
		if record.LockedUntil != nil && now.After(*record.LockedUntil) {
			delete(rl.attempts, key)
		}
	}
}
