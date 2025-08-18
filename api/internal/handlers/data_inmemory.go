package handlers

import (
	"errors"
	"slices"

	"api.wellbeingquest.app/internal/dtos"
	"api.wellbeingquest.app/internal/models"
)

type InMemoryDataHandler struct {
	entries []models.Entry
	activitiesByName map[string]models.Activity
	feelingsByName map[string]models.Feeling
	weeksByName map[string]models.Week
}

func NewInMemoryDataHandler() *InMemoryDataHandler {
	return &InMemoryDataHandler{
		entries: []models.Entry{},
		activitiesByName: map[string]models.Activity{},
		feelingsByName: map[string]models.Feeling{},
		weeksByName: map[string]models.Week{},
	}
}

func (handler *InMemoryDataHandler) AddActivity(activityRequested *dtos.Activity) models.Entry {
	// Create or get activity, feelings, and week
	activity := createOrGetByName(handler.activitiesByName, activityRequested.Name, func(name string) models.Activity {
		return models.Activity {
			Id: len(handler.activitiesByName) + 1,
			Name: activityRequested.Name,
			Feelings: []string{},
		}
	})

	feelings := []models.Feeling{}
	for _, feelingRequested := range activityRequested.Feelings {
		feeling := createOrGetByName(handler.feelingsByName, feelingRequested, func(name string) models.Feeling {
			return models.Feeling {
				Id: len(handler.feelingsByName) + 1,
				Name: feelingRequested,
				Activities: []string{},
			}
		})
		feelings = append(feelings, feeling)
	}

	weekRequested := getWeek(&activityRequested.Created)
	start, _ := getTime(weekRequested)
	week := createOrGetByName(handler.weeksByName, weekRequested, func(name string) models.Week {
		return models.Week {
			Id: len(handler.weeksByName) + 1,
			Name: weekRequested,
			Start: start,
			End: start.AddDate(0, 0, 6),
		}
	})

	// Link activities and feelings
	for _, feeling := range feelings {
		if !slices.Contains(activity.Feelings, feeling.Name) {
			activity.Feelings = append(activity.Feelings, feeling.Name)
		}

		if !slices.Contains(feeling.Activities, activity.Name) {
		  feeling.Activities = append(feeling.Activities, activity.Name)
		}

		handler.feelingsByName[feeling.Name] = feeling
	}
	
	handler.activitiesByName[activity.Name] = activity
	
	// Create entry
	entry := models.Entry {
		Id: len(handler.entries) + 1,
		Activity: activity.Name,
		Feelings: mapSlice(feelings, func(f models.Feeling) string {
			return f.Name;
		}),
		Week: week.Name,
		Created: activityRequested.Created,
	}
	handler.entries = append(handler.entries, entry)
	return entry
}

func (handler *InMemoryDataHandler) GetWeekAndActivities(weekRequested *dtos.Week) (models.Week, []models.Entry, error) {
	// get week
	if _, exists := handler.weeksByName[weekRequested.Name]; !exists {
		return models.Week{}, []models.Entry{}, errors.New("week does not exist")
	}

	week := handler.weeksByName[weekRequested.Name]

	// get entries
	entriesByWeek := filter(handler.entries, func(e models.Entry) bool {
		return e.Week == week.Name
	})

	if len(entriesByWeek) == 0 {
		return week, entriesByWeek, errors.New("entries are empty")
	}

	return week, entriesByWeek, nil
}

func (handler *InMemoryDataHandler) GetActivitiesAndFeelings() ([]models.Activity, []models.Feeling, error) {
	// get activities
	activities := values(handler.activitiesByName)
	if len(activities) == 0 {
		return []models.Activity{}, []models.Feeling{}, errors.New("activities are empty")
	}

	// get feelings
	feelings := values(handler.feelingsByName)
	if len(feelings) == 0 {
		return []models.Activity{}, []models.Feeling{}, errors.New("feelings are empty")
	}

	// return activities and feelings
	return activities, feelings, nil
}
