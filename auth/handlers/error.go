package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not found")
)

func ErrorHandler(c *gin.Context, status int, message string) {
	c.HTML(status, "error.html", gin.H{
		"StatusCode": status,
		"StatusText": http.StatusText(status),
		"Message":    message,
	})
	c.Abort()
}

func UnauthorizedHandler(c *gin.Context) {
	ErrorHandler(c, http.StatusUnauthorized, "Please log in to access this page")
}

func NotFoundHandler(c *gin.Context) {
	ErrorHandler(c, http.StatusNotFound, "The page you're looking for doesn't exist")
}

func InternalServerErrorHandler(c *gin.Context) {
	ErrorHandler(c, http.StatusInternalServerError, "Something went wrong on our end. We're working to fix it!")
}

func ForbiddenHandler(c *gin.Context) {
	ErrorHandler(c, http.StatusForbidden, "You don't have permission to access this resource")
}

func BadRequestHandler(c *gin.Context, message string) {
	ErrorHandler(c, http.StatusBadRequest, message)
}
