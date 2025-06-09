package handlers

import (
	"leaflang/internal/database"
	"leaflang/internal/models"
	"leaflang/internal/models/context"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type StatisticsResponse struct {
	TotalWords       int                 `json:"total_words"`
	LearnedWords     int                 `json:"learned_words"`
	KnownWords       int                 `json:"known_words"`
	ToLearnWords     int                 `json:"to_learn_words"`
	SuccessRate      float64             `json:"success_rate"`
	LastActive       time.Time           `json:"last_active"`
	RegistrationDate time.Time           `json:"registration_date"`
	LastPractice     time.Time           `json:"last_practice"`
	SessionsWeek     int                 `json:"sessions-week"`
	WordsWeek        int                 `json:"words_week"`
	DailyActivity    []int               `json:"daily_activity"`
	RecentWords      []RecentWord        `json:"recent_words"`
	User             context.UserContext `json:"user"`
}

type RecentWord struct {
	Word        string `json:"word"`
	Translation string `json:"translation"`
	Status      string `json:"status"`
}

func StatisticsPageHandler(c *gin.Context) {
	userCtx, exists := c.MustGet("user").(context.UserContext)
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.HTML(http.StatusOK, "statistics.html", gin.H{
		"User": userCtx,
	})
}

func StatisticsHandler(c *gin.Context) {
	userCtx, exists := c.MustGet("user").(context.UserContext)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID := userCtx.ID

	var wordStats struct {
		Total      int
		Learned    int
		Known      int
		ToLearn    int
		SuccessSum int
		Attempts   int
	}

	database.DB.Model(&models.UserWord{}).
		Select("COUNT(*) as total, "+
			"SUM(CASE WHEN status_id = (SELECT id FROM word_status WHERE name = 'Learned') THEN 1 ELSE 0 END) as learned, "+
			"SUM(CASE WHEN status_id = (SELECT id FROM word_status WHERE name = 'Known') THEN 1 ELSE 0 END) as known, "+
			"SUM(CASE WHEN status_id = (SELECT id FROM word_status WHERE name = 'To learn') THEN 1 ELSE 0 END) as to_learn, "+
			"SUM(success_count) as success_sum, "+
			"SUM(success_count + fail_count) as attempts").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Scan(&wordStats)

	var user models.User
	database.DB.Select("last_active", "created_at").First(&user, userID)

	var lastPractice time.Time
	database.DB.Model(&models.LearningAttempt{}).
		Select("MAX(attempt_time)").
		Joins("JOIN user_word ON learning_attempt.user_word_id = user_word.id").
		Where("user_word.user_id = ?", userID).
		Scan(&lastPractice)

	var weeklyStats struct {
		Sessions int
		Words    int
	}

	weekAgo := time.Now().AddDate(0, 0, -7)
	database.DB.Model(&models.LearningAttempt{}).
		Select("COUNT(DISTINCT DATE(attempt_time)) as sessions, "+
			"COUNT(DISTINCT user_word_id) as words").
		Joins("JOIN user_word ON learning_attempt.user_word_id = user_word.id").
		Where("user_word.user_id = ? AND attempt_time >= ?", userID, weekAgo).
		Scan(&weeklyStats)

	dailyActivity := make([]int, 7)
	var dailyStats []struct {
		Day   time.Time
		Count int
	}

	database.DB.Model(&models.LearningAttempt{}).
		Select("DATE(attempt_time) as day, COUNT(*) as count").
		Joins("JOIN user_word ON learning_attempt.user_word_id = user_word.id").
		Where("user_word.user_id = ? AND attempt_time >= ?", userID, weekAgo).
		Group("DATE(attempt_time)").
		Scan(&dailyStats)

	for _, stat := range dailyStats {
		daysAgo := int(time.Since(stat.Day).Hours() / 24)
		if daysAgo < 7 {
			dailyActivity[6-daysAgo] = stat.Count
		}
	}

	var recentWords []RecentWord
	database.DB.Model(&models.UserWord{}).
		Select("original_word as word, translation, word_status.name as status").
		Joins("JOIN word_status ON user_word.status_id = word_status.id").
		Where("user_word.user_id = ? AND user_word.last_reviewed IS NOT NULL", userID).
		Order("last_reviewed DESC").
		Limit(5).
		Scan(&recentWords)

	successRate := 0.0
	if wordStats.Attempts > 0 {
		successRate = float64(wordStats.SuccessSum) / float64(wordStats.Attempts) * 100
	}

	c.JSON(http.StatusOK, StatisticsResponse{
		TotalWords:       int(wordStats.Total),
		LearnedWords:     int(wordStats.Learned),
		KnownWords:       int(wordStats.Known),
		ToLearnWords:     int(wordStats.ToLearn),
		SuccessRate:      math.Round(successRate*10) / 10,
		LastActive:       user.LastActive,
		RegistrationDate: user.CreatedAt,
		LastPractice:     lastPractice,
		SessionsWeek:     int(weeklyStats.Sessions),
		WordsWeek:        int(weeklyStats.Words),
		DailyActivity:    dailyActivity,
		RecentWords:      recentWords,
		User:             userCtx,
	})
}
