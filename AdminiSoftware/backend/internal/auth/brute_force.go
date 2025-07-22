
package auth

import (
	"net"
	"sync"
	"time"
)

type BruteForceProtection struct {
	attempts map[string]*AttemptInfo
	mutex    sync.RWMutex
	config   BruteForceConfig
}

type AttemptInfo struct {
	Count      int
	LastAttempt time.Time
	BlockedUntil time.Time
}

type BruteForceConfig struct {
	MaxAttempts   int
	WindowMinutes int
	BlockMinutes  int
}

func NewBruteForceProtection() *BruteForceProtection {
	return &BruteForceProtection{
		attempts: make(map[string]*AttemptInfo),
		config: BruteForceConfig{
			MaxAttempts:   5,
			WindowMinutes: 15,
			BlockMinutes:  30,
		},
	}
}

func (bfp *BruteForceProtection) IsBlocked(ip string) bool {
	bfp.mutex.RLock()
	defer bfp.mutex.RUnlock()
	
	info, exists := bfp.attempts[ip]
	if !exists {
		return false
	}
	
	return time.Now().Before(info.BlockedUntil)
}

func (bfp *BruteForceProtection) RecordAttempt(ip string, successful bool) {
	bfp.mutex.Lock()
	defer bfp.mutex.Unlock()
	
	now := time.Now()
	
	if successful {
		delete(bfp.attempts, ip)
		return
	}
	
	info, exists := bfp.attempts[ip]
	if !exists {
		info = &AttemptInfo{}
		bfp.attempts[ip] = info
	}
	
	windowStart := now.Add(-time.Duration(bfp.config.WindowMinutes) * time.Minute)
	if info.LastAttempt.Before(windowStart) {
		info.Count = 0
	}
	
	info.Count++
	info.LastAttempt = now
	
	if info.Count >= bfp.config.MaxAttempts {
		info.BlockedUntil = now.Add(time.Duration(bfp.config.BlockMinutes) * time.Minute)
	}
}

func (bfp *BruteForceProtection) GetClientIP(remoteAddr string) string {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return ip
}
package auth

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type BruteForceProtection struct {
	redis   *redis.Client
	enabled bool
}

func NewBruteForceProtection(redis *redis.Client) *BruteForceProtection {
	return &BruteForceProtection{
		redis:   redis,
		enabled: redis != nil,
	}
}

func (bf *BruteForceProtection) IsBlocked(ip string) bool {
	if !bf.enabled {
		return false
	}

	ctx := context.Background()
	key := fmt.Sprintf("bf:ip:%s", ip)
	
	attempts, err := bf.redis.Get(ctx, key).Int()
	if err != nil {
		return false
	}

	return attempts >= 5
}

func (bf *BruteForceProtection) RecordAttempt(ip string, success bool) error {
	if !bf.enabled {
		return nil
	}

	ctx := context.Background()
	key := fmt.Sprintf("bf:ip:%s", ip)

	if success {
		// Clear failed attempts on successful login
		return bf.redis.Del(ctx, key).Err()
	}

	// Increment failed attempts
	pipe := bf.redis.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, 15*time.Minute)
	_, err := pipe.Exec(ctx)

	return err
}

func (bf *BruteForceProtection) GetAttempts(ip string) int {
	if !bf.enabled {
		return 0
	}

	ctx := context.Background()
	key := fmt.Sprintf("bf:ip:%s", ip)
	
	attempts, err := bf.redis.Get(ctx, key).Int()
	if err != nil {
		return 0
	}

	return attempts
}

func (bf *BruteForceProtection) GetTimeUntilUnblock(ip string) time.Duration {
	if !bf.enabled {
		return 0
	}

	ctx := context.Background()
	key := fmt.Sprintf("bf:ip:%s", ip)
	
	ttl, err := bf.redis.TTL(ctx, key).Result()
	if err != nil {
		return 0
	}

	return ttl
}

func (bf *BruteForceProtection) UnblockIP(ip string) error {
	if !bf.enabled {
		return nil
	}

	ctx := context.Background()
	key := fmt.Sprintf("bf:ip:%s", ip)
	
	return bf.redis.Del(ctx, key).Err()
}
