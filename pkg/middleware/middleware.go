package middleware

import (
	"fmt"
	"leaflang/internal/handlers"
	"leaflang/internal/models/context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	const op string = "pkg.middleware.TokenAuthMiddleware"
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			log.Printf("%s: %s\n", op, err)
			handlers.ErrorResponseHandle(c, http.StatusUnauthorized, "Non authorized")
			return
		}
		fmt.Printf("token: %v\n", tokenString)

		jwtSecret := []byte("clappy")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("%s: %s", op, t.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			log.Printf("%s\tToken parse error: %s\t", op, err)
			handlers.ErrorResponseHandle(c, http.StatusUnauthorized, "Incorrect token")
			return
		}

		if !token.Valid {
			log.Printf("%s\tToken is invalid: %s\n", op, tokenString)
			handlers.ErrorResponseHandle(c, http.StatusUnauthorized, "Incorrect token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			handlers.ErrorResponseHandle(c, http.StatusUnauthorized, "Token validation error")
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			handlers.ErrorResponseHandle(c, http.StatusUnauthorized, "Incorrect format user_id")
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
