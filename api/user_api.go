package api

import (
	"github.com/JokerYuan-lang/MyNoteBook/internal/service"
	"github.com/JokerYuan-lang/MyNoteBook/pkg/errcode"
	"github.com/JokerYuan-lang/MyNoteBook/pkg/response"
	"github.com/JokerYuan-lang/MyNoteBook/pkg/validator"
	"github.com/gin-gonic/gin"
)

// 注册请求参数

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"` // 用户名3-20位
	Password string `json:"password" binding:"required,min=8"`        // 密码至少8位
	Email    string `json:"email" binding:"required,email"`           // 邮箱格式
}

// 登录请求参数

type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

// UserAPI 用户接口
type UserAPI struct {
	userService *service.UserService
}

// NewUserAPI 创建 UserAPI 实例
func NewUserAPI(userService *service.UserService) *UserAPI {
	return &UserAPI{userService: userService}
}

// Register 用户注册接口
func (a *UserAPI) Register(c *gin.Context) {
	var req RegisterRequest
	// 绑定并校验参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.InvalidParam, validator.GetErrorMsg(err))
		return
	}

	// 调用业务逻辑
	err := a.userService.Register(req.Username, req.Password, req.Email)
	if err != nil {
		response.Error(c, errcode.DuplicateData, err.Error())
		return
	}

	// 返回成功
	response.SuccessWithoutData(c)
}

// Login 用户登录接口
func (a *UserAPI) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.InvalidParam, validator.GetErrorMsg(err))
		return
	}

	// 调用业务逻辑生成 Token
	token, err := a.userService.Login(req.Username, req.Password)
	if err != nil {
		response.Error(c, errcode.PasswordError, err.Error())
		return
	}

	// 返回 Token
	response.Success(c, gin.H{"token": token})
}
