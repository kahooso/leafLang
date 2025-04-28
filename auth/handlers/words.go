package handlers

import (
	database "auth/db"
	"auth/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddWordHandler(c *gin.Context) {
	var input struct {
		OriginalWord string `json:"original_word" binding:"required"`
		Translation  string `json:"translation" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid input"})
		return
	}

	userID := c.MustGet("userID").(uint)

	var status models.WordStatus
	if err := database.DB.Where("name = ?", "To learn").First(&status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Status not found"})
		return
	}

	userWord := models.UserWord{
		UserID:       userID,
		OriginalWord: input.OriginalWord,
		Translation:  input.Translation,
		StatusID:     status.ID,
	}

	if err := database.DB.Create(&userWord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to add word"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "word_id": userWord.ID})
}
