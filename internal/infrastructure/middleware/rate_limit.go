package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/response"
)

type visitor struct {
	tokens    int
	lastSeen  time.Time
}

type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     int
	burst    int
	window   time.Duration
}

func newRateLimiter(rate int, burst int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
		window:   window,
	}

	go rl.cleanup()
	return rl
}

func (rl *rateLimiter) cleanup() {
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

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	now := time.Now()

	if !exists {
		rl.visitors[ip] = &visitor{tokens: rl.burst - 1, lastSeen: now}
		return true
	}

	elapsed := now.Sub(v.lastSeen)
	v.lastSeen = now

	refill := int(elapsed / rl.window) * rl.rate
	v.tokens += refill
	if v.tokens > rl.burst {
		v.tokens = rl.burst
	}

	if v.tokens <= 0 {
		return false
	}

	v.tokens--
	return true
}

func RateLimit(requestsPerSecond int, burst int) gin.HandlerFunc {
	limiter := newRateLimiter(requestsPerSecond, burst, time.Second)

	return func(c *gin.Context) {
		if !limiter.allow(c.ClientIP()) {
			response.Error(c, http.StatusTooManyRequests, constant.MsgRateLimitExceeded, nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
