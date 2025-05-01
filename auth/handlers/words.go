package handlers

import (
	"auth/context"
	database "auth/db"
	"auth/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AddWordHandler(c *gin.Context) {
	var input struct {
		OriginalWord string `json:"original_word" binding:"required,min=1,max=100"`
		Translation  string `json:"translation" binding:"required,min=1,max=100"`
		Example      string `json:"example" max:"500"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Incorrect data",
			"details": err.Error(),
		})
		return
	}

	user := c.MustGet("user").(context.UserContext)

	word := models.UserWord{
		UserID:       user.ID,
		OriginalWord: input.OriginalWord,
		Translation:  input.Translation,
		Example:      input.Example,
		StatusID:     1,
		LastReviewed: time.Now(),
		NextReviewAt: time.Now().Add(6 * time.Hour),
	}

	if err := database.DB.Create(&word).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error while saving word",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Word has successfully added",
		"word_id": word.ID,
	})
}
