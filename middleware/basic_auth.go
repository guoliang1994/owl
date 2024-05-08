package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// BasicAuth 中间件
func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, password, ok := c.Request.BasicAuth()

		// 检查用户名和密码是否正确
		if !ok || user != "9918" || password != "mbox@99#18" {
			c.Header("WWW-Authenticate", "Basic realm=Restricted")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Next()
	}
}
