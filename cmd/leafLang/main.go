package main

import (
	"leaflang/internal/config"
	"leaflang/internal/database"
	"leaflang/internal/handlers"
	"leaflang/internal/routes"
	"leaflang/pkg/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	gin.SetMode(gin.ReleaseMode)

	database.ConnectDB(cfg)

	handlers.InitAuthHandlers(cfg.JWTSecret)
	middleware.SetJWTSecret(cfg.JWTSecret)

	r := gin.Default()
	routes.Set(r)

	log.Printf("Server is running on http://localhost:%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
