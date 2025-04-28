package middleware

import (
	"auth/context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			log.Println("Ошибка получения токена из cookies:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен не предоставлен"})
			c.Abort()
			return
		}
		fmt.Printf("token: %v\n", tokenString)

		jwtSecret := []byte("clappy")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неверный алгоритм подписи: %v", t.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			log.Println("Ошибка при парсинге токена:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
			c.Abort()
			return
		}

		if !token.Valid {
			log.Println("Токен невалиден:", tokenString)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка обработки токена"})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат user_id"})
			c.Abort()
			return
		}

		user := context.UserContext{
			ID:    uint(userID),
			Email: claims["email"].(string),
			Role:  claims["role"].(string),
		}

		c.Set("user", user)

		c.Next()
	}
}
