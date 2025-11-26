package main

import (
	"fmt"
	"os"

	"github.com/JokerYuan-lang/MyNoteBook/internal/config"
	"github.com/JokerYuan-lang/MyNoteBook/pkg/db"
	"github.com/JokerYuan-lang/MyNoteBook/pkg/redis"
	"github.com/JokerYuan-lang/MyNoteBook/router"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//项目启动入口

// 全局变量（仅在 main 中初始化，其他地方通过参数传递）
var globalConf config.Config

func main() {
	// 1. 初始化日志（最先初始化，确保其他步骤能打日志）
	initLogger()
	defer zap.L().Sync() // 退出时刷新日志缓冲区

	// 2. 加载配置文件
	if err := loadConfig(); err != nil {
		zap.S().Fatalf("加载配置失败: %v", err)
		os.Exit(1)
	}
	zap.S().Info("加载配置成功")

	// 3. 初始化 MySQL
	mysqlDB, err := db.InitMySQL(globalConf.Mysql, globalConf.Debug)
	if err != nil {
		zap.S().Fatalf("MySQL 初始化失败: %v", err)
		os.Exit(1)
	}

	// 4. 初始化 Redis
	redisClient, err := redis.InitRedis(globalConf.Redis)
	if err != nil {
		zap.S().Fatalf("Redis 初始化失败: %v", err)
		os.Exit(1)
	}

	// 5. 初始化路由
	r := router.InitRouter(mysqlDB, redisClient, globalConf.Jwt, globalConf.Debug)

	// 6. 启动服务
	zap.S().Infof("服务启动成功，监听端口: %d", globalConf.Port)
	if err := r.Run(fmt.Sprintf(":%d", globalConf.Port)); err != nil {
		zap.S().Fatalf("服务启动失败: %v", err)
		os.Exit(1)
	}
}

// initLogger 初始化 Zap 日志
func initLogger() {
	config := zap.NewProductionConfig()
	// 调整日志格式（保留彩色输出，方便本地调试）
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("日志初始化失败: %v", err))
	}
	zap.ReplaceGlobals(logger)
}

// loadConfig 加载配置文件（config.yaml）
func loadConfig() error {
	viper.SetConfigName("config") // 配置文件名（无后缀）
	viper.SetConfigType("yaml")   // 配置文件类型
	viper.AddConfigPath(".")      // 配置文件路径（当前目录）
	viper.AutomaticEnv()          // 支持环境变量覆盖

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// 反序列化到全局配置
	return viper.Unmarshal(&globalConf)
}
