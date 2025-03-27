package services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint, role string) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 测试
	fmt.Println("jwtSecret:", jwtSecret)
	fmt.Println("JWT_SECRET环境变量:", os.Getenv("JWT_SECRET"))
	fmt.Println("JWT_SECRET字节数组:", []byte(os.Getenv("JWT_SECRET")))

	return token.SignedString(jwtSecret)
}
