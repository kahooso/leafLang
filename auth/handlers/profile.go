package handlers

import (
	"auth/context"
	database "auth/db"
	"auth/models"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateProfileRequest struct {
	Phone string                `form:"phone"`
	Image *multipart.FileHeader `form:"image"`
}

func ProfileHandler(c *gin.Context) {
	userRaw, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorized!"})
		return
	}

	userCtx := userRaw.(context.UserContext)

	var user models.User
	if err := database.DB.Preload("Role").Where("id = ?", userCtx.ID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось загрузить данные пользователя"})
		return
	}

	imageUrl := ""
	if user.ImageUrl != "" {
		imageUrl = user.ImageUrl
	}

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"id":        user.ID,
		"email":     user.Email,
		"phone":     user.Phone,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"role":      user.Role,
		"image":     imageUrl,
	})
}

func UpdateProfileHandler(c *gin.Context) {
	userRaw, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не Authorized!"})
		return
	}

	userCtx := userRaw.(context.UserContext)

	var req UpdateProfileRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect input"})
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", userCtx.ID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User is not found"})
		return
	}

	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if req.Image != nil {
		if !strings.HasPrefix(req.Image.Header.Get("Content-Type"), "image/") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Можно загружать только изображения"})
			return
		}

		imageDir := "./static/user-images"
		if err := os.MkdirAll(imageDir, 0755); err != nil {
			log.Printf("Ошибка создания папки: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
			return
		}

		ext := filepath.Ext(req.Image.Filename)
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		imagePath := filepath.Join(imageDir, filename)

		if err := c.SaveUploadedFile(req.Image, imagePath); err != nil {
			log.Printf("Ошибка сохранения изображения: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения изображения"})
			return
		}

		if user.ImageUrl != "" {
			oldImage := strings.TrimPrefix(user.ImageUrl, "/static/user-images/")
			if err := os.Remove(filepath.Join(imageDir, oldImage)); err != nil {
				log.Printf("Ошибка удаления старого изображения: %v", err)
			}
		}

		user.ImageUrl = "/static/user-images/" + filename
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения профиля"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Profile is updated",
		"imageUrl": user.ImageUrl,
		"phone":    user.Phone,
	})
}
