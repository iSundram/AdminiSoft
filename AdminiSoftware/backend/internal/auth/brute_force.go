
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
