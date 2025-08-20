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
	// Resource setup
	loadDotEnv()
	dataHandler := getDataHandler()


	// Gin setup
	router := gin.Default()
	router.POST("/activities/v1", handlers.PostActivity(dataHandler))
	router.GET("/activities/v1/weeks/:week", handlers.GetActivitiesByWeek(dataHandler))
	router.GET("/activities/v1/suggestions", handlers.GetActivitySuggestions(dataHandler))

	router.Run("localhost:8080")
}

func loadDotEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env file", err)
	}

	log.Println("Loaded .env")
}

func getDataHandler() handlers.DataHandler {
	switch getDatabaseMode() {
	case "localhost":
		localhostDataHandler, err := handlers.NewLocalHostDataHandler(getDatabaseDriver(), getConnectionString(), getSetup())
		if err != nil {
			log.Fatal("Failed to initialize localhost database handler:", err)
		}

		log.Println("Setup localhost data handler")
		return localhostDataHandler
	default:
		log.Println("Setup inmemory data handler")
		return handlers.NewInMemoryDataHandler()
	}
}

func getDatabaseMode() string {
	return os.Getenv("DB_MODE")
}

func getDatabaseReset() string {
	return os.Getenv("DB_RESET")
}

func getDatabaseDriver() string {
	return os.Getenv("DB_DRIVER")
}

func getConnectionString() string {
	return fmt.Sprintf(
		"host=localhost port=5432 user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
}

func getSetup() []string {
	setup := []string{}

	if getDatabaseReset() == "reset" {
		setup = append(setup, "cmd/service/database_reset.sql")
	}

	setup = append(setup, "cmd/service/database.sql")

	return setup
}