package router

import (
	"time"

	"github.com/JokerYuan-lang/MyNoteBook/api"
	"github.com/JokerYuan-lang/MyNoteBook/internal/config"
	"github.com/JokerYuan-lang/MyNoteBook/internal/middlewares"
	"github.com/JokerYuan-lang/MyNoteBook/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// InitRouter 初始化路由
func InitRouter(
	db *gorm.DB,
	rdb *redis.Client,
	jwtConf config.JwtConfig,
	debug bool,
) *gin.Engine {
	// 设置 Gin 模式（调试/生产）
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// 1. 全局中间件
	r.Use(middlewares.RequestLogger()) // 日志中间件
	// 跨域中间件（前端对接必备）
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500"}, // 生产环境替换为前端实际域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Notebook"}, // Notebook 是 Token 头
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 2. 初始化服务和 API
	userService := service.NewUserService(db, jwtConf)
	userAPI := api.NewUserAPI(userService)

	noteService := service.NewNoteService(db)
	noteAPI := api.NewNoteAPI(noteService)

	// 3. 路由分组
	apiGroup := r.Group("/api/v1")
	{
		// 公开接口（无需登录）
		publicGroup := apiGroup.Group("/public")
		{
			publicGroup.POST("/register", middlewares.RateLimit(rdb, 5, time.Minute), userAPI.Register) // 注册（1分钟5次）
			publicGroup.POST("/login", middlewares.RateLimit(rdb, 5, time.Minute), userAPI.Login)       // 登录（1分钟5次）
		}

		// 需登录接口（AuthCheck 中间件）
		authGroup := apiGroup.Group("/note")
		authGroup.Use(middlewares.AuthCheck(jwtConf)) // 统一认证
		{
			authGroup.POST("/create", noteAPI.CreateNote)   // 创建笔记
			authGroup.GET("/list", noteAPI.GetNoteList)     // 笔记列表（分页）
			authGroup.GET("/detail", noteAPI.GetNoteByID)   // 笔记详情
			authGroup.PUT("/update", noteAPI.UpdateNote)    // 更新笔记
			authGroup.DELETE("/delete", noteAPI.DeleteNote) // 删除笔记
		}
	}

	return r
}
