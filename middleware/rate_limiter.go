package middleware

import (
	"github.com/gin-gonic/gin"
	"sync"
	"time"
)

// TokenBucket 结构表示令牌桶
type TokenBucket struct {
	capacity       int           // 令牌桶容量
	tokens         int           // 当前令牌数量
	refillInterval time.Duration // 令牌补充间隔
	mu             sync.Mutex    // 互斥锁，用于保护 tokens 的并发访问
}

// NewTokenBucket 创建一个新的令牌桶
func NewTokenBucket(capacity int, refillRate time.Duration) *TokenBucket {
	tb := &TokenBucket{
		capacity:       capacity,
		tokens:         capacity,
		refillInterval: refillRate,
	}
	go tb.startRefill()
	return tb
}

// startRefill 启动补充令牌的 goroutine
func (tb *TokenBucket) startRefill() {
	refillTicker := time.NewTicker(tb.refillInterval)
	defer refillTicker.Stop()
	for {
		<-refillTicker.C
		tb.mu.Lock()
		if tb.tokens < tb.capacity {
			tb.tokens++
		}
		tb.mu.Unlock()
	}
}

// Take 尝试获取一个令牌，成功返回 true，失败返回 false
func (tb *TokenBucket) Take() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

// RateLimiter 是一个 Gin 中间件，用于实现限流
func RateLimiter(capacity int, refillRate time.Duration) gin.HandlerFunc {
	tb := NewTokenBucket(capacity, refillRate)

	return func(c *gin.Context) {
		if tb.Take() {
			c.Next()
		} else {
			c.JSON(429, gin.H{"error": "Too Many Requests"})
			c.Abort()
		}
	}
}
