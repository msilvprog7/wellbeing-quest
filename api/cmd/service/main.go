package main

import (
	"api.wellbeingquest.app/internal/handlers"

	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Activity struct {
	Name string `json:"name"`
	Feelings []string `json:"feelings"`
	Week string `json:"week"`
	Created time.Time `json:"created"`
}

type ActivitiesByWeek struct {
	Week string `json:"week"`
	Activities []Activity `json:"activities"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type Suggestions struct {
	Activities []string `json:"activities"`
	Feelings []string `json:"feelings"`
}

var activities = []Activity {
	createActivity("Read", []string{ "Relaxed" }),
}

func main() {
	router := gin.Default()
	router.POST("/activities/v1", handlers.PostActivity)
	router.GET("/activities/v1/weeks/:week", handlers.GetActivitiesByWeek)
	router.GET("/activities/v1/suggestions", getActivitySuggestions)
	router.Run("localhost:8080")
}

func postActivity(c *gin.Context) {
  // bind and validate activity
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

	// set time
	setTime(&activity)

	// append to in-memory list
	activities = append(activities, activity)

	// return json
	c.IndentedJSON(http.StatusCreated, activity)
}

func getActivitiesByWeek(c *gin.Context) {
	// bind and validate week
	week := c.Param("week")

	if week == "" {
		c.IndentedJSON(http.StatusBadRequest, ErrorMessage{ Message: "Get activities expected a non-empty 'week'" })
		return
	}

	activitiesByWeek := ActivitiesByWeek{ Week: week, Activities: filter(activities, week) }

	if len(activitiesByWeek.Activities) == 0 {
		c.IndentedJSON(http.StatusNotFound, ErrorMessage{ Message: fmt.Sprintf("No activities found for week '%s'", week) })
		return
	}

	c.IndentedJSON(http.StatusOK, activitiesByWeek)
}

func getActivitySuggestions(c *gin.Context) {
	suggestions := getSuggestions(activities)

	if len(suggestions.Activities) == 0 || len(suggestions.Feelings) == 0 {
		c.IndentedJSON(http.StatusNotFound, ErrorMessage{ Message: "No suggestions found" })
		return
	}

	c.IndentedJSON(http.StatusOK, suggestions)
}

func createActivity(name string, feelings []string) Activity {
	activity := Activity{ Name: name, Feelings: feelings }
	setTime(&activity)
	return activity
}

func setTime(activity *Activity) {
	activity.Created = time.Now()
	activity.Week = getWeek(&activity.Created)
}

func getWeek(time *time.Time) string {
	// Calculate the start of the week (Sunday)
	weekday := int(time.Weekday())
	startOfWeek := time.AddDate(0, 0, -weekday)
	return startOfWeek.Format("2006-01-02")
}

func filter(activities []Activity, week string) []Activity {
	var results []Activity

	for _, activity := range activities {
		if activity.Week == week {
			results = append(results, activity)
		}
	}

	return results
}

func getSuggestions(activities []Activity) Suggestions {
	result := Suggestions{ Activities: []string{}, Feelings: []string{} }
	activitiesSet := map[string]struct{}{}
	feelingsSet := map[string]struct{}{}

	for _, activity := range activities {
		if _, exists := activitiesSet[activity.Name]; !exists {
			result.Activities = append(result.Activities, activity.Name)
			activitiesSet[activity.Name] = struct{}{}
		}

		for _, feeling := range activity.Feelings {
			if _, exists := feelingsSet[feeling]; !exists {
				result.Feelings = append(result.Feelings, feeling)
				feelingsSet[feeling] = struct{}{}
			}
		}
	}

	return result
}