package main

import (
	"api.wellbeingquest.app/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/activities/v1", handlers.PostActivity)
	router.GET("/activities/v1/weeks/:week", handlers.GetActivitiesByWeek)
	router.GET("/activities/v1/suggestions", handlers.GetActivitySuggestions)
	router.Run("localhost:8080")
}
