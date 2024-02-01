package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HeaderAddRequestId(c *gin.Context) {
	// 生成一个新的 UUID
	requestID := uuid.New().String()

	// 将生成的 UUID 添加到请求头中
	c.Header("X-Request-ID", requestID)

	fmt.Println(c.GetHeader("X-Request-ID"))
	// 将 RequestID 存储到 Gin 上下文中，以便后续处理函数可以访问
	c.Set("RequestID", requestID)

	// 执行下一个中间件或处理函数
	c.Next()
}
