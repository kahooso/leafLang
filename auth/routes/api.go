package routes

import (
	"auth/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/auth/signin", handlers.LoginHandler)
		api.POST("/auth/signup", handlers.RegisterHandler)

		// api.PUT("/users/:id/admin", handlers.ToggleAdminHandler)
		// api.DELETE("/users/:id", handlers.DeleteUserHandler)
	}
}
