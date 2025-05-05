package handlers

import (
	database "leaflang/internal/database"
	"leaflang/internal/models"
	"leaflang/internal/models/context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetWordList(c *gin.Context) {
	user := c.MustGet("user").(context.UserContext)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	statusFilter := c.Query("status")
	searchQuery := c.Query("search")

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
		query = query.Where("original_word ILIKE ? OR translation ILIKE ?",
			"%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	var words []models.UserWord
	var total int64

	query.Count(&total)
	query.Offset(offset).Limit(limit).Find(&words)

	pages := (int(total) + limit - 1) / limit

	c.HTML(http.StatusOK, "word-list.html", gin.H{
		"Words":        words,
		"StatusFilter": statusFilter,
		"SearchQuery":  searchQuery,
		"Pagination": gin.H{
			"Total": total,
			"Page":  page,
			"Limit": limit,
			"Pages": pages,
		},
	})
}

func DeleteWord(c *gin.Context) {
	user := c.MustGet("user").(context.UserContext)
	wordID := c.Param("id")

	if err := database.DB.Where("id = ? AND user_id = ?", wordID, user.ID).Delete(&models.UserWord{}).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": "Failed to delete word",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
