package errmsg

const (
	Success = 200
	Error   = 500

	// Follow System
	ErrorServerCommon     = 5001
	ErrorDbWrite          = 6001
	ErrorDbRead           = 6002
	ErrorCannotFollowSelf = 6101
	ErrorCannotBlockSelf  = 6102
)

var codeMsg = map[int]string{
	Success: "OK",
	Error:   "FAIL",

	ErrorServerCommon:     "系统内部错误",
	ErrorDbWrite:          "写入数据库失败",
	ErrorDbRead:           "查询数据库失败",
	ErrorCannotFollowSelf: "不能关注自己",
	ErrorCannotBlockSelf:  "不能拉黑自己",
}

func GetErrMsg(code int) string {
	msg, ok := codeMsg[code]
	if !ok {
		return codeMsg[Error]
	}
	return msg
}
