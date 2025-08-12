package handlers

import (
	"api.wellbeingquest.app/internal/dtos"
	"api.wellbeingquest.app/internal/models"
)

func FormatActivitiesByFeelings(entries []models.Entry) []dtos.Feeling {
	// create a map for feelings to quickly lookup
	feelingsByName := map[string]dtos.Feeling{}

	// enumerate the entries
	for _, entry := range entries {
		// enumerate the feelings for each entry
		for _, feeling := range entry.Feelings {
			// create new feelings in result
			feelingStruct, exists := feelingsByName[feeling]
			if !exists {
				feelingStruct = dtos.Feeling{
					Name: feeling,
					Activities: []dtos.Activity{},
				}
			}

			// add activities to result
			activity := dtos.Activity{
				Name: entry.Activity,
				Feelings: entry.Feelings,
				Created: entry.Created,
				Week: entry.Week,
			}
			feelingStruct.Activities = append(feelingStruct.Activities, activity)
			feelingsByName[feeling] = feelingStruct
		}
	}

	return values(feelingsByName)
}

func FormatSuggestions(activities []models.Activity, feelings []models.Feeling) dtos.Suggestions {
	// format activities
	activitiesToSuggest := make([]dtos.Activity, len(activities))
	for i, activity := range activities {
		activitiesToSuggest[i] = dtos.Activity{
			Name: activity.Name,
			Feelings: activity.Feelings,
		}
	}

	// format feelings
	feelingsToSuggest := make([]dtos.Feeling, len(feelings))
	for i, feeling := range feelings {
		feelingsToSuggest[i] = dtos.Feeling{
			Name: feeling.Name,
			Activities: mapSlice(feeling.Activities, func (a string) dtos.Activity {
				return dtos.Activity{
					Name: a,
				}
			}),
		}
	}

	// returnsuggestions
	suggestions := dtos.Suggestions{
		Activities: activitiesToSuggest,
		Feelings: feelingsToSuggest,
	}
	return suggestions
}

// Helper functions related to formatting or arranging generic maps and lists
func createOrGetByName[T any](mapByName map[string]T, name string, createFunc func(string) T) T {
	if existing, exists := mapByName[name]; exists {
		return existing
	}

	newItem := createFunc(name)
	mapByName[name] = newItem
	return newItem
}

func mapSlice[T any, R any](slice []T, mapFunc func(T) R) []R {
	result := make([]R, len(slice))

	for i, item := range slice {
		result[i] = mapFunc(item)
	}

	return result
}

func values[K comparable, V any](m map[K]V) []V {
    values := make([]V, 0, len(m))
    for _, v := range m {
        values = append(values, v)
    }
    return values
}