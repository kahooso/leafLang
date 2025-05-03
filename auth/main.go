package main

import (
	"fmt"
	"log"
	"os"

	database "auth/db"
	"auth/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден")
	}

	database.ConnectDB()

	r := gin.Default()
	routes.Set(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	id := os.Getenv("ID")
	if id == "" {
		id = "localhost"
	}

	fmt.Println("Сервер работает на http://localhost:" + port)
	if err := r.Run(id + ":" + port); err != nil {
		log.Fatal("Не удалось запустить сервер:", err)
	}
}
