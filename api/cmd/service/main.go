package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Activity struct {
	Name string `json:"name"`
	Feelings []string `json:"feelings"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

var activities = []Activity {
	{ Name: "Read", Feelings: []string{ "Relaxed" } },
}

func main() {
	router := gin.Default()
	router.POST("/activities/v1", postActivity)
	router.GET("/activities/v1/weeks/:week", getActivitiesByWeek)
	router.GET("/activities/v1/suggestions", getActivitySuggestions)
	router.Run("localhost:8080")
}

func postActivity(c *gin.Context) {
	var activity Activity

	if err := c.BindJSON(&activity); err != nil {
		c.IndentedJSON(http.StatusBadRequest, ErrorMessage{ Message: fmt.Sprintf("Post activity expected an activity, error: %s", err) })
		return
	}

	if activity.Name == "" {
		c.IndentedJSON(http.StatusBadRequest, ErrorMessage{ Message: "Post activity expected an activity with non-empty 'name'" })
		return
	}

	if len(activity.Feelings) == 0 {
		c.IndentedJSON(http.StatusBadRequest, ErrorMessage{ Message: "Post activity expected an activity with non-empty 'feelings'" })
		return
	}

	activities = append(activities, activity)
	c.IndentedJSON(http.StatusCreated, activity)
}

func getActivitiesByWeek(c *gin.Context) {
	c.IndentedJSON(http.StatusNotImplemented, ErrorMessage{ Message: "Get activities by week has not yet been implemented" })
}

func getActivitySuggestions(c *gin.Context) {
	c.IndentedJSON(http.StatusNotImplemented, ErrorMessage{ Message: "Get activity suggestions has not yet been implemented" })
}