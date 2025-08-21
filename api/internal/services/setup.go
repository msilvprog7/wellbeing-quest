package services

import (
	"fmt"
	"log"
	"os"

	"api.wellbeingquest.app/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func LoadDotEnv(filename string) {
	if err := godotenv.Load(filename); err != nil {
		log.Fatal("Failed to load .env file", err)
	}

	log.Printf("Loaded .env filename, %s\n", filename)
}

func GetDataHandler() handlers.DataHandler {
	switch os.Getenv("DB_MODE") {
		case "localhost":
			// connection string
			connectionString := fmt.Sprintf(
				"host=localhost port=5432 user=%s password=%s dbname=%s sslmode=disable",
				os.Getenv("DB_USER"),
				os.Getenv("DB_PASSWORD"),
				os.Getenv("DB_NAME"))

			// setup scripts
			setup := []string{}
			directory := os.Getenv("DB_SQLDIRECTORY")
			if os.Getenv("DB_RESET") == "reset" {
				setup = append(setup, fmt.Sprintf("%s/database_reset.sql", directory))
			}
			setup = append(setup, fmt.Sprintf("%s/database.sql", directory))
			
			// create handler
			localhostDataHandler, err := handlers.NewLocalHostDataHandler(os.Getenv("DB_DRIVER"), connectionString, setup)

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

func SetupRouter(dataHandler handlers.DataHandler) *gin.Engine {
	router := gin.Default()
	router.POST("/activities/v1", handlers.PostActivity(dataHandler))
	router.GET("/activities/v1/weeks/:week", handlers.GetActivitiesByWeek(dataHandler))
	router.GET("/activities/v1/suggestions", handlers.GetActivitySuggestions(dataHandler))
	return router
}
