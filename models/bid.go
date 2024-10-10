// models/bid.go
package models

import (
	"time"

	"github.com/google/uuid"
)

// Bid представляет предложение (Bid)
type Bid struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	TenderID      uint      `json:"tender_id"`
	OrganizationID uint      `json:"organization_id"`
	CreatedBy     uuid.UUID `json:"created_by"`
	Version       int       `json:"version" gorm:"default:1"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
