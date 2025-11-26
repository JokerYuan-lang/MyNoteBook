package service

import (
	"errors"
	"fmt"

	"github.com/JokerYuan-lang/MyNoteBook/internal/config"
	"github.com/JokerYuan-lang/MyNoteBook/internal/model"
	"github.com/JokerYuan-lang/MyNoteBook/pkg/errcode"
	"github.com/JokerYuan-lang/MyNoteBook/pkg/jwt"
	"github.com/JokerYuan-lang/MyNoteBook/pkg/validator"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// UserService 用户业务逻辑
type UserService struct {
	db      *gorm.DB
	jwtConf config.JwtConfig
}

// NewUserService 创建 UserService 实例
func NewUserService(db *gorm.DB, jwtConf config.JwtConfig) *UserService {
	return &UserService{db: db, jwtConf: jwtConf}
}

// Register 用户注册
func (s *UserService) Register(username, password, email string) error {
	// 1. 校验密码强度
	if !validator.CheckPasswordStrength(password) {
		return fmt.Errorf("密码强度不足（需8位以上，包含字母和数字）")
	}

	// 2. 检查用户名/邮箱是否已存在
	var existUser model.User
	err := s.db.Where("username = ?", username).First(&existUser).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		zap.S().Errorf("查询用户名失败: %v", err)
		return fmt.Errorf(errcode.GetMsg(errcode.ServerError))
	}
	if existUser.ID > 0 {
		return fmt.Errorf(errcode.GetMsg(errcode.DuplicateData))
	}

	err = s.db.Where("email = ?", email).First(&existUser).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		zap.S().Errorf("查询邮箱失败: %v", err)
		return fmt.Errorf(errcode.GetMsg(errcode.ServerError))
	}
	if existUser.ID > 0 {
		return fmt.Errorf(errcode.GetMsg(errcode.DuplicateData))
	}

	// 3. 创建用户（密码会在 BeforeSave 钩子中自动加密）
	user := model.User{
		Username: username,
		Password: password,
		Email:    email,
	}
	if err := s.db.Create(&user).Error; err != nil {
		zap.S().Errorf("创建用户失败: %v", err)
		return fmt.Errorf(errcode.GetMsg(errcode.ServerError))
	}

	return nil
}

// Login 用户登录（返回 Token）
func (s *UserService) Login(username, password string) (string, error) {
	// 1. 查询用户
	var user model.User
	err := s.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf(errcode.GetMsg(errcode.PasswordError))
		}
		zap.S().Errorf("查询用户失败: %v", err)
		return "", fmt.Errorf(errcode.GetMsg(errcode.ServerError))
	}

	// 2. 验证密码
	if !user.CheckPassword(password) {
		return "", fmt.Errorf(errcode.GetMsg(errcode.PasswordError))
	}

	// 3. 生成 JWT Token
	token, err := jwt.GenerateToken(user.ID, user.Username, s.jwtConf)
	if err != nil {
		zap.S().Errorf("生成 Token 失败: %v", err)
		return "", fmt.Errorf(errcode.GetMsg(errcode.ServerError))
	}

	return token, nil
}
