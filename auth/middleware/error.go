package middleware

import (
	"auth/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Продолжаем выполнение цепочки middleware и обработчиков

		// Проверяем статус ответа и имеющиеся ошибки
		status := c.Writer.Status()

		// Если есть ошибки в контексте
		if len(c.Errors) > 0 {
			// Берем последнюю ошибку
			lastError := c.Errors.Last()

			switch lastError.Err {
			case handlers.ErrUnauthorized:
				handlers.UnauthorizedHandler(c)
				return
			case handlers.ErrForbidden:
				handlers.ForbiddenHandler(c)
				return
			case handlers.ErrUnauthorized:
				handlers.UnauthorizedHandler(c)
				return
			}

			if status >= http.StatusInternalServerError {
				handlers.InternalServerErrorHandler(c)
				return
			}
		}

		switch status {
		case http.StatusNotFound:
			handlers.NotFoundHandler(c)
		case http.StatusUnauthorized:
			handlers.UnauthorizedHandler(c)
		case http.StatusForbidden:
			handlers.ForbiddenHandler(c)
		case http.StatusInternalServerError:
			handlers.InternalServerErrorHandler(c)
		}
	}
}
