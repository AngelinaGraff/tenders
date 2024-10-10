// controllers/tender_controller.go
package controllers

import (
	"net/http"
	"tender-service/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateTender(c *gin.Context) {
	// Получение экземпляра базы данных из контекста
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database instance"})
		return
	}

	var input struct {
		Name           string `json:"name" binding:"required"`
		Description    string `json:"description" binding:"required"`
		ServiceType    string `json:"serviceType" binding:"required"`
		Status         string `json:"status" binding:"required"`
		OrganizationID uint   `json:"organizationId" binding:"required"`
		CreatorUsername string `json:"creatorUsername" binding:"required"`
	}

	// Привязка данных из запроса к структуре
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Поиск пользователя по username
	var creator models.Employee
	if err := db.Where("username = ?", input.CreatorUsername).First(&creator).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Creator not found"})
		return
	}

	// Создание нового тендера
	tender := models.Tender{
		Name:           input.Name,
		Description:    input.Description,
		ServiceType:    input.ServiceType,
		Status:         input.Status,
		OrganizationID: input.OrganizationID,
		CreatedBy:      creator.ID, // Используем ID найденного пользователя
	}

	// Создание записи тендера в базе данных
	if err := db.Create(&tender).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Tender created", "tender": tender})
}


func EditTender(c *gin.Context) {
	// Получение экземпляра базы данных из контекста
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database instance"})
		return
	}

	// Получаем ID тендера из параметра URL
	id := c.Param("id")

	var tender models.Tender
	if err := db.First(&tender, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tender not found"})
		return
	}

	// Обновляем данные тендера из JSON-запроса
	if err := c.ShouldBindJSON(&tender); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Сохраняем обновления в базе данных
	if err := db.Save(&tender).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tender"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tender updated", "tender": tender})
}

// GetTenders возвращает список всех тендеров
func GetTenders(c *gin.Context) {
	// Получение экземпляра базы данных из контекста
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database instance"})
		return
	}

	var tenders []models.Tender

	// Извлечение всех тендеров из базы данных
	if err := db.Find(&tenders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем список тендеров
	c.JSON(http.StatusOK, gin.H{"tenders": tenders})
}

// GetUserTenders возвращает список тендеров пользователя
func GetUserTenders(c *gin.Context) {
	// Получаем имя пользователя из параметров запроса
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	// Получаем экземпляр базы данных
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database instance"})
		return
	}

	// Находим пользователя по имени
	var user models.Employee
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Извлекаем тендеры, созданные этим пользователем
	var tenders []models.Tender
	if err := db.Where("created_by = ?", user.ID).Find(&tenders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tenders"})
		return
	}

	// Возвращаем список тендеров
	c.JSON(http.StatusOK, gin.H{"tenders": tenders})
}

// RollbackTenderVersion откатывает тендер к указанной версии
func RollbackTenderVersion(c *gin.Context) {
	// Получение экземпляра базы данных из контекста
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database instance"})
		return
	}

	// Получение ID тендера и версии из параметров
	tenderID := c.Param("tenderId")
	version := c.Param("version")

	// Проверка на наличие тендера
	var tender models.Tender
	if err := db.Where("id = ?", tenderID).First(&tender).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tender not found"})
		return
	}

	// Поиск версии тендера
	var tenderVersion models.TenderVersion
	if err := db.Where("tender_id = ? AND version = ?", tenderID, version).First(&tenderVersion).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tender version not found"})
		return
	}

	// Обновление тендера до указанной версии
	tender.Name = tenderVersion.Name
	tender.Description = tenderVersion.Description
	tender.ServiceType = tenderVersion.ServiceType
	tender.Status = tenderVersion.Status
	tender.OrganizationID = tenderVersion.OrganizationID
	tender.CreatedBy = tenderVersion.CreatedBy
	tender.Version = tenderVersion.Version

	// Сохранение обновленного тендера
	if err := db.Save(&tender).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tender"})
		return
	}

	// Возвращение обновленного тендера
	c.JSON(http.StatusOK, gin.H{"message": "Tender rolled back", "tender": tender})
}

func GetTenderStatus(c *gin.Context) {
	// Получение экземпляра базы данных из контекста
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database instance"})
		return
	}

	// Получение ID тендера из параметров URL
	tenderID := c.Param("tenderId")

	// Поиск тендера по ID
	var tender models.Tender
	if err := db.Where("id = ?", tenderID).First(&tender).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tender not found"})
		return
	}

	// Возвращаем статус тендера
	c.JSON(http.StatusOK, gin.H{
		"tender_id": tender.ID,
		"status":    tender.Status,
	})
}