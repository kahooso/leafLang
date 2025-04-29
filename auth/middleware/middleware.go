package middleware

import (
	"auth/context"
	"auth/handlers"
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
			handlers.ErrorResponseHandle(c, http.StatusUnauthorized, "Токен не предоставлен")
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
			handlers.ErrorResponseHandle(c, http.StatusUnauthorized, "Неверный токен")
			return
		}

		if !token.Valid {
			log.Println("Токен невалиден:", tokenString)
			handlers.ErrorResponseHandle(c, http.StatusUnauthorized, "Неверный токен")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			handlers.ErrorResponseHandle(c, http.StatusUnauthorized, "Ошибка обработки токена")
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			handlers.ErrorResponseHandle(c, http.StatusUnauthorized, "Неверный формат user_id")
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
