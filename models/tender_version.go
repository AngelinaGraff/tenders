package models

import (
	"time"

	"github.com/google/uuid"
)

// TenderVersion представляет версию тендера
type TenderVersion struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	TenderID      uint      `json:"tender_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	ServiceType   string    `json:"service_type"`
	Status        string    `json:"status"`
	OrganizationID uint      `json:"organization_id"`
	CreatedBy     uuid.UUID `json:"created_by"`
	Version       int       `json:"version"`
	CreatedAt     time.Time `json:"created_at"`
}