package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

const (
	// JWT claims 中的 userID 字段名称
	UserIDClaimKey = "user_id"
	// Context 中存储 userID 的 key
	UserIDCtxKey = "user_id_ctx_key"
)

// GetUserIDFromRequest 从 HTTP 请求中提取 userID
// 工作流程：
// 1. 获取 Authorization header
// 2. 解析 JWT token
// 3. 从 claims 中提取 userID
func GetUserIDFromRequest(r *http.Request, secret string) (int64, error) {
	// 步骤 1: 获取 Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("missing authorization header")
	}

	// 步骤 2: 提取 token（格式：Bearer <token>）
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, fmt.Errorf("invalid authorization header format")
	}
	tokenString := parts[1]

	// 步骤 3: 解析 JWT token
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	// 步骤 4: 从 claims 中提取 userID
	userIDInterface, ok := claims[UserIDClaimKey]
	if !ok {
		return 0, fmt.Errorf("user_id not found in token claims")
	}

	// 处理不同的类型（float64 或 string）
	var userID int64
	switch v := userIDInterface.(type) {
	case float64:
		userID = int64(v)
	case string:
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid user_id format: %w", err)
		}
		userID = id
	default:
		return 0, fmt.Errorf("unexpected user_id type: %T", userIDInterface)
	}

	if userID <= 0 {
		return 0, fmt.Errorf("invalid user_id: %d", userID)
	}

	return userID, nil
}

// WriteErrorResponse 写入错误响应
func WriteErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    statusCode,
		"message": errMsg,
	})
}
