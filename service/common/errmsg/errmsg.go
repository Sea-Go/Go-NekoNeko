package errmsg

const (
	Success           = 200
	Error             = 500
	CodeServerBusy    = 1015
	ErrorServerCommon = 5001
	ErrorDbUpdate     = 5002
	ErrorDbSelect     = 5003
	ErrInvalidParam   = 2001
)

var codeMsg = map[int]string{
	Success:           "OK",
	Error:             "FAIL",
	CodeServerBusy:    "服务繁忙",
	ErrorServerCommon: "系统内部错误",
	ErrorDbUpdate:     "更新数据库失败",
	ErrorDbSelect:     "查询数据库失败",
	ErrInvalidParam:   "请求参数错误",
}

func GetErrMsg(code int) string {
	msg, ok := codeMsg[code]
	if !ok {
		return codeMsg[Error]
	}
	return msg
}
