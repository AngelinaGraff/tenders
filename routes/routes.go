package routes

import (
	"tender-service/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")

	api.GET("/ping", controllers.Ping) // Маршрут для проверки доступности сервера

	// Маршруты для тендеров
	api.POST("/tenders/new", controllers.CreateTender)	// Создание нового тендера

	api.GET("/tenders", controllers.GetTenders)	// Получение списка тендеров
	api.GET("/tenders/my", controllers.GetUserTenders)	// Получение тендеров пользователя
	
	api.GET("/tenders/:tenderId/status", controllers.GetTenderStatus) // Получение статуса тендера

	api.PATCH("/tenders/:id/edit", controllers.EditTender) // Редактирование тендера
	api.PUT("/tenders/:tenderId/rollback/:version", controllers.RollbackTenderVersion) // Откат версии тендера

	// Маршруты для предложений
	api.POST("/bids/new", controllers.CreateBid) // Создание нового предложения
	
	api.PUT("/bids/:bidId/submit_decision", controllers.SubmitDecision) // Подача решения по предложению
	
	api.GET("/bids/:tenderId/list", controllers.GetBidsForTender) // Получение списка предложений для тендера
	api.GET("/bids/my", controllers.GetUserBids) // Получение предложений пользователя
	
	api.GET("/bids/status/:bidId", controllers.GetBidStatus) // Получение статуса предложения

	api.PATCH("/bids/:bidId/edit", controllers.EditBid) // Редактирование предложения
	api.PUT("/bids/:bidId/rollback/:version", controllers.RollbackBidVersion) // Откат версии предложения

	api.GET("/bids/:tenderId/reviews", controllers.GetReviews) // Получение отзывов на предложения автора
	api.POST("/bids/feedback", controllers.CreateReview) // Создание отзыва на предложение
}
