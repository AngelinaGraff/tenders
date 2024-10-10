package models

import (
	"time"
	"github.com/google/uuid"
)

type BidVersion struct {
	ID             uint      `gorm:"primaryKey"`
	BidID          uint      `gorm:"not null"`
	Version        int       `gorm:"not null"`
	Name           string    `gorm:"size:100;not null"`
	Description    string
	Status         string    `gorm:"size:20;not null"`
	CreatedBy      uuid.UUID `gorm:"type:uuid;not null"`
	OrganizationID uint      `gorm:"not null"`
	TenderID       uint      `gorm:"not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}
