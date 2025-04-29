package middleware

import (
	"auth/handlers"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Written() {
			return
		}

		status := c.Writer.Status()
		if status >= 400 {
			c.Abort()
			handlers.ErrorResponseHandle(c, status, "")
		}
	}
}
