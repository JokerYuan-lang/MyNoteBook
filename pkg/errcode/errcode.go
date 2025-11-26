package errcode

// 业务错误码（前端根据 code 判断状态）
const (
	Success       = 200 // 成功
	InvalidParam  = 400 // 参数错误
	Unauthorized  = 401 // 未登录（token 无效/缺失）
	Forbidden     = 403 // 无权限
	NotFound      = 404 // 资源不存在（如笔记不存在）
	ServerError   = 500 // 服务器内部错误
	DuplicateData = 601 // 数据重复（如账号已注册）
	PasswordError = 602 // 密码错误
)

// 错误码对应提示信息

func GetMsg(code int) string {
	switch code {
	case Success:
		return "操作成功"
	case InvalidParam:
		return "参数无效"
	case Unauthorized:
		return "请先登录"
	case Forbidden:
		return "无权限操作"
	case NotFound:
		return "资源不存在"
	case ServerError:
		return "服务器内部错误"
	case DuplicateData:
		return "数据已存在"
	case PasswordError:
		return "密码错误"
	default:
		return "未知错误"
	}
}
