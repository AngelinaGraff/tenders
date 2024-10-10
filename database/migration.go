package database

import (
    "gorm.io/gorm"
    "tender-service/models"
)

func Migrate(db *gorm.DB) error {
	// Миграция таблиц
	err := db.AutoMigrate(&models.Tender{}, &models.Bid{}, &models.Employee{}, &models.TenderVersion{}, &models.BidVersion{}, &models.Review{})
	if err != nil {
		return err
	}

	// Добавление данных о пользователях (проверка, если данные не добавлены)
	var count int64
	db.Model(&models.Employee{}).Count(&count)
	if count == 0 {
		// Добавляем пользователей
		db.Create(&models.Employee{Username: "user1", FirstName: "John", LastName: "Doe"})
		db.Create(&models.Employee{Username: "user2", FirstName: "Jane", LastName: "Smith"})
		db.Create(&models.Employee{Username: "user3", FirstName: "Alice", LastName: "Johnson"})
	}

	return nil
}