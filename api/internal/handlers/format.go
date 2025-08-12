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

func values[K comparable, V any](m map[K]V) []V {
    values := make([]V, 0, len(m))
    for _, v := range m {
        values = append(values, v)
    }
    return values
}