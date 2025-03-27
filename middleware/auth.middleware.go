package middleware

import (
	"net/http"
	"os"
	"r2-gallery/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			utils.SendError(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			jwtSecret := []byte(os.Getenv("JWT_SECRET"))
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			utils.LogError("Invalid token", err)
			utils.SendError(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		c.Set("user", token.Claims)
		c.Next()
	}
}
