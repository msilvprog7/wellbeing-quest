package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

func main() {
	router := gin.Default()
	router.POST("/activities/v1", postActivity)
	router.GET("/activities/v1/weeks/:week", getActivitiesByWeek)
	router.GET("/activities/v1/suggestions", getActivitySuggestions)
	router.Run("localhost:8080")
}

func postActivity(c *gin.Context) {
	c.IndentedJSON(http.StatusNotImplemented, ErrorMessage{ Message: "Post activity has not yet been implemented" })
}

func getActivitiesByWeek(c *gin.Context) {
	c.IndentedJSON(http.StatusNotImplemented, ErrorMessage{ Message: "Get activities by week has not yet been implemented" })
}

func getActivitySuggestions(c *gin.Context) {
	c.IndentedJSON(http.StatusNotImplemented, ErrorMessage{ Message: "Get activity suggestions has not yet been implemented" })
}