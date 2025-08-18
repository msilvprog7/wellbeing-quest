package main

import (
	"api.wellbeingquest.app/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	dataHandler := handlers.NewInMemoryDataHandler()

	router := gin.Default()
	router.POST("/activities/v1", handlers.PostActivity(dataHandler))
	router.GET("/activities/v1/weeks/:week", handlers.GetActivitiesByWeek(dataHandler))
	router.GET("/activities/v1/suggestions", handlers.GetActivitySuggestions(dataHandler))

	router.Run("localhost:8080")
}
