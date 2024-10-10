// models/tender.go
package models

import (
	"time"

	"github.com/google/uuid"
)

// Tender модель тендера
type Tender struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	ServiceType   string    `json:"service_type"`
	Status        string    `json:"status"`
	OrganizationID uint      `json:"organization_id"`
	CreatedBy     uuid.UUID `json:"created_by"` // Изменено на uuid.UUID
	Version       int       `json:"version"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
