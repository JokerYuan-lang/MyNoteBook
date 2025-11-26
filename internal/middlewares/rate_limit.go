package middlewares

import (
	"context"
	"fmt"
	"time"

	"github.com/JokerYuan-lang/MyNoteBook/pkg/errcode"
	"github.com/JokerYuan-lang/MyNoteBook/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// RateLimit 接口限流（基于 Redis）
// limit: 单位时间内最大请求数
// duration: 时间窗口（如 time.Minute 表示1分钟）
func RateLimit(rdb *redis.Client, limit int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 构造限流 Key（IP + 接口路径）
		key := fmt.Sprintf("rate_limit:%s:%s", c.ClientIP(), c.FullPath())

		// Redis 计数器自增
		ctx := context.Background()
		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			zap.S().Errorf("限流计数器失败: %v", err)
			c.Next()
			return
		}

		// 第一次请求时设置过期时间
		if count == 1 {
			rdb.Expire(ctx, key, duration)
		}

		// 超过限制则拒绝请求
		if count > int64(limit) {
			response.Error(c, errcode.Forbidden, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}
