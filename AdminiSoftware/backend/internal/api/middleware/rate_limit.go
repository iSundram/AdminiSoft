
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
