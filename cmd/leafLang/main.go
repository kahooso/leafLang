package main

import (
	"fmt"
	"leaflang/internal/database"
	"leaflang/internal/routes"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	if err := godotenv.Load(filepath.Join("../../", ".env")); err != nil {
		log.Println("File .env is not found")
	}
	database.ConnectDB()

	r := gin.Default()
	routes.Set(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server is working on http://localhost:" + port)
	if err := r.Run("localhost:" + port); err != nil {
		log.Fatal("Error while starting server:", err)
	}
}
