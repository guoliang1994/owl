package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 用https把这个中间件在router里面use一下就好

func RedirectHTTPtoHTTPS() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("X-Forwarded-Proto") == "http" {
			c.Redirect(http.StatusMovedPermanently, "https://"+c.Request.Host+c.Request.RequestURI)
		}
		c.Next()
	}
}
