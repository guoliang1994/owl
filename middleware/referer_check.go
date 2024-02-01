package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"owl"
	"strings"
)

// RefererCheck 是一个中间件，用于检查 Referer 头部
func RefererCheck(stage *owl.Stage) gin.HandlerFunc {
	appCfg := stage.ConfManager.GetConfig("app").Get
	return func(c *gin.Context) {
		referer := c.GetHeader("Referer")
		if referer == "" {
			// 如果 Referer 头部为空，表示直接通过浏览器地址栏访问，放行
			c.Next()
			return
		}

		// 检查 Referer 是否包含当前域名
		if strings.Contains(referer, appCfg("domain").ToString()) {
			c.Next()
			return
		}

		// 如果不符合条件，返回 403 Forbidden
		c.AbortWithStatus(http.StatusForbidden)
	}
}
