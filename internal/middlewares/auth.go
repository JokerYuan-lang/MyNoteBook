package middlewares

import (
	"github.com/JokerYuan-lang/MyNoteBook/internal/config"
	"github.com/JokerYuan-lang/MyNoteBook/pkg/errcode"
	"github.com/JokerYuan-lang/MyNoteBook/pkg/jwt"
	"github.com/JokerYuan-lang/MyNoteBook/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthCheck 登录认证中间件（需要登录的接口使用）
func AuthCheck(conf config.JwtConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 Token（前端传参：Header -> Notebook: token值）
		tokenStr := c.GetHeader("Notebook")
		if tokenStr == "" {
			response.ErrorWithDefaultMsg(c, errcode.Unauthorized)
			c.Abort() // 终止请求
			return
		}

		// 解析 Token
		claims, err := jwt.ParseToken(tokenStr, conf)
		if err != nil {
			zap.S().Errorf("Token 解析失败: %v", err)
			response.Error(c, errcode.Unauthorized, "登录已过期，请重新登录")
			c.Abort()
			return
		}

		// 将用户信息存入上下文（后续接口可直接获取）
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next() // 继续执行后续接口
	}
}
