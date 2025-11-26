package redis

import (
	"context"
	"fmt"

	"github.com/JokerYuan-lang/MyNoteBook/internal/config"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func InitRedis(conf config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.Database,
		PoolSize: 100, // 连接池大小
	})

	// 测试连接
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		zap.S().Errorf("Redis 连接失败: %v", err)
		return nil, err
	}

	zap.S().Info("Redis 初始化成功")
	return client, nil
}
