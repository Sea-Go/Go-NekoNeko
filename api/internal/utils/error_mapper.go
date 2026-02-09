package utils

import (
	"errors"
	"net/http"

	"sea-try-go/service/favorite/favorite_item"
)

// ErrorResponse API 错误响应结构
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MapBusinessErrorToHTTPStatus 将业务错误映射到 HTTP 状态码和错误信息
// 负责将应用层错误转换为 RESTful HTTP 响应
func MapBusinessErrorToHTTPStatus(err error) (int, string) {
	if err == nil {
		return http.StatusOK, "success"
	}

	switch {
	// 404 Not Found 错误
	case errors.Is(err, favorite_item.ErrFolderNotFound):
		return http.StatusNotFound, "收藏夹不存在"

	case errors.Is(err, favorite_item.ErrItemNotFound):
		return http.StatusNotFound, "收藏项不存在"

	// 403 Forbidden 错误（权限相关）
	case errors.Is(err, favorite_item.ErrFolderNotOwned):
		return http.StatusForbidden, "收藏夹不属于当前用户"

	case errors.Is(err, favorite_item.ErrPermissionDenied):
		return http.StatusForbidden, "权限不足"

	// 409 Conflict 错误（重复资源等）
	case errors.Is(err, favorite_item.ErrItemAlreadyExists):
		return http.StatusConflict, "该对象已被收藏"

	case errors.Is(err, favorite_item.ErrDuplicateFavorite):
		return http.StatusConflict, "不能重复收藏同一对象"

	// 400 Bad Request 错误（参数验证等）
	// 这里可以处理其他类型的错误
	default:
		// 其他未知错误默认返回 500
		return http.StatusInternalServerError, "服务器错误，请稍后重试"
	}
}

// WriteBusinessErrorResponse 根据业务错误写入 HTTP 错误响应
func WriteBusinessErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	statusCode, message := MapBusinessErrorToHTTPStatus(err)
	WriteErrorResponse(w, r, statusCode, message)
}

// ErrorResponseWrapper 包装标准错误响应
func ErrorResponseWrapper(code int, message string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
	}
}
