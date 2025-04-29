package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ErrorResponseHandle(c *gin.Context, status int, customMessage string) {
	defaultMessages := map[int]string{
		http.StatusBadRequest:          "Bad request",
		http.StatusUnauthorized:        "Unauthorized",
		http.StatusForbidden:           "Access denied",
		http.StatusNotFound:            "Page is not found",
		http.StatusInternalServerError: "Internal server error",
	}

	isAPIRequest := strings.Contains(c.GetHeader("Accept"), "application/json") ||
		strings.HasPrefix(c.Request.URL.Path, "/api/")

	message := customMessage
	if message == "" {
		message = defaultMessages[status]
	}

	if isAPIRequest {
		c.JSON(status, gin.H{
			"error":   http.StatusText(status),
			"message": message,
		})
	} else {
		c.HTML(status, "error.html", gin.H{
			"StatusCode": status,
			"StatusText": http.StatusText(status),
			"Message":    message,
		})
	}

	c.Abort()
}

func AbortWithError(c *gin.Context, status int, message string) {
	c.Status(status)
	ErrorResponseHandle(c, status, message)
}

func AbortUnauthorized(c *gin.Context) {
	AbortWithError(c, http.StatusUnauthorized, "")
}

func AbortForbidden(c *gin.Context) {
	AbortWithError(c, http.StatusForbidden, "")
}

func AbortNotFound(c *gin.Context) {
	AbortWithError(c, http.StatusNotFound, "")
}

func AbortInternalError(c *gin.Context) {
	AbortWithError(c, http.StatusInternalServerError, "")
}
