package handlers

import (
	"errors"
	"slices"

	"api.wellbeingquest.app/internal/dtos"
	"api.wellbeingquest.app/internal/models"
)

type InMemoryDataHandler struct {
	Entries []models.Entry
	ActivitiesByName map[string]models.Activity
	FeelingsByName map[string]models.Feeling
	WeeksByName map[string]models.Week
}

func NewInMemoryDataHandler() *InMemoryDataHandler {
	return &InMemoryDataHandler{
		Entries: []models.Entry{},
		ActivitiesByName: map[string]models.Activity{},
		FeelingsByName: map[string]models.Feeling{},
		WeeksByName: map[string]models.Week{},
	}
}

func (handler *InMemoryDataHandler) AddActivity(activityRequested *dtos.Activity) (models.Entry, error) {
	// Create or get activity, feelings, and week
	activity := createOrGetByName(handler.ActivitiesByName, activityRequested.Name, func(name string) models.Activity {
		return models.Activity {
			Id: len(handler.ActivitiesByName) + 1,
			Name: activityRequested.Name,
			Feelings: []string{},
		}
	})

	feelings := []models.Feeling{}
	for _, feelingRequested := range activityRequested.Feelings {
		feeling := createOrGetByName(handler.FeelingsByName, feelingRequested, func(name string) models.Feeling {
			return models.Feeling {
				Id: len(handler.FeelingsByName) + 1,
				Name: feelingRequested,
				Activities: []string{},
			}
		})
		feelings = append(feelings, feeling)
	}

	weekRequested := GetWeek(&activityRequested.Created)
	start, _ := GetTime(weekRequested)
	week := createOrGetByName(handler.WeeksByName, weekRequested, func(name string) models.Week {
		return models.Week {
			Id: len(handler.WeeksByName) + 1,
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

		handler.FeelingsByName[feeling.Name] = feeling
	}
	
	handler.ActivitiesByName[activity.Name] = activity
	
	// Create entry
	entry := models.Entry {
		Id: len(handler.Entries) + 1,
		Activity: activity.Name,
		Feelings: mapSlice(feelings, func(f models.Feeling) string {
			return f.Name;
		}),
		Week: week.Name,
		Created: activityRequested.Created,
	}
	handler.Entries = append(handler.Entries, entry)
	return entry, nil
}

func (handler *InMemoryDataHandler) GetWeekAndActivities(weekRequested *dtos.Week) (models.Week, []models.Entry, error) {
	// get week
	if _, exists := handler.WeeksByName[weekRequested.Name]; !exists {
		return models.Week{}, []models.Entry{}, errors.New("week does not exist")
	}

	week := handler.WeeksByName[weekRequested.Name]

	// get entries
	entriesByWeek := Filter(handler.Entries, func(e models.Entry) bool {
		return e.Week == week.Name
	})

	if len(entriesByWeek) == 0 {
		return week, entriesByWeek, errors.New("entries are empty")
	}

	return week, entriesByWeek, nil
}

func (handler *InMemoryDataHandler) GetActivitiesAndFeelings() ([]models.Activity, []models.Feeling, error) {
	// get activities
	activities := Values(handler.ActivitiesByName)
	if len(activities) == 0 {
		return []models.Activity{}, []models.Feeling{}, errors.New("activities are empty")
	}

	// get feelings
	feelings := Values(handler.FeelingsByName)
	if len(feelings) == 0 {
		return []models.Activity{}, []models.Feeling{}, errors.New("feelings are empty")
	}

	// return activities and feelings
	return activities, feelings, nil
}
