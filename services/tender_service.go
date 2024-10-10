package services

import (
    "tender-service/models"
    "gorm.io/gorm"
)

func CreateTender(db *gorm.DB, tender *models.Tender) error {
    return db.Create(tender).Error
}

// Другие функции для работы с тендерами
