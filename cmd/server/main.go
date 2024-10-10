package main

import (
	"log"
	"tender-service/config"
	"tender-service/database"
	"tender-service/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключение к базе данных
	db, err := database.ConnectDatabase(cfg.PostgresConn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Миграция базы данных
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Инициализация роутера Gin
	router := gin.Default()

	// Передача экземпляра базы данных через middleware в контекст
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Установка маршрутов
	routes.SetupRoutes(router)

	// Запуск сервера
	router.Run(cfg.ServerAddress)
}
