package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 用于生成测试 JWT token 的工具
// 使用方式：go run jwt_generator.go
func main() {
	// 需要与 api/etc/favorite.yaml 中 UserAuth.AccessSecret 匹配
	secret := "favorite-secret-key"

	// 创建 Claims
	claims := jwt.MapClaims{
		"user_id": int64(1), // 测试用户 ID
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}

	// 生成 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatalf("Failed to sign token: %v", err)
	}

	fmt.Println("JWT Token (用于 Authorization header 中):")
	fmt.Printf("Bearer %s\n\n", tokenString)

	// 解码验证
	fmt.Println("Token Claims:")
	fmt.Printf("  user_id: %d\n", claims["user_id"])
	fmt.Printf("  exp: %d\n", claims["exp"])
	fmt.Printf("  iat: %d\n", claims["iat"])
}
