package controllers

import (
	"log"
	"net/http"
	"tender-service/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

// CreateBid создает новое предложение (Bid)
func CreateBid(c *gin.Context) {
	// Получение экземпляра базы данных из контекста
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database instance"})
		return
	}

	// Привязка входных данных
	var input struct {
		Name           string `json:"name" binding:"required"`
		Description    string `json:"description" binding:"required"`
		Status         string `json:"status" binding:"required"`
		TenderID       uint   `json:"tenderId" binding:"required"`
		OrganizationID uint   `json:"organizationId" binding:"required"`
		CreatorUsername string `json:"creatorUsername" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error binding JSON:", err) // Логируем ошибку
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Поиск тендера по TenderID
	var tender models.Tender
	if err := db.Where("id = ?", input.TenderID).First(&tender).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tender not found"})
		return
	}

	// Поиск пользователя по CreatorUsername
	var creator models.Employee
	if err := db.Where("username = ?", input.CreatorUsername).First(&creator).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Creator not found"})
		return
	}

	// Создание нового предложения (Bid)
	bid := models.Bid{
		Name:           input.Name,
		Description:    input.Description,
		Status:         input.Status,
		TenderID:       input.TenderID,
		OrganizationID: input.OrganizationID,
		CreatedBy:      creator.ID,
	}

	// Сохранение предложения в базе данных
	if err := db.Create(&bid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bid"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Bid created", "bid": bid})
}

// GetUserBids возвращает список предложений пользователя
func GetUserBids(c *gin.Context) {
	// Получение экземпляра базы данных из контекста
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database instance"})
		return
	}

	// Получение параметра username из запроса
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	// Поиск пользователя по имени пользователя
	var user models.Employee
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Получение списка предложений пользователя
	var bids []models.Bid
	if err := db.Where("created_by = ?", user.ID).Find(&bids).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bids"})
		return
	}

	// Возвращаем список предложений
	c.JSON(http.StatusOK, gin.H{"bids": bids})
}

// GetBidsForTender возвращает список предложений для указанного тендера
func GetBidsForTender(c *gin.Context) {
	// Получение экземпляра базы данных из контекста
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database instance"})
		return
	}

	// Получение ID тендера из параметров URL
	tenderID := c.Param("tenderId")

	// Проверка, существует ли тендер
	var tender models.Tender
	if err := db.Where("id = ?", tenderID).First(&tender).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tender not found"})
		return
	}

	// Получение списка предложений для указанного тендера
	var bids []models.Bid
	if err := db.Where("tender_id = ?", tenderID).Find(&bids).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bids"})
		return
	}

	// Возвращаем список предложений
	c.JSON(http.StatusOK, gin.H{"bids": bids})
}

// EditBid редактирует предложение и создает новую версию
func EditBid(c *gin.Context) {
	// Получение экземпляра базы данных из контекста
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database instance"})
		return
	}

	// Получение ID предложения из параметров URL
	bidID := c.Param("bidId")

	// Поиск предложения по ID
	var bid models.Bid
	if err := db.Where("id = ?", bidID).First(&bid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bid not found"})
		return
	}

	// Привязка входных данных
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создание новой версии предложения
	bidVersion := models.BidVersion{
		BidID:          bid.ID,
		Version:        bid.Version,  // сохраняем текущую версию
		Name:           bid.Name,
		Description:    bid.Description,
		Status:         bid.Status,
		CreatedBy:      bid.CreatedBy,
		OrganizationID: bid.OrganizationID,
		TenderID:       bid.TenderID,
	}

	// Сохранение версии предложения в базе данных
	if err := db.Create(&bidVersion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bid version"})
		return
	}

	// Обновление текущего предложения
	bid.Name = input.Name
	bid.Description = input.Description
	bid.Version += 1 // Увеличиваем версию

	// Сохранение изменений в базе данных
	if err := db.Save(&bid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bid"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Bid updated", "bid": bid})
}

func RollbackBidVersion(c *gin.Context) {
	// Получение экземпляра базы данных из контекста
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database instance"})
		return
	}

	// Получение ID предложения и версии из параметров URL
	bidID := c.Param("bidId")
	versionParam := c.Param("version")

	// Преобразование версии из строки в число
	version, err := strconv.Atoi(versionParam)
	if err != nil || version < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version"})
		return
	}

	// Поиск предложения по ID
	var bid models.Bid
	if err := db.Where("id = ?", bidID).First(&bid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bid not found"})
		return
	}

	// Поиск версии предложения в таблице bid_versions
	var bidVersion models.BidVersion
	if err := db.Where("bid_id = ? AND version = ?", bidID, version).First(&bidVersion).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bid version not found"})
		return
	}

	// Обновление текущего предложения до параметров версии
	bid.Name = bidVersion.Name
	bid.Description = bidVersion.Description
	bid.Status = bidVersion.Status
	bid.Version = bidVersion.Version

	// Сохранение изменений в базе данных
	if err := db.Save(&bid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rollback bid"})
		return
	}

	// Возвращаем успешный ответ с откатной версией предложения
	c.JSON(http.StatusOK, gin.H{
		"message": "Bid rolled back to version",
		"bid":     bid,
	})
}

// SubmitDecision позволяет согласовать или отклонить предложение
func SubmitDecision(c *gin.Context) {
	// Получение экземпляра базы данных из контекста
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database instance"})
		return
	}

	// Получение ID предложения из URL
	bidID := c.Param("bidId")

	// Привязка данных из тела запроса
	var input struct {
		Decision string `json:"decision" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Поиск предложения по ID
	var bid models.Bid
	if err := db.Where("id = ?", bidID).First(&bid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bid not found"})
		return
	}

	// Поиск тендера по ID предложения
	var tender models.Tender
	if err := db.Where("id = ?", bid.TenderID).First(&tender).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tender not found"})
		return
	}

	// Проверка решения (approved или rejected)
	if input.Decision == "approved" {
		// Обновление статуса предложения
		bid.Status = "PUBLISHED"

		// Обновление статуса тендера (закрываем тендер)
		tender.Status = "CLOSED"

		// Сохранение изменений в базе данных
		if err := db.Save(&bid).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bid status"})
			return
		}
		if err := db.Save(&tender).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close tender"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Bid approved and tender closed",
			"bid":     bid,
			"tender":  tender,
		})

	} else if input.Decision == "rejected" {
		// Обновление статуса предложения
		bid.Status = "CANCELED"

		// Сохранение изменений в базе данных
		if err := db.Save(&bid).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bid status"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Bid rejected",
			"bid":     bid,
		})

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid decision, must be 'approved' or 'rejected'"})
	}
}

func GetBidStatus(c *gin.Context) {
	// Получение экземпляра базы данных из контекста
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database instance"})
		return
	}

	// Получение ID предложения из параметров URL
	bidID := c.Param("bidId")

	// Поиск предложения по ID
	var bid models.Bid
	if err := db.Where("id = ?", bidID).First(&bid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bid not found"})
		return
	}

	// Возвращаем статус предложения
	c.JSON(http.StatusOK, gin.H{
		"bid_id":  bid.ID,
		"status":  bid.Status,
	})
}