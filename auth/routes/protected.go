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
		// main
		protected.GET("/main", handlers.MainHandler)

		// profile
		protected.GET("/profile", handlers.ProfileHandler)
		protected.PUT("/profile/update/:id", handlers.UpdateProfileHandler)

		// word add
		protected.POST("/word/add", handlers.AddWordHandler)

		// word list
		protected.GET("/words", handlers.GetWordList)
		protected.DELETE("/words/:id", handlers.DeleteWord)

		// admin
		// protected.GET("/UserList", handlers.UserListHandler)
	}
}
