package db

import (
	"fmt"
	"time"

	"github.com/JokerYuan-lang/MyNoteBook/internal/config"
	"github.com/JokerYuan-lang/MyNoteBook/internal/model"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitMySQL 初始化 MySQL 连接
func InitMySQL(conf config.MysqlConfig, debug bool) (*gorm.DB, error) {
	// 构造 DSN（数据源名称）
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.UserName,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
	)

	// 配置 GORM 日志级别（调试模式输出 SQL）
	logLevel := logger.Warn
	if debug {
		logLevel = logger.Info
	}

	// 连接 MySQL
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		zap.S().Errorf("MySQL 连接失败: %v", err)
		return nil, err
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		zap.S().Errorf("MySQL 连接池配置失败: %v", err)
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)                  // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)                 // 最大打开连接数
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // 连接最大生命周期

	// 自动迁移数据表（创建/更新表结构）
	err = db.AutoMigrate(
		&model.User{},
		&model.Note{},
		&model.Tag{},
		&model.NoteTag{},
	)
	if err != nil {
		zap.S().Errorf("MySQL 数据表迁移失败: %v", err)
		return nil, err
	}

	zap.S().Info("MySQL 初始化成功")
	return db, nil
}
