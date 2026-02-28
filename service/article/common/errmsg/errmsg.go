package errmsg

const (
	Success           = 200
	Error             = 500
	CodeServerBusy    = 1015
	ErrorServerCommon = 5001
	ErrorDbUpdate     = 5002
	ErrorDbSelect     = 5003
	ErrorMinioUpload  = 5004
	ErrorMinioDelete  = 5005
	ErrorMinioDownload = 5006
	ErrorArticleExist = 1001
	ErrorArticleNone  = 1002
)

var codeMsg = map[int]string{
	Success:           "OK",
	Error:             "FAIL",
	CodeServerBusy:    "服务繁忙",
	ErrorServerCommon: "系统内部错误",
	ErrorDbUpdate:     "更新数据库失败",
	ErrorDbSelect:     "查询数据库失败",
	ErrorMinioUpload:  "文件上传失败",
	ErrorMinioDelete:  "文件删除失败",
	ErrorMinioDownload: "文件下载失败",
	ErrorArticleExist: "文章已存在",
	ErrorArticleNone:  "文章不存在",
}

func GetErrMsg(code int) string {
	msg, ok := codeMsg[code]
	if !ok {
		return codeMsg[Error]
	}
	return msg
}
