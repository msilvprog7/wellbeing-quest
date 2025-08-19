package main

import (
	"fmt"
	"log"
	"os"

	"api.wellbeingquest.app/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// .env setup
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env file", err)
	}

	// Data handler setup
	dataHandler := handlers.NewInMemoryDataHandler()

	_, err := handlers.NewLocalHostDataHandler(getConnectionString(), "cmd/service/database.sql", getDatabaseDriver())
	if err != nil {
		log.Fatal("Failed to initialize localhost database handler:", err)
	}


	// Gin setup
	router := gin.Default()
	router.POST("/activities/v1", handlers.PostActivity(dataHandler))
	router.GET("/activities/v1/weeks/:week", handlers.GetActivitiesByWeek(dataHandler))
	router.GET("/activities/v1/suggestions", handlers.GetActivitySuggestions(dataHandler))

	router.Run("localhost:8080")
}

func getConnectionString() string {
	return fmt.Sprintf(
		"host=localhost port=5432 user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
}

func getDatabaseDriver() string {
	return os.Getenv("DB_DRIVER")
}