package handlers

import (
	"time"

	"api.wellbeingquest.app/internal/dtos"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostActivity(c *gin.Context) {
  // bind and validate activity
	var activity dtos.Activity
	err := c.BindJSON(&activity)
	if validationErr := ValidateActivity(&activity, err); validationErr != nil {
		c.IndentedJSON(http.StatusBadRequest, dtos.ErrorMessage {
			Message: fmt.Sprintf("Post activity request is invalid. error: %s", validationErr.Error()),
		})
		return
	}

	// add entry
	activity.Created = time.Now()
	entry := AddActivity(&activity)

	// set week
	activity.Week = entry.Week

	// return json
	c.IndentedJSON(http.StatusCreated, activity)
}

func GetActivitiesByWeek(c *gin.Context) {
	// bind and validate week
	week := dtos.Week{}
	week.Name = c.Param("week")

	if validationErr := ValidateWeek(&week, nil); validationErr != nil {
		c.IndentedJSON(http.StatusBadRequest, dtos.ErrorMessage {
			Message: fmt.Sprintf("Get activities by week request is invalid. error: %s", validationErr.Error()),
		})
		return
	}

	// get week and activities
	weekRetrieved, entriesByWeek, dataErr := GetWeekAndActivities(&week)

	if dataErr != nil {
		c.IndentedJSON(http.StatusNotFound, dtos.ErrorMessage {
			Message: fmt.Sprintf("Get activities by week request returned no results. error: %s", dataErr.Error()),
		})
		return
	}

	// populate the week
	week.Start = weekRetrieved.Start
	week.End = weekRetrieved.End
	week.Feelings = FormatActivitiesByFeelings(entriesByWeek)

	// return week
	c.IndentedJSON(http.StatusOK, week)
}

func GetActivitySuggestions(c *gin.Context) {
	// get activities and feelings
	activities, feelings, dataErr := GetActivitiesAndFeelings()
	if dataErr != nil {
		c.IndentedJSON(http.StatusNotFound, dtos.ErrorMessage{
			Message: fmt.Sprintf("Get activity suggestions returned no results. error: %s", dataErr.Error()),
		})
		return
	}

	// return suggestions
	suggestions := FormatSuggestions(activities, feelings)
	c.IndentedJSON(http.StatusOK, suggestions)
}