package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"api.wellbeingquest.app/internal/dtos"

	"github.com/gin-gonic/gin"
)

func getUser(c *gin.Context) dtos.User {
	provider, _ := c.Get("user_provider")
	id, _ := c.Get("user_id")
	user := dtos.User{
		Provider: provider.(string),
		Id:       id.(string),
	}
	return user
}

func PostActivity(dataHandler DataHandler) gin.HandlerFunc {
  return func (c *gin.Context) {
		// bind and validate activity
		var activity dtos.Activity
		err := c.BindJSON(&activity)
		if validationErr := ValidateActivity(&activity, err); validationErr != nil {
			log.Printf("PostActivity validation error: %v", validationErr)
			c.IndentedJSON(http.StatusBadRequest, dtos.ErrorMessage {
				Message: fmt.Sprintf("Post activity request is invalid. error: %s", validationErr.Error()),
			})
			return
		}

		// add entry
		activity.Created = time.Now()
		user := getUser(c)
		entry, dataErr := dataHandler.AddActivity(&activity, &user)

		if dataErr != nil {
			log.Printf("PostActivity data error: %v", dataErr)
			c.IndentedJSON(http.StatusInternalServerError, dtos.ErrorMessage {
				Message: fmt.Sprintf("Post activity request failed. error: %s", dataErr.Error()),
			})
			return
		}

		// set week
		activity.Week = entry.Week

		// return json
		c.IndentedJSON(http.StatusCreated, activity)
	}
}

func GetActivitiesByWeek(dataHandler DataHandler) gin.HandlerFunc {
	return func (c *gin.Context) {
		// bind and validate week
		week := dtos.Week{}
		week.Name = c.Param("week")

		if validationErr := ValidateWeek(&week, nil); validationErr != nil {
			log.Printf("GetActivitiesByWeek validation error: %v", validationErr)
			c.IndentedJSON(http.StatusBadRequest, dtos.ErrorMessage {
				Message: fmt.Sprintf("Get activities by week request is invalid. error: %s", validationErr.Error()),
			})
			return
		}

		// get week and activities
		user := getUser(c)
		weekRetrieved, entriesByWeek, dataErr := dataHandler.GetWeekAndActivities(&week, &user)

		if dataErr != nil {
			log.Printf("GetActivitiesByWeek data error: %v", dataErr)
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
}

func GetActivitySuggestions(dataHandler DataHandler) gin.HandlerFunc {
	return func (c *gin.Context) {
		// get activities and feelings
		user := getUser(c)
		activities, feelings, dataErr := dataHandler.GetActivitiesAndFeelings(&user)
		if dataErr != nil {
			log.Printf("GetActivitySuggestions data error: %v", dataErr)
			c.IndentedJSON(http.StatusNotFound, dtos.ErrorMessage{
				Message: fmt.Sprintf("Get activity suggestions returned no results. error: %s", dataErr.Error()),
			})
			return
		}

		// return suggestions
		suggestions := FormatSuggestions(activities, feelings)
		c.IndentedJSON(http.StatusOK, suggestions)
	}
}