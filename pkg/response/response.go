package response

import (
	"github.com/JokerYuan-lang/MyNoteBook/pkg/errcode"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"` //业务码
	Message string      `json:"msg"`  //提示信息
	Data    interface{} `json:"data"` //响应数据
}

// Success 成功响应（带数据）
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    errcode.Success,
		Message: errcode.GetMsg(errcode.Success),
		Data:    data,
	})
}

// SuccessWithoutData 成功响应（无数据）
func SuccessWithoutData(c *gin.Context) {
	c.JSON(200, Response{
		Code:    errcode.Success,
		Message: errcode.GetMsg(errcode.Success),
		Data:    nil,
	})
}

// Error 错误响应（自定义提示）
func Error(c *gin.Context, code int, msg string) {
	if msg == "" {
		msg = errcode.GetMsg(code)
	}
	c.JSON(200, Response{
		Code:    code,
		Message: msg,
		Data:    nil,
	})
}

// ErrorWithDefaultMsg 错误响应（使用默认提示）
func ErrorWithDefaultMsg(c *gin.Context, code int) {
	c.JSON(200, Response{
		Code:    code,
		Message: errcode.GetMsg(code),
		Data:    nil,
	})
}
