package handlers

import (
	"time"

	"api.wellbeingquest.app/internal/dtos"
	"api.wellbeingquest.app/internal/models"
)

var entries = []models.Entry{}
var activitiesByName = map[string]models.Activity{}
var feelingsByName = map[string]models.Feeling{}
var weeksByName = map[string]models.Week{}

func AddActivity(activityRequested dtos.Activity) models.Entry {
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

func createOrGetByName[T any](mapByName map[string]T, name string, createFunc func(string) T) T {
	if existing, exists := mapByName[name]; exists {
		return existing
	}

	newItem := createFunc(name)
	mapByName[name] = newItem
	return newItem
}

func getWeek(time *time.Time) string {
	// Calculate the start of the week (Sunday)
	weekday := int(time.Weekday())
	startOfWeek := time.AddDate(0, 0, -weekday)
	return startOfWeek.Format("2006-01-02")
}

func mapSlice[T any, R any](slice []T, mapFunc func(T) R) []R {
	result := make([]R, len(slice))

	for i, item := range slice {
		result[i] = mapFunc(item)
	}

	return result
}