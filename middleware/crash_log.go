package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"owl"
	"runtime"
)

// CrashRecover 系统 panic 捕获并记录日志
func CrashRecover(stack bool, l *owl.LoggerFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := make([]byte, 4096)
				length := runtime.Stack(stack, false)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Internal Server Error",
					"error":   fmt.Sprintf("Recovered from panic: %v\n%s", err, stack[:length]),
				})
				l.RuntimeLogger().Error("[PANIC RECOVERED] %s")
			}
		}()

		c.Next()
	}
}
