package handlers

import (
	"errors"
	"time"

	"api.wellbeingquest.app/internal/dtos"
	"api.wellbeingquest.app/internal/models"
)

var entries = []models.Entry{}
var activitiesByName = map[string]models.Activity{}
var feelingsByName = map[string]models.Feeling{}
var weeksByName = map[string]models.Week{}

func AddActivity(activityRequested *dtos.Activity) models.Entry {
	// Create or get activity, feelings, and week
	activity := createOrGetByName(activitiesByName, activityRequested.Name, func(name string) models.Activity {
		return models.Activity {
			Id: len(activitiesByName) + 1,
			Name: activityRequested.Name,
		}
	})

	feelings := []models.Feeling{}
	for _, feelingRequested := range activityRequested.Feelings {
		feeling := createOrGetByName(feelingsByName, feelingRequested, func(name string) models.Feeling {
			return models.Feeling {
				Id: len(feelingsByName) + 1,
				Name: feelingRequested,
			}
		})
		feelings = append(feelings, feeling)
	}

	weekRequested := getWeek(&activityRequested.Created)
	week := createOrGetByName(weeksByName, weekRequested, func(name string) models.Week {
		return models.Week {
			Id: len(weeksByName) + 1,
			Name: weekRequested,
		}
	})
	
	// Create entry
	entry := models.Entry {
		Id: len(entries) + 1,
		Activity: activity.Name,
		Feelings: mapSlice(feelings, func(f models.Feeling) string {
			return f.Name;
		}),
		Week: week.Name,
	}
	entries = append(entries, entry)
	return entry
}

func GetWeekAndActivities(weekRequested *dtos.Week) (models.Week, []models.Entry, error) {
	// get week
	if _, exists := weeksByName[weekRequested.Name]; !exists {
		return models.Week{}, []models.Entry{}, errors.New("week does not exist")
	}

	week := weeksByName[weekRequested.Name]

	// get entries
	entriesByWeek := []models.Entry{}
	for _, entry := range entries {
		if entry.Week == week.Name {
			entriesByWeek = append(entriesByWeek, entry)
		}
	}

	if len(entriesByWeek) == 0 {
		return week, entriesByWeek, errors.New("entries are empty")
	}

	return week, entriesByWeek, nil
}

func GetActivitiesAndFeelings() ([]models.Activity, []models.Feeling, error) {
	// get activities
	activities := values(activitiesByName)
	if len(activities) == 0 {
		return []models.Activity{}, []models.Feeling{}, errors.New("activities are empty")
	}

	// get feelings
	feelings := values(feelingsByName)
	if len(feelings) == 0 {
		return []models.Activity{}, []models.Feeling{}, errors.New("feelings are empty")
	}

	// return activities and feelings
	return activities, feelings, nil
}

// Helper functions related to representing data
func getWeek(time *time.Time) string {
	// Calculate the start of the week (Sunday)
	weekday := int(time.Weekday())
	startOfWeek := time.AddDate(0, 0, -weekday)
	return startOfWeek.Format("2006-01-02")
}
