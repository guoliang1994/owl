package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors(c *gin.Context) {

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,X-Token,X-UserName-Id,Lang,X-Requested-With")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Expose-Headers", "*")
	c.Header("Access-Control-Allow-Credentials", "true")

	// 放行所有OPTIONS方法
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}
	// 处理请求
	c.Next()

}
