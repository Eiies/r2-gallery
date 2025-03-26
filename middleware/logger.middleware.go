package middleware

import (
	"r2-gallery/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// 记录 API 请求日志
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		duration := time.Since(startTime)

		utils.Info("Method: %s, Path: %s, Status: %d, Duration: %v",
			c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration)
	}
}
