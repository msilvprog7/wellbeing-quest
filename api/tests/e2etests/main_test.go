package e2etests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"

	"api.wellbeingquest.app/internal/dtos"
	"api.wellbeingquest.app/internal/handlers"
	"api.wellbeingquest.app/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAddActivityWhenNoEntries(t *testing.T) {
	// Arrange
	router := setupRouter()
	activity := dtos.Activity{
		Name: "Read",
		Feelings: []string{"Relaxed", "Accomplished"},
	}

	// Act
	w := postActivity(router, &activity)

	// Assert
	assert.Equal(t, 201, w.Code)
	assertActivity(t, w.Body, &activity, time.Minute)
}

func TestAddActivityWhenMultipleEntries(t *testing.T) {
	// Arrange
	router := setupRouter()
	activities := []dtos.Activity{
		{
			Name: "Read",
			Feelings: []string{"Relaxed", "Accomplished"},
		},
		{
			Name: "Write",
			Feelings: []string{"Relaxed", "Creative", "Excited"},
		},
		{
			Name: "Coffee",
			Feelings: []string{"Relaxed"},
		},
	}

	// Act and Assert
	for i := range activities {
		w := postActivity(router, &activities[i])

		assert.Equal(t, 201, w.Code)
		assertActivity(t, w.Body, &activities[i], time.Minute)
	}
}

func TestGetActivitiesByWeekWhenNoEntries(t *testing.T) {
	// Arrange
	router := setupRouter()
	time := time.Now()
	week := handlers.GetWeek(&time)

	// Act
	w := getActivitiesByWeek(router, week)

	// Assert
	assert.Equal(t, 404, w.Code)
	assertErrorMessage(t, w.Body, "Get activities by week request returned no results.")
}

func TestGetActivitiesByWeekWhenMultipleEntries(t *testing.T) {
	// Arrange
	router := setupRouter()
	time := time.Now()
	week := handlers.GetWeek(&time)

	activities := []dtos.Activity{
		{
			Name: "Read",
			Feelings: []string{"Relaxed", "Accomplished"},
		},
		{
			Name: "Write",
			Feelings: []string{"Relaxed", "Creative", "Excited"},
		},
		{
			Name: "Coffee",
			Feelings: []string{"Relaxed"},
		},
	}

	start, _ := handlers.GetTime(week)
	expected := dtos.Week{
		Name: week,
		Start: start,
		End: start.AddDate(0, 0, 6),
		Feelings: []dtos.Feeling{
			{
				Name: "Relaxed",
				Activities: []dtos.Activity{
					{
						Name: "Read",
					},
					{
						Name: "Write",
					},
					{
						Name: "Coffee",
					},
				},
			},
			{
				Name: "Accomplished",
				Activities: []dtos.Activity{
					{
						Name: "Read",
					},
				},
			},
			{
				Name: "Creative",
				Activities: []dtos.Activity{
					{
						Name: "Write",
					},
				},
			},
			{
				Name: "Excited",
				Activities: []dtos.Activity{
					{
						Name: "Write",
					},
				},
			},
		},
	}

	// Act
	for i := range activities {
		postActivity(router, &activities[i])
	}

	w := getActivitiesByWeek(router, week)

	// Assert
	assert.Equal(t, 200, w.Code)
	assertWeek(t, w.Body, &expected)
}

func TestGetActivitySuggestionsWhenNoEntries(t *testing.T) {
	// Arrange
	router := setupRouter()

	// Act
	w := getActivitySuggestions(router)

	// Assert
	assert.Equal(t, 404, w.Code)
	assertErrorMessage(t, w.Body, "Get activity suggestions returned no results.")
}

func TestGetActivitySuggestionsWhenMultipleEntries(t *testing.T) {
	// Arrange
	router := setupRouter()
	activities := []dtos.Activity{
		{
			Name: "Read",
			Feelings: []string{"Relaxed", "Accomplished"},
		},
		{
			Name: "Write",
			Feelings: []string{"Relaxed", "Creative", "Excited"},
		},
		{
			Name: "Coffee",
			Feelings: []string{"Relaxed"},
		},
	}

	expected := dtos.Suggestions{
		Activities: []dtos.Activity{
			{
				Name: "Read",
			},
			{
				Name: "Write",
			},
			{
				Name: "Coffee",
			},
		},
		Feelings: []dtos.Feeling{
			{
				Name: "Relaxed",
			},
			{
				Name: "Accomplished",
			},
			{
				Name: "Creative",
			},
			{
				Name: "Excited",
			},
		},
	}

	// Act
	for i := range activities {
		postActivity(router, &activities[i])
	}

	w := getActivitySuggestions(router)

	// Assert
	assert.Equal(t, 200, w.Code)
	assertSuggestions(t, w.Body, &expected)
}

/**
 * Arrangements
 */
func setupRouter() *gin.Engine {
	services.LoadDotEnv(".env")
	dataHandler := services.GetDataHandler()
	return services.SetupRouter(dataHandler)
}

/**
 * Acts
 */
func postActivity(router *gin.Engine, activity *dtos.Activity) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	body, _ := json.Marshal(activity)
	req, _ := http.NewRequest("POST", "/activities/v1", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)
	return w
}

func getActivitiesByWeek(router *gin.Engine, week string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/activities/v1/weeks/%s", week), nil)
	router.ServeHTTP(w, req)
	return w
}

func getActivitySuggestions(router *gin.Engine) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/activities/v1/suggestions", nil)
	router.ServeHTTP(w, req)
	return w
}

/**
 * Assertions
 */
func assertActivity(t *testing.T, body *bytes.Buffer, expected *dtos.Activity, within time.Duration) {
	var actual dtos.Activity
	err := json.Unmarshal(body.Bytes(), &actual)
	assert.NoError(t, err)

	assert.Equal(t, actual.Name, expected.Name)
	assert.Equal(t, len(actual.Feelings), len(expected.Feelings))
	for i := range actual.Feelings {
		assert.Equal(t, actual.Feelings[i], expected.Feelings[i])
	}

	assert.WithinDuration(t, time.Now(), actual.Created, within)
	assert.Equal(t, actual.Week, handlers.GetWeek(&actual.Created))
}

func assertErrorMessage(t *testing.T, body *bytes.Buffer, expected string) {
	var actual dtos.ErrorMessage
	err := json.Unmarshal(body.Bytes(), &actual)
	assert.NoError(t, err)

	assert.Contains(t, actual.Message, expected)
}

func assertWeek(t *testing.T, body *bytes.Buffer, expected *dtos.Week) {
	var actual dtos.Week
	err := json.Unmarshal(body.Bytes(), &actual)
	assert.NoError(t, err)

	assert.Equal(t, actual.Name, expected.Name)

	expectedStart, _ := handlers.GetTime(expected.Name)
	assert.Equal(t, actual.Start, expectedStart)
	assert.Equal(t, actual.End, expectedStart.AddDate(0, 0, 6))

	// Sort feelings by name for consistent comparison
	sort.Slice(actual.Feelings, func(i, j int) bool {
			return actual.Feelings[i].Name < actual.Feelings[j].Name
	})
	sort.Slice(expected.Feelings, func(i, j int) bool {
			return expected.Feelings[i].Name < expected.Feelings[j].Name
	})
	
	assert.Equal(t, len(actual.Feelings), len(expected.Feelings))
	for i := range actual.Feelings {
		assert.Equal(t, actual.Feelings[i].Name, expected.Feelings[i].Name)
		
		// Sort activities by name for consistent comparison
		sort.Slice(actual.Feelings[i].Activities, func(x, y int) bool {
				return actual.Feelings[i].Activities[x].Name < actual.Feelings[i].Activities[y].Name
		})
		sort.Slice(expected.Feelings[i].Activities, func(x, y int) bool {
				return expected.Feelings[i].Activities[x].Name < expected.Feelings[i].Activities[y].Name
		})
		
		assert.Equal(t, len(actual.Feelings[i].Activities), len(expected.Feelings[i].Activities))
		for j := range actual.Feelings[i].Activities {
			assert.Equal(t, actual.Feelings[i].Activities[j].Name, expected.Feelings[i].Activities[j].Name)
		}
	}
}

func assertSuggestions(t *testing.T, body *bytes.Buffer, expected *dtos.Suggestions) {
	var actual dtos.Suggestions
	err := json.Unmarshal(body.Bytes(), &actual)
	assert.NoError(t, err)

	// Sort activities by name for consistent comparison
	sort.Slice(actual.Activities, func(i, j int) bool {
		return actual.Feelings[i].Name < actual.Feelings[j].Name
	})

	sort.Slice(expected.Activities, func(i, j int) bool {
		return expected.Feelings[i].Name < expected.Feelings[j].Name
	})

	assert.Equal(t, len(actual.Activities), len(expected.Activities))
	for i := range actual.Activities {
		assert.Equal(t, actual.Activities[i].Name, expected.Activities[i].Name)
	}

	// Sort feelings by name for consistent comparison
	sort.Slice(actual.Feelings, func(i, j int) bool {
		return actual.Feelings[i].Name < actual.Feelings[j].Name
	})

	sort.Slice(expected.Feelings, func(i, j int) bool {
		return expected.Feelings[i].Name < expected.Feelings[j].Name
	})

	assert.Equal(t, len(actual.Feelings), len(expected.Feelings))
	for i := range actual.Feelings {
		assert.Equal(t, actual.Feelings[i].Name, expected.Feelings[i].Name)
	}
}