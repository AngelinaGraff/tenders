package services

import (
    "tender-service/models"
    "gorm.io/gorm"
)

func CreateBid(db *gorm.DB, bid *models.Bid) error {
    return db.Create(bid).Error
}

// Другие функции для работы с предложениями
