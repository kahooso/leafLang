package handlers

import (
	"leaflang/internal/database"
	"leaflang/internal/models"
	"leaflang/internal/models/context"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PracticeWord struct {
	ID           uint   `json:"id"`
	OriginalWord string `json:"original_word"`
	Translation  string `json:"translation"`
	Example      string `json:"example"`
	QuestionType string `json:"question_type"`
	Question     string `json:"question"`
	Answer       string `json:"answer"`
}

type PracticeResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message,omitempty"`
	Words   []PracticeWord `json:"words,omitempty"`
}

func PracticePageHandler(c *gin.Context) {
	user := c.MustGet("user").(context.UserContext)

	var wordsCount int64
	database.DB.Model(&models.UserWord{}).
		Where("user_id = ? AND next_review_at <= ?", user.ID, time.Now()).
		Count(&wordsCount)

	c.HTML(http.StatusOK, "practice.html", gin.H{
		"Title":      "Word Practice",
		"WordsCount": wordsCount,
	})
}

func GetPracticeWordHandler(c *gin.Context) {
	user := c.MustGet("user").(context.UserContext)
	const wordsList int = 10

	var words []models.UserWord
	result := database.DB.Where("user_id = ? AND next_review_at <= ?", user.ID, time.Now()).
		Order("next_review_at ASC").
		Limit(wordsList).
		Find(&words)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch words for practice",
		})
		return
	}

	if len(words) == 0 {
		c.JSON(http.StatusOK, PracticeResponse{
			Success: true,
			Message: "No words to practice right now",
			Words:   []PracticeWord{},
		})
		return
	}

	PracticeWords := make([]PracticeWord, 0, len(words))
	for _, word := range words {
		if word.OriginalWord == "" || word.Translation == "" {
			continue
		}

		pw := createPracticeWord(word)
		PracticeWords = append(PracticeWords, pw)
	}

	c.JSON(http.StatusOK, PracticeResponse{
		Success: true,
		Words:   PracticeWords,
	})
}

func SubmitPracticeResultHandler(c *gin.Context) {
	user := c.MustGet("user").(context.UserContext)

	var input struct {
		Answers []struct {
			WordID  uint `json:"word_id" binding:"required"`
			Success bool `json:"success" binding:"required"`
		} `json:"answers" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid input data",
		})
		return
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, answer := range input.Answers {
		var count int64
		if err := tx.Model(&models.UserWord{}).
			Where("id = ? AND user_id = ?", answer.WordID, user.ID).
			Count(&count).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Word doesn't belong to user",
			})
			return
		}
	}

	for _, answer := range input.Answers {
		var word models.UserWord
		if err := tx.Where("id = ? AND user_id = ?", answer.WordID, user.ID).First(&word).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Word not found",
			})
			return
		}

		updates := map[string]interface{}{}
		if answer.Success {
			updates["success_count"] = word.SuccessCount + 1
		} else {
			updates["fail_count"] = word.FailCount + 1
		}

		newStatusID := determineWordStatus(word.SuccessCount, word.FailCount, answer.Success)

		if newStatusID != word.StatusID {
			updates["status_id"] = newStatusID
		}

		nextReview := calculateNextReview(word.SuccessCount, word.FailCount, answer.Success)

		updates["next_review_at"] = nextReview
		updates["last_reviewed"] = time.Now()

		if err := tx.Model(&word).Updates(updates).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to update word",
			})
		}

		attempt := models.LearningAttempt{
			UserWordID: answer.WordID,
			Success:    answer.Success,
		}

		if err := tx.Create(&attempt).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to save attempt",
			})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to save results",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Practice results saved successfully",
	})
}

func determineWordStatus(successCount, failCount int, lastSuccess bool) uint {
	totalAttempts := successCount + failCount
	if totalAttempts < 3 {
		return models.StatusToLearn
	}

	successRate := float64(successCount) / float64(totalAttempts)

	lastResultWeight := 0.3
	adjustedRate := successRate * (1 - lastResultWeight)
	if lastSuccess {
		adjustedRate += lastResultWeight
	}

	switch {
	case successRate > 0.8:
		return models.StatusLearned
	case successRate > 0.5:
		return models.StatusKnown
	default:
		return models.StatusToLearn
	}
}

func createPracticeWord(word models.UserWord) PracticeWord {
	questionTypes := []string{
		"original_to_translation",
		"translation_to_original",
		"fill_in_gap",
	}

	questionType := questionTypes[rand.Intn(len(questionTypes))]

	pw := PracticeWord{
		ID:           word.ID,
		OriginalWord: word.OriginalWord,
		Translation:  word.Translation,
		Example:      word.Example,
		QuestionType: questionType,
		Question:     "",
		Answer:       "",
	}

	switch questionType {
	case "original_to_translation":
		pw.Question = word.OriginalWord
		pw.Answer = word.Translation
	case "translation_to_original":
		pw.Question = word.Translation
		pw.Answer = word.OriginalWord
	case "fill_in_gap":
		if word.Example != "" {
			pw.Question = replaceWordWithGap(word.Example, word.OriginalWord)
			pw.Answer = word.OriginalWord
		} else {
			pw.Question = word.OriginalWord
			pw.Answer = word.Translation
			pw.QuestionType = "original_to_translation"
		}
	}

	if pw.Answer == "" {
		pw.Answer = word.Translation
	}

	return pw
}

func calculateNextReview(successCount, failCount int, success bool) time.Time {
	now := time.Now()

	switch {
	case successCount == 0 && failCount == 0:
		if success {
			return now.Add(6 * time.Hour)
		}
		return now.Add(1 * time.Hour)
	case success:
		interval := time.Duration(successCount) * 12 * time.Hour
		if interval > 30*24*time.Hour {
			interval = 30 * 24 * time.Hour
		}
		return now.Add(interval)
	default:
		return now.Add(1 * time.Hour)
	}
}

func replaceWordWithGap(example, word string) string {
	pattern := `(?i)\b` + regexp.QuoteMeta(word) + `\b`
	re := regexp.MustCompile(pattern)

	result := re.ReplaceAllString(example, "[_____]")

	if result == example {
		result = strings.Replace(example, word, "[_____]", -1)
	}

	if result == example {
		return example + "[_____]"
	}

	return result
}
