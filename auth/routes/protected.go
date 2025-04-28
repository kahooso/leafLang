package routes

import (
	"auth/handlers"
	"auth/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterProtectedRoutes(r *gin.Engine) {
	protected := r.Group("/")
	protected.Use(middleware.TokenAuthMiddleware())
	{
		protected.GET("/main", handlers.MainHandler)
		protected.GET("/profile", handlers.ProfileHandler)
		protected.PUT("/profile/update/:id", handlers.UpdateProfileHandler)
		protected.POST("/word/add")

		// admin
		// protected.GET("/UserList", handlers.UserListHandler)
	}
}
