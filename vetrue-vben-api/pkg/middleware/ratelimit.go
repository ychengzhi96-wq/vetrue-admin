package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 限流器接口
type RateLimiter interface {
	Allow(key string) bool
}

// TokenBucketLimiter 令牌桶限流器
type TokenBucketLimiter struct {
	rate       int // 每秒产生的令牌数
	capacity   int // 桶容量
	tokens     map[string]int
	lastRefill map[string]time.Time
	mu         sync.RWMutex
}

// NewTokenBucketLimiter 创建令牌桶限流器
func NewTokenBucketLimiter(rate, capacity int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		rate:       rate,
		capacity:   capacity,
		tokens:     make(map[string]int),
		lastRefill: make(map[string]time.Time),
	}
}

// Allow 判断是否允许请求
func (tbl *TokenBucketLimiter) Allow(key string) bool {
	tbl.mu.Lock()
	defer tbl.mu.Unlock()

	now := time.Now()

	// 初始化
	if _, exists := tbl.tokens[key]; !exists {
		tbl.tokens[key] = tbl.capacity
		tbl.lastRefill[key] = now
	}

	// 计算应该补充的令牌数
	elapsed := now.Sub(tbl.lastRefill[key])
	tokensToAdd := int(elapsed.Seconds()) * tbl.rate
	if tokensToAdd > 0 {
		tbl.tokens[key] += tokensToAdd
		if tbl.tokens[key] > tbl.capacity {
			tbl.tokens[key] = tbl.capacity
		}
		tbl.lastRefill[key] = now
	}

	// 尝试获取令牌
	if tbl.tokens[key] > 0 {
		tbl.tokens[key]--
		return true
	}

	return false
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Rate     int                       // 每秒请求数
	Capacity int                       // 令牌桶容量
	KeyFunc  func(*gin.Context) string // 生成限流key的函数
}

var globalLimiter *TokenBucketLimiter

// RateLimit 限流中间件
func RateLimit(config RateLimitConfig) gin.HandlerFunc {
	if globalLimiter == nil {
		globalLimiter = NewTokenBucketLimiter(config.Rate, config.Capacity)
	}

	// 默认使用 IP 作为限流 key
	keyFunc := config.KeyFunc
	if keyFunc == nil {
		keyFunc = func(c *gin.Context) string {
			return c.ClientIP()
		}
	}

	// 启动清理协程,定期清理过期的限流记录
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			globalLimiter.mu.Lock()
			now := time.Now()
			for key, lastTime := range globalLimiter.lastRefill {
				if now.Sub(lastTime) > 10*time.Minute {
					delete(globalLimiter.tokens, key)
					delete(globalLimiter.lastRefill, key)
				}
			}
			globalLimiter.mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		key := keyFunc(c)
		if !globalLimiter.Allow(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    http.StatusTooManyRequests,
				"message": "请求过于频繁,请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// IPRateLimit IP限流中间件
func IPRateLimit(rate, capacity int) gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		Rate:     rate,
		Capacity: capacity,
		KeyFunc: func(c *gin.Context) string {
			return c.ClientIP()
		},
	})
}

// UserRateLimit 用户限流中间件(需要在JWT中间件之后使用)
func UserRateLimit(rate, capacity int) gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		Rate:     rate,
		Capacity: capacity,
		KeyFunc: func(c *gin.Context) string {
			// 从上下文中获取用户ID
			if userId, exists := c.Get("userId"); exists {
				return "user:" + userId.(string)
			}
			// 如果没有用户ID,回退到IP限流
			return c.ClientIP()
		},
	})
}
