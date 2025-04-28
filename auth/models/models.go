package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint      `gorm:"primaryKey"`
	Email      string    `gorm:"unique;not null"`
	Phone      string    `gorm:"unique;nullable"`
	FirstName  string    `gorm:"not null"`
	LastName   string    `gorm:"not null"`
	Password   string    `gorm:"not null"`
	RoleID     uint      `gorm:"not null"`
	Role       Role      `gorm:"foreignKey:RoleID"`
	LastActive time.Time `gorm:"default:current_timestamp"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	ImageUrl   string         `gorm:"nullable"`
}

func (User) TableName() string {
	return "user"
}

type Role struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"unique;not null"`
	Users []User `gorm:"foreignKey:RoleID"`
}

func (Role) TableName() string {
	return "role"
}

type WordStatus struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

func (WordStatus) TableName() string {
	return "word_status"
}

type UserWord struct {
	ID           uint       `gorm:"primaryKey"`
	UserID       uint       `gorm:"not null"`
	User         User       `gorm:"foreignKey:UserID"`
	OriginalWord string     `gorm:"not null"`
	Translation  string     `gorm:"not null"`
	StatusID     uint       `gorm:"not null"`
	Status       WordStatus `gorm:"foreignKey:StatusID"`
	SuccessCount int        `gorm:"default:0"`
	LastReviewed *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (UserWord) TableName() string {
	return "user_word"
}

type LearningAttempt struct {
	ID          uint      `gorm:"primaryKey"`
	UserWordID  uint      `gorm:"not null"`
	UserWord    UserWord  `gorm:"foreignKey:UserWordID"`
	AttemptTime time.Time `gorm:"autoCreateTime"`
	Success     bool      `gorm:"not null"`
}

func (LearningAttempt) TableName() string {
	return "learning_attempt"
}
