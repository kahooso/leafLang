package handlers

import (
	database "leaflang/internal/database"
	"leaflang/internal/models"
	"leaflang/internal/models/context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MainHandler(c *gin.Context) {
	user := c.MustGet("user").(context.UserContext)

	var toLearnCount int64
	var knownCount int64
	var learnedCount int64

	database.DB.Model(&models.UserWord{}).
		Joins("JOIN word_status ON user_word.status_id = word_status.id").
		Where("user_word.user_id = ? AND word_status.name = ?", user.ID, "To learn").
		Count(&toLearnCount)

	database.DB.Model(&models.UserWord{}).
		Joins("JOIN word_status ON user_word.status_id = word_status.id").
		Where("user_word.user_id = ? AND word_status.name = ?", user.ID, "Known").
		Count(&knownCount)

	database.DB.Model(&models.UserWord{}).
		Joins("JOIN word_status ON user_word.status_id = word_status.id").
		Where("user_word.user_id = ? AND word_status.name = ?", user.ID, "Learned").
		Count(&learnedCount)

	c.HTML(http.StatusOK, "main.html", gin.H{
		"username":     user.Email,
		"role":         user.Role,
		"toLearnCount": toLearnCount,
		"knownCount":   knownCount,
		"learnedCount": learnedCount,
	})
}

// func ToggleAdminHandler(c *gin.Context) {
// 	userId, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid user ID"})
// 		return
// 	}

// 	var input struct {
// 		Admin bool `json:"admin"`
// 	}
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid input"})
// 		return
// 	}

// 	var user models.User
// 	if err := database.DB.First(&user, userId).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
// 		return
// 	}

// 	var role models.Role
// 	roleName := "Student"

// 	if input.Admin {
// 		roleName = "Admin"
// 	}
// 	if err := database.DB.Where("name = ?", roleName).First(&role).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Role not found"})
// 		return
// 	}

// 	user.RoleID = role.ID
// 	if err := database.DB.Save(&user).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update role"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"success": true})
// }

// func DeleteUserHandler(c *gin.Context) {
// 	userId, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid user ID"})
// 		return
// 	}

// 	if err := database.DB.Delete(&models.User{}, userId).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete user"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"success": true})
// }
