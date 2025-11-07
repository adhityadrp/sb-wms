package middlewares

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var JwtKey []byte

func InitJWT() {
	_ = godotenv.Load()
	JwtKey = []byte(os.Getenv("JWT_SECRET"))
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		// allow "Bearer <token>"
		tokenString := strings.TrimSpace(authHeader)
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not valid"})
			c.Abort()
			return
		}

		// check expiration
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
				c.Abort()
				return
			}
		}

		// attach user_id to context
		if uid, ok := claims["user_id"].(float64); ok {
			c.Set("user_id", uint(uid))
		}

		c.Next()
	}
}
