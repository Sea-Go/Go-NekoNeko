package errmsg

const (
	Success             = 200
	Error               = 500
	CodeServerBusy      = 1015
	ErrorServerCommon   = 5001
	ErrorDbUpdate       = 5002
	ErrorDbSelect       = 5003
	ErrorUserExist      = 1001
	ErrorLoginWrong     = 1002
	ErrorUserNotExist   = 1003
	ErrorTokenNotExist  = 1004
	ErrorTokenTypeWrong = 1005
	ErrorTokenRuntime   = 1006
	ErrorTokenRefresh   = 1007
	ErrorUserNoRight    = 1008
	ErrorUserNoLogin    = 1009
	ErrorUserLogined    = 1010
	ErrorUserBanned     = 1011
)

var codeMsg = map[int]string{
	Success:             "OK",
	Error:               "FAIL",
	CodeServerBusy:      "服务繁忙",
	ErrorServerCommon:   "系统内部错误",
	ErrorDbUpdate:       "更新数据库失败",
	ErrorDbSelect:       "查询数据库失败",
	ErrorUserExist:      "用户名已存在",
	ErrorLoginWrong:     "用户名或密码错误",
	ErrorUserNotExist:   "用户不存在",
	ErrorTokenNotExist:  "TOKEN不存在",
	ErrorTokenTypeWrong: "TOKEN格式错误",
	ErrorTokenRuntime:   "TOKEN已过期",
	ErrorTokenRefresh:   "TOKEN刷新失败",
	ErrorUserNoRight:    "权限不足",
	ErrorUserNoLogin:    "未登录",
	ErrorUserLogined:    "已登录",
	ErrorUserBanned:     "用户已被封禁",
}

func GetErrMsg(code int) string {
	msg, ok := codeMsg[code]
	if !ok {
		return codeMsg[Error]
	}
	return msg
}
