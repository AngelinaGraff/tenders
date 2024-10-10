package models

import (
	"time"
)

type Review struct {
	ID             uint      `gorm:"primaryKey"`
	BidID          uint      `gorm:"not null"` // Связь с предложением
	ReviewerID     uint      `gorm:"not null"` // Кто оставил отзыв (ответственный)
	AuthorUsername string    `gorm:"not null"` // Автор предложения, на которого оставлен отзыв
	Content        string    `gorm:"type:text;not null"` // Текст отзыва
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}