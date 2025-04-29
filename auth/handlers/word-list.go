package handlers

import (
	"auth/context"
	database "auth/db"
	"auth/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetWordList(c *gin.Context) {
	user := c.MustGet("user").(context.UserContext)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	statusFilter := c.Query("status")
	searchQuery := c.Query("seacrh")

	offset := (page - 1) * limit

	query := database.DB.Model(&models.UserWord{}).
		Where("user_id = ?", user.ID).
		Preload("Status").
		Order("next_review_at ASC")

	if statusFilter != "" {
		query = query.Joins("JOIN word_status ON user_word.status_id = word_status.id").
			Where("word_status.name = ?", statusFilter)
	}

	if searchQuery != "" {
		query = query.Joins("original_word ILIKE ? OR translation ILIKE = ?",
			"%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	var words []models.UserWord
	var total int64

	query.Count(&total)
	query.Offset(offset).Limit(limit).Find(&words)

	var response []gin.H
	for _, word := range words {
		response = append(response, gin.H{
			"id":             word.ID,
			"original_word":  word.OriginalWord,
			"translation":    word.Translation,
			"example":        word.Example,
			"status":         word.Status.Name,
			"success_count":  word.SuccessCount,
			"fail_count":     word.FailCount,
			"next_review_at": word.NextReviewAt.Format("2006-01-02"),
			"last_reviewed":  word.LastReviewed.Format("2006-01-02 15:04"),
		})
	}

	c.HTML(http.StatusOK, "word-list.html", gin.H{
		"success": true,
		"data": gin.H{
			"words": response,
			"pagination": gin.H{
				"total": total,
				"page":  page,
				"limit": limit,
				"pages": (int(total) + limit - 1) / limit,
			},
		},
	})
}

func DeleteWord(c *gin.Context) {
	user := c.MustGet("user").(context.UserContext)
	wordID := c.Param("id")

	if err := database.DB.Where("id = ? AND user_id = ?", wordID, user.ID).Delete(&models.UserWord{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to delete word"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Word deleted!"})
}
