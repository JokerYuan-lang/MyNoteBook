package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestLogger 记录请求日志（耗时、路径、方法等）
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()

		// 执行后续接口
		c.Next()

		// 记录结束时间，计算耗时
		costTime := time.Since(startTime)

		// 打印日志
		zap.S().Infof(
			"请求信息 - 方法: %s, 路径: %s, 状态码: %d, 耗时: %v, IP: %s",
			c.Request.Method,
			c.FullPath(),
			c.Writer.Status(),
			costTime,
			c.ClientIP(),
		)
	}
}
