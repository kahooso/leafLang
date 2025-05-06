package routes

import (
	"leaflang/internal/handlers"
	"leaflang/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterProtectedRoutes(r *gin.Engine) {
	protected := r.Group("/")
	protected.Use(middleware.TokenAuthMiddleware())
	{
		protected.GET("/main", handlers.MainHandler)

		protected.GET("/profile", handlers.ProfileHandler)
		protected.PUT("/profile/update/:id", handlers.UpdateProfileHandler)

		protected.POST("/word/add", handlers.AddWordHandler)

		protected.GET("/words", handlers.GetWordList)
		protected.DELETE("/words/:id", handlers.DeleteWord)
		protected.POST("/words/:id", handlers.UpdateWord)

		// admin
		// protected.GET("/UserList", handlers.UserListHandler)
	}
}
