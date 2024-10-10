package controllers

import (
	"net/http"
	"tender-service/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetReviews возвращает список отзывов на предложения указанного автора
func GetReviews(c *gin.Context) {
	// Получение экземпляра базы данных из контекста
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database instance"})
		return
	}

	// Получение ID тендера из параметров URL
	tenderID := c.Param("tenderId")

	// Получение параметров запроса (имя автора и организация)
	authorUsername := c.Query("authorUsername")
	organizationID := c.Query("organizationId")

	// Проверка наличия параметров
	if authorUsername == "" || organizationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorUsername and organizationId are required"})
		return
	}

	// Поиск сотрудника по имени пользователя
	var employee models.Employee
	if err := db.Where("username = ?", authorUsername).First(&employee).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		return
	}

	// Поиск отзывов на предложения указанного автора для тендера
	var reviews []models.Review
	if err := db.Joins("JOIN bids ON reviews.bid_id = bids.id").
		Where("bids.tender_id = ? AND bids.created_by = ?", tenderID, employee.ID).
		Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	// Возвращаем список отзывов
	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

func CreateReview(c *gin.Context) {
	// Получение экземпляра базы данных из контекста
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database instance"})
		return
	}

	// Привязка входных данных
	var input struct {
		BidID          uint   `json:"bidId" binding:"required"`
		ReviewerID     uint   `json:"reviewerId" binding:"required"`
		AuthorUsername string `json:"authorUsername" binding:"required"`
		Content        string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка существования предложения
	var bid models.Bid
	if err := db.Where("id = ?", input.BidID).First(&bid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bid not found"})
		return
	}

	// Создание нового отзыва
	review := models.Review{
		BidID:          input.BidID,
		ReviewerID:     input.ReviewerID,
		AuthorUsername: input.AuthorUsername,
		Content:        input.Content,
	}

	// Сохранение отзыва в базе данных
	if err := db.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Review created", "review": review})
}