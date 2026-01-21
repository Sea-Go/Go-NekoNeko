// package errmsg provides business error codes and messages.
// 注意：这里是“业务码”，不是 HTTP status code。
// HTTP 状态码建议在网关/handler 层决定（200/400/401/500...）。

package errmsg

import (
	"errors"
	"fmt"
)

// Code 业务错误码
type Code int32

const (
	// 通用
	CodeOK           Code = 0
	CodeBadRequest   Code = 1000 // 参数/请求不合法
	CodeUnauthorized Code = 1001 // 未登录/鉴权失败
	CodeForbidden    Code = 1002 // 权限不足
	CodeInternal     Code = 9001 // 内部错误(兜底)
	CodeServerBusy   Code = 9000 // 服务繁忙/限流/依赖不可用

	// 用户/登录 (1100~1299)
	CodeUserAlreadyExists       Code = 1100
	CodeUserNotFound            Code = 1101
	CodeUsernameOrPasswordWrong Code = 1102
	CodeUserAlreadyLoggedIn     Code = 1103

	// Token/JWT (1200~1299)
	CodeTokenMissing       Code = 1200
	CodeTokenInvalid       Code = 1201
	CodeTokenExpired       Code = 1202
	CodeTokenRefreshFailed Code = 1203

	// 社区 (2000~2099)
	CodeCommunityAlreadyExists Code = 2000
	CodeCommunityNotFound      Code = 2001

	// 帖子 (3000~3099)
	CodePostAlreadyExists Code = 3000
	CodePostNotFound      Code = 3001

	// 投票 (4000~4099)
	CodeVoteRepeated    Code = 4000
	CodeVoteTimeExpired Code = 4001
)

var msg = map[Code]string{
	CodeOK:           "OK",
	CodeBadRequest:   "请求参数错误",
	CodeUnauthorized: "未登录或登录已失效",
	CodeForbidden:    "权限不足",
	CodeServerBusy:   "服务繁忙，请稍后再试",
	CodeInternal:     "内部错误",

	CodeUserAlreadyExists:       "用户名已存在",
	CodeUserNotFound:            "用户不存在",
	CodeUsernameOrPasswordWrong: "用户名或密码错误",
	CodeUserAlreadyLoggedIn:     "已登录",

	CodeTokenMissing:       "TOKEN缺失",
	CodeTokenInvalid:       "TOKEN无效",
	CodeTokenExpired:       "TOKEN已过期",
	CodeTokenRefreshFailed: "TOKEN刷新失败",

	CodeCommunityAlreadyExists: "社区已存在",
	CodeCommunityNotFound:      "社区不存在",

	CodePostAlreadyExists: "帖子已存在",
	CodePostNotFound:      "帖子不存在",

	CodeVoteRepeated:    "请勿重复投票",
	CodeVoteTimeExpired: "投票时间已过",
}

// BizError 表示业务错误
type BizError struct {
	Code  Code
	Msg   string
	Cause error
}

// Message 根据 Code 返回默认 Msg , Code 对应 Msg 不存在时返回内部错误
func Message(code Code) string {
	if s, ok := msg[code]; ok {
		return s
	}
	return "内部错误"
}

// Error 接口,返回格式化字符串便于输出日志
func (e *BizError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Cause != nil {
		// 给日志用：带 cause；对外回包一般只用 Code/Msg
		return fmt.Sprintf("code=%d msg=%s cause=%v", e.Code, e.Msg, e.Cause)
	}
	return fmt.Sprintf("code=%d msg=%s", e.Code, e.Msg)
}

// Option 用于 New 的可选配置
type Option func(*BizError)

// WithMsg 设置自定义 Msg
func WithMsg(m string) Option {
	return func(e *BizError) {
		e.Msg = m
	}
}

// WithMsgf 设置格式化 Msg
func WithMsgf(format string, args ...any) Option {
	return func(e *BizError) {
		e.Msg = fmt.Sprintf(format, args...)
	}
}

// WithCause 设置底层原因错误 Cause
func WithCause(cause error) Option {
	return func(e *BizError) {
		e.Cause = cause
	}
}

// New 创建 BizError,可选用 WithMsg / WithMsgf, WithCause 自定义消息和原因,不填则使用 Code 对于的默认配置
func New(code Code, opts ...Option) *BizError {
	e := &BizError{Code: code, Msg: Message(code)}
	for _, opt := range opts {
		if opt != nil {
			opt(e)
		}
	}
	return e
}

// Unwrap 返回底层原因错误,用于 errors.Is / errors.As 的错误链解析
func (e *BizError) Unwrap() error {
	return e.Cause
}

// FromError 将任意 error 解析为 业务码,消息 由用内部错误 CodeInternal 兜底
func FromError(err error) (Code, string) {
	if err == nil {
		return CodeOK, Message(CodeOK)
	}
	var be *BizError
	if errors.As(err, &be) {
		return be.Code, be.Msg
	}
	return CodeInternal, Message(CodeInternal)
}

// IsBiz 判断 err 是否为 BizError 或错误链中包含 BizError
func IsBiz(err error) bool {
	var be *BizError
	return errors.As(err, &be)
}

// CodeOf 获取 err 对应的业务码,由内部错误 CodeInternal 兜底
func CodeOf(err error) Code {
	c, _ := FromError(err)
	return c
}
