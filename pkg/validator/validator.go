package validator

import (
	"regexp"

	"github.com/JokerYuan-lang/MyNoteBook/pkg/errcode"
	"github.com/go-playground/validator/v10"
)

// 全局校验器实例
var validate = validator.New()

// ValidateStruct 校验结构体（带自定义规则）

func ValidateStruct(obj interface{}) error {
	return validate.Struct(obj)
}

// GetErrorMsg 提取校验错误信息（给前端看）

func GetErrorMsg(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			// 格式：字段名 + 校验规则 + 提示
			return e.Field() + "参数不合法（" + e.Tag() + "）"
		}
	}
	return errcode.GetMsg(errcode.InvalidParam)
}

// CheckPasswordStrength 密码强度校验（8位以上+字母+数字）

func CheckPasswordStrength(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)
	return hasLetter && hasNumber
}
