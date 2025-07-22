
package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
}

type Visitor struct {
	requests int
	lastSeen time.Time
}

func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
	}

	// Cleanup routine
	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Allow(ip string, limit int, window time.Duration) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	visitor, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &Visitor{
			requests: 1,
			lastSeen: time.Now(),
		}
		return true
	}

	// Reset if window has passed
	if time.Since(visitor.lastSeen) > window {
		visitor.requests = 1
		visitor.lastSeen = time.Now()
		return true
	}

	if visitor.requests >= limit {
		return false
	}

	visitor.requests++
	visitor.lastSeen = time.Now()
	return true
}

func RateLimit(requestsPerMinute int) gin.HandlerFunc {
	limiter := NewRateLimiter()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		if !limiter.Allow(ip, requestsPerMinute, time.Minute) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"retry_after": "60 seconds",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func AuthRateLimit() gin.HandlerFunc {
	return RateLimit(5) // 5 requests per minute for auth endpoints
}

func APIRateLimit() gin.HandlerFunc {
	return RateLimit(60) // 60 requests per minute for API endpoints
}

func GeneralRateLimit() gin.HandlerFunc {
	return RateLimit(100) // 100 requests per minute for general endpoints
}
package middleware

import (
	"AdminiSoftware/internal/utils"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type RateLimiter struct {
	redis   *redis.Client
	enabled bool
}

func NewRateLimiter(redis *redis.Client) *RateLimiter {
	return &RateLimiter{
		redis:   redis,
		enabled: redis != nil,
	}
}

func (rl *RateLimiter) Middleware(requests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.enabled {
			c.Next()
			return
		}

		ip := utils.GetClientIP(c.Request.RemoteAddr, c.GetHeader("X-Forwarded-For"), c.GetHeader("X-Real-IP"))
		key := fmt.Sprintf("rate_limit:%s", ip)

		ctx := context.Background()
		current, err := rl.redis.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			c.Next()
			return
		}

		if current >= requests {
			c.Header("X-RateLimit-Limit", strconv.Itoa(requests))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(window).Unix(), 10))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		pipe := rl.redis.Pipeline()
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, window)
		pipe.Exec(ctx)

		remaining := requests - current - 1
		if remaining < 0 {
			remaining = 0
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(requests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(window).Unix(), 10))

		c.Next()
	}
}

func (rl *RateLimiter) AuthRateLimit() gin.HandlerFunc {
	return rl.Middleware(5, 15*time.Minute) // 5 requests per 15 minutes for auth endpoints
}

func (rl *RateLimiter) APIRateLimit() gin.HandlerFunc {
	return rl.Middleware(100, time.Hour) // 100 requests per hour for general API
}

func (rl *RateLimiter) UploadRateLimit() gin.HandlerFunc {
	return rl.Middleware(10, time.Hour) // 10 uploads per hour
}
