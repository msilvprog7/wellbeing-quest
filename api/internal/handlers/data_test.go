package handlers

import (
	"testing"
	"time"

	"api.wellbeingquest.app/internal/dtos"
	"api.wellbeingquest.app/internal/models"
)

func TestAddActivityWhenNoEntries(t *testing.T) {
	// Arrange
	dataHandler := NewInMemoryDataHandler()
	activity := dtos.Activity {
		Name: "Read",
		Feelings: []string{"Relaxed", "Accomplished"},
		Created: time.Date(2025, 8, 18, 0, 0, 0, 0, time.UTC),
	}

	// Act
	entry := dataHandler.AddActivity(&activity)

	// Assert
	assertEntry(t, &entry, &activity, 1)
	assertActivity(t, dataHandler, activity.Name, 1, activity.Feelings)
	assertFeeling(t, dataHandler, activity.Feelings[0], 1, []string{activity.Name})
	assertFeeling(t, dataHandler, activity.Feelings[1], 2, []string{activity.Name})
	assertWeek(t, dataHandler, getWeek(&activity.Created), 1)
}

func TestAddActivityWhenExistingActivity(t *testing.T) {
	// Arrange
	dataHandler := NewInMemoryDataHandler()
	dataHandler.entries = []models.Entry{
		models.Entry{
			Id: 1,
			Activity: "Read",
			Feelings: []string{"Relaxed", "Accomplished"},
			Week: "2025-08-17",
			Created: time.Date(2025, 8, 18, 0, 0, 0, 0, time.UTC),
		},
	}
	dataHandler.activitiesByName = map[string]models.Activity{
		"Read": models.Activity{
			Id: 1,
			Name: "Read",
			Feelings: []string{"Relaxed", "Accomplished"},
		},
	}
  dataHandler.feelingsByName = map[string]models.Feeling{
		"Relaxed": models.Feeling{
			Id: 1,
			Name: "Relaxed",
			Activities: []string{"Read"},
		},
		"Accomplished": models.Feeling{
			Id: 2,
			Name: "Accomplished",
			Activities: []string{"Read"},
		},
	}
  dataHandler.weeksByName = map[string]models.Week{
		"2025-08-17": models.Week{
			Id: 1,
			Name: "2025-08-17",
			Start: time.Date(2025, 8, 17, 0, 0, 0, 0, time.UTC),
			End: time.Date(2025, 8, 23, 0, 0, 0, 0, time.UTC),
		},
	}

	activity := dtos.Activity {
		Name: "Read",
		Feelings: []string{"Focused"},
		Created: time.Date(2025, 8, 24, 0, 0, 0, 0, time.UTC),
	}

	// Act
	entry := dataHandler.AddActivity(&activity)

	// Assert
	assertEntry(t, &entry, &activity, 2)
	assertActivity(t, dataHandler, activity.Name, 1, []string{"Relaxed", "Accomplished", "Focused"})
	assertFeeling(t, dataHandler, "Focused", 3, []string{activity.Name})
	assertWeek(t, dataHandler, getWeek(&activity.Created), 2)
}

func TestAddActivityWhenExistingFeeling(t *testing.T) {
	// Arrange
	dataHandler := NewInMemoryDataHandler()
	dataHandler.entries = []models.Entry{
		models.Entry{
			Id: 1,
			Activity: "Read",
			Feelings: []string{"Relaxed", "Accomplished"},
			Week: "2025-08-17",
			Created: time.Date(2025, 8, 18, 0, 0, 0, 0, time.UTC),
		},
	}
	dataHandler.activitiesByName = map[string]models.Activity{
		"Read": models.Activity{
			Id: 1,
			Name: "Read",
			Feelings: []string{"Relaxed", "Accomplished"},
		},
	}
  dataHandler.feelingsByName = map[string]models.Feeling{
		"Relaxed": models.Feeling{
			Id: 1,
			Name: "Relaxed",
			Activities: []string{"Read"},
		},
		"Accomplished": models.Feeling{
			Id: 2,
			Name: "Accomplished",
			Activities: []string{"Read"},
		},
	}
  dataHandler.weeksByName = map[string]models.Week{
		"2025-08-17": models.Week{
			Id: 1,
			Name: "2025-08-17",
			Start: time.Date(2025, 8, 17, 0, 0, 0, 0, time.UTC),
			End: time.Date(2025, 8, 23, 0, 0, 0, 0, time.UTC),
		},
	}

	activity := dtos.Activity {
		Name: "Write",
		Feelings: []string{"Relaxed"},
		Created: time.Date(2025, 8, 24, 0, 0, 0, 0, time.UTC),
	}

	// Act
	entry := dataHandler.AddActivity(&activity)

	// Assert
	assertEntry(t, &entry, &activity, 2)
	assertActivity(t, dataHandler, activity.Name, 2, []string{"Relaxed"})
	assertFeeling(t, dataHandler, "Relaxed", 1, []string{"Read", "Write"})
	assertWeek(t, dataHandler, getWeek(&activity.Created), 2)
}

func TestAddActivityWhenExistingWeek(t *testing.T) {
	// Arrange
	dataHandler := NewInMemoryDataHandler()
	dataHandler.entries = []models.Entry{
		models.Entry{
			Id: 1,
			Activity: "Read",
			Feelings: []string{"Relaxed", "Accomplished"},
			Week: "2025-08-17",
			Created: time.Date(2025, 8, 18, 0, 0, 0, 0, time.UTC),
		},
	}
	dataHandler.activitiesByName = map[string]models.Activity{
		"Read": models.Activity{
			Id: 1,
			Name: "Read",
			Feelings: []string{"Relaxed", "Accomplished"},
		},
	}
  dataHandler.feelingsByName = map[string]models.Feeling{
		"Relaxed": models.Feeling{
			Id: 1,
			Name: "Relaxed",
			Activities: []string{"Read"},
		},
		"Accomplished": models.Feeling{
			Id: 2,
			Name: "Accomplished",
			Activities: []string{"Read"},
		},
	}
  dataHandler.weeksByName = map[string]models.Week{
		"2025-08-17": models.Week{
			Id: 1,
			Name: "2025-08-17",
			Start: time.Date(2025, 8, 17, 0, 0, 0, 0, time.UTC),
			End: time.Date(2025, 8, 23, 0, 0, 0, 0, time.UTC),
		},
	}

	activity := dtos.Activity {
		Name: "Write",
		Feelings: []string{"Focused"},
		Created: time.Date(2025, 8, 23, 0, 0, 0, 0, time.UTC),
	}

	// Act
	entry := dataHandler.AddActivity(&activity)

	// Assert
	assertEntry(t, &entry, &activity, 2)
	assertActivity(t, dataHandler, activity.Name, 2, []string{"Focused"})
	assertFeeling(t, dataHandler, "Focused", 3, []string{"Write"})
	assertWeek(t, dataHandler, getWeek(&activity.Created), 1)
}

func TestGetWeekAndActivitiesWhenNoWeek(t *testing.T) {
	// Arrange
	dataHandler := NewInMemoryDataHandler()
	dataHandler.weeksByName = map[string]models.Week{}

	weekRequested := dtos.Week{
		Name: "2025-08-17",
	}

	// Act
	_, _, err := dataHandler.GetWeekAndActivities(&weekRequested)

	// Assert
	if err.Error() != "week does not exist" {
		t.Errorf("Error was '%s' but should have been '%s'", err.Error(), "week does not exist")
	}
}

func TestGetWeekAndActivitiesWhenNoEntries(t *testing.T) {
	// Arrange
	dataHandler := NewInMemoryDataHandler()
	dataHandler.entries = []models.Entry{}
  dataHandler.weeksByName = map[string]models.Week{
		"2025-08-17": models.Week{
			Id: 1,
			Name: "2025-08-17",
			Start: time.Date(2025, 8, 17, 0, 0, 0, 0, time.UTC),
			End: time.Date(2025, 8, 23, 0, 0, 0, 0, time.UTC),
		},
	}

	weekRequested := dtos.Week{
		Name: "2025-08-17",
	}

	// Act
	_, _, err := dataHandler.GetWeekAndActivities(&weekRequested)

	// Assert
	if err.Error() != "entries are empty" {
		t.Errorf("Error was '%s' but should have been '%s'", err.Error(), "entries are empty")
	}
}

func TestGetWeekAndActivitiesWhenEntries(t *testing.T) {
	// Arrange
	dataHandler := NewInMemoryDataHandler()
	dataHandler.entries = []models.Entry{
		models.Entry{
			Id: 1,
			Activity: "Read",
			Feelings: []string{"Relaxed", "Accomplished"},
			Week: "2025-08-17",
			Created: time.Date(2025, 8, 18, 0, 0, 0, 0, time.UTC),
		},
		models.Entry{
			Id: 2,
			Activity: "Write",
			Feelings: []string{"Relaxed", "Creative"},
			Week: "2025-08-17",
			Created: time.Date(2025, 8, 19, 0, 0, 0, 0, time.UTC),
		},
		models.Entry{
			Id: 3,
			Activity: "Read",
			Feelings: []string{"Relaxed", "Accomplished"},
			Week: "2025-08-24",
			Created: time.Date(2025, 8, 24, 0, 0, 0, 0, time.UTC),
		},
	}
  dataHandler.weeksByName = map[string]models.Week{
		"2025-08-17": models.Week{
			Id: 1,
			Name: "2025-08-17",
			Start: time.Date(2025, 8, 17, 0, 0, 0, 0, time.UTC),
			End: time.Date(2025, 8, 23, 0, 0, 0, 0, time.UTC),
		},
	}

	weekRequested := dtos.Week{
		Name: "2025-08-17",
	}

	expectedEntries := filter(dataHandler.entries, func(e models.Entry) bool {
		return e.Week == weekRequested.Name
	})

	// Act
	week, actualEntries, err := dataHandler.GetWeekAndActivities(&weekRequested)

	// Assert
	if err != nil {
		t.Errorf("Error was '%s' but should have been nil", err.Error())
	}

	assertWeekFields(t, &week, weekRequested.Name, 1)
	assertEntries(t, actualEntries, expectedEntries)
}

func TestGetActivitiesAndFeelingsWhenNoActivities(t *testing.T) {
	// Arrange
	dataHandler := NewInMemoryDataHandler()
	dataHandler.activitiesByName = map[string]models.Activity{}
  dataHandler.feelingsByName = map[string]models.Feeling{
		"Relaxed": models.Feeling{
			Id: 1,
			Name: "Relaxed",
			Activities: []string{"Read"},
		},
		"Accomplished": models.Feeling{
			Id: 2,
			Name: "Accomplished",
			Activities: []string{"Read"},
		},
		"Creative": models.Feeling{
			Id: 3,
			Name: "Creative",
			Activities: []string{"Write"},
		},
	}

	// Act
	_, _, err := dataHandler.GetActivitiesAndFeelings()

	// Assert
	if err.Error() != "activities are empty" {
		t.Errorf("Error message is '%s' but expected '%s'", err.Error(), "activities are empty")
	}
}

func TestGetActivitiesAndFeelingsWhenNoFeelings(t *testing.T) {
	// Arrange
	dataHandler := NewInMemoryDataHandler()
	dataHandler.activitiesByName = map[string]models.Activity{
		"Read": models.Activity{
			Id: 1,
			Name: "Read",
			Feelings: []string{"Relaxed", "Accomplished"},
		},
		"Write": models.Activity{
			Id: 2,
			Name: "Write",
			Feelings: []string{"Creative"},
		},
	}
  dataHandler.feelingsByName = map[string]models.Feeling{}

	// Act
	_, _, err := dataHandler.GetActivitiesAndFeelings()

	// Assert
	if err.Error() != "feelings are empty" {
		t.Errorf("Error message is '%s' but expected '%s'", err.Error(), "feelings are empty")
	}
}

func TestGetActivitiesAndFeelingsWheActivitiesAndFeelings(t *testing.T) {
	// Arrange
	dataHandler := NewInMemoryDataHandler()
	dataHandler.activitiesByName = map[string]models.Activity{
		"Read": models.Activity{
			Id: 1,
			Name: "Read",
			Feelings: []string{"Relaxed", "Accomplished"},
		},
		"Write": models.Activity{
			Id: 2,
			Name: "Write",
			Feelings: []string{"Creative"},
		},
	}
  dataHandler.feelingsByName = map[string]models.Feeling{
		"Relaxed": models.Feeling{
			Id: 1,
			Name: "Relaxed",
			Activities: []string{"Read"},
		},
		"Accomplished": models.Feeling{
			Id: 2,
			Name: "Accomplished",
			Activities: []string{"Read"},
		},
		"Creative": models.Feeling{
			Id: 3,
			Name: "Creative",
			Activities: []string{"Write"},
		},
	}

	// Act
	activities, feelings, err := dataHandler.GetActivitiesAndFeelings()

	// Assert
	if err != nil {
		t.Errorf("Error message is '%s' but expected nil", err.Error())
	}

	assertActivities(t, activities, sort(values(dataHandler.activitiesByName), func(a models.Activity, b models.Activity) int {
		return a.Id - b.Id
	}))
	assertFeelings(t, feelings, sort(values(dataHandler.feelingsByName), func(a models.Feeling, b models.Feeling) int {
		return a.Id - b.Id
	}))
}

/**
 * Assertions to help test the in-memory model
 */
func assertEntry(t *testing.T, entry *models.Entry, activity *dtos.Activity, id int) {
	if entry.Id != id {
		t.Errorf("Entry has id %d but should have %d", entry.Id, id)
	}

	if entry.Activity != activity.Name {
		t.Errorf("Entry has name '%s' but should have name '%s'", entry.Activity, activity.Name)
	}

	if len(entry.Feelings) != len(activity.Feelings) {
		t.Errorf("Entry has %d feelings but should have %d", len(entry.Feelings), len(activity.Feelings))
	}

	for i := range entry.Feelings {
		if entry.Feelings[i] != activity.Feelings[i] {
			t.Errorf("Entry has feeling '%s' but should have '%s'", entry.Feelings[i], activity.Feelings[i])
		}
	}

	if entry.Week != getWeek(&activity.Created) {
		t.Errorf("Entry has week '%s' but should have '%s'", entry.Week, getWeek(&activity.Created))
	}

	if entry.Created != activity.Created {
		t.Errorf("Entry has created '%s' but should have '%s'", entry.Created, activity.Created)
	}
}

func assertEntries(t *testing.T, actual []models.Entry, expected []models.Entry) {
	if len(actual) != len(expected) {
		t.Errorf("Entries has %d but expected %d", len(actual), len(expected))
	}

	for i := range actual {
		if actual[i].Id != expected[i].Id {
			t.Errorf("Entry %d has id %d but expected %d", i, actual[i].Id, expected[i].Id)
		}

		if actual[i].Activity != expected[i].Activity {
			t.Errorf("Entry %d has activity '%s' but expected '%s'", i, actual[i].Activity, expected[i].Activity)
		}

		if len(actual[i].Feelings) != len(expected[i].Feelings) {
			t.Errorf("Entry %d has %d feelings but expected %d", i, len(actual[i].Feelings), len(expected[i].Feelings))
		}

		for j := range actual[i].Feelings {
			if actual[i].Feelings[j] != expected[i].Feelings[j] {
				t.Errorf("Entry %d has feeling %d of '%s' but expected '%s'", i, j, actual[i].Feelings[j], expected[i].Feelings[j])
			}
		}

		if actual[i].Week != expected[i].Week {
			t.Errorf("Entry %d has week '%s' but expected '%s'", i, actual[i].Week, expected[i].Week)
		}

		if actual[i].Created != expected[i].Created {
			t.Errorf("Entry %d has created '%s' but expected '%s'", i, actual[i].Created, expected[i].Created)
		}
	}
}

func assertActivity(t *testing.T, dataHandler *InMemoryDataHandler, name string, id int, feelings []string) {
	activity, exists := dataHandler.activitiesByName[name]

	if !exists {
		t.Errorf("Activities by name should contain key for '%s'", name)
	}

	if activity.Id != id {
		t.Errorf("Activity has id %d but should have %d", activity.Id, id)
	}

	if activity.Name != name {
		t.Errorf("Activity has name '%s' but should have '%s'", activity.Name, name)
	}

	if len(activity.Feelings) != len(feelings) {
		t.Errorf("Activity has %d feelings but should have %d", len(activity.Feelings), len(feelings))
	}

	for i := range activity.Feelings {
		if activity.Feelings[i] != feelings[i] {
			t.Errorf("Activity has feeling '%s' but should have '%s'", activity.Feelings[i], feelings[i])
		}
	}
}

func assertActivities(t *testing.T, actual []models.Activity, expected []models.Activity) {
	if len(actual) != len(expected) {
		t.Errorf("Activities has %d entries but should have %d", len(actual), len(expected))
	}

	for i := range actual {
		if actual[i].Id != expected[i].Id {
			t.Errorf("Activity %d has id %d but should have %d", i, actual[i].Id, expected[i].Id)
		}

		if actual[i].Name != expected[i].Name {
			t.Errorf("Activity %d has name '%s' but should have '%s'", i, actual[i].Name, expected[i].Name)
		}

		if len(actual[i].Feelings) != len(expected[i].Feelings) {
			t.Errorf("Activity %d has %d feelings but should have %d", i, len(actual[i].Feelings), len(expected[i].Feelings))
		}

		for j := range actual[i].Feelings {
			if actual[i].Feelings[j] != expected[i].Feelings[j] {
				t.Errorf("Activity %d has feeling %d with name '%s' but should have '%s'", i, j, actual[i].Feelings[j], expected[i].Feelings[j])
			}
		}
	}
}

func assertFeeling(t *testing.T, dataHandler *InMemoryDataHandler, name string, id int, activities []string) {
	feeling, exists := dataHandler.feelingsByName[name]

	if !exists {
		t.Errorf("Feelings by name should contain key for '%s'", name)
	}

	if feeling.Id != id {
		t.Errorf("Feeling has id %d but should have %d", feeling.Id, id)
	}

	if feeling.Name != name {
		t.Errorf("Feeling has name '%s' but should have '%s'", feeling.Name, name)
	}

	if len(feeling.Activities) != len(activities) {
		t.Errorf("Feeling has %d activities but should have %d", len(feeling.Activities), len(activities))
	}

	for i := range feeling.Activities {
		if feeling.Activities[i] != activities[i] {
			t.Errorf("Feeling has activity '%s' but should have '%s'", feeling.Activities[i], activities[i])
		}
	}
}

func assertFeelings(t *testing.T, actual []models.Feeling, expected []models.Feeling) {
	if len(actual) != len(expected) {
		t.Errorf("Feelings has %d entries but should have %d", len(actual), len(expected))
	}

	for i := range actual {
		if actual[i].Id != expected[i].Id {
			t.Errorf("Feeling %d has id %d but should have %d", i, actual[i].Id, expected[i].Id)
		}

		if actual[i].Name != expected[i].Name {
			t.Errorf("Feeling %d has name '%s' but should have '%s'", i, actual[i].Name, expected[i].Name)
		}

		if len(actual[i].Activities) != len(expected[i].Activities) {
			t.Errorf("Feeling %d has %d activities but should have %d", i, len(actual[i].Activities), len(expected[i].Activities))
		}

		for j := range actual[i].Activities {
			if actual[i].Activities[j] != expected[i].Activities[j] {
				t.Errorf("Feeling %d has activity %d with name '%s' but should have '%s'", i, j, actual[i].Activities[j], expected[i].Activities[j])
			}
		}
	}
}

func assertWeek(t *testing.T, dataHandler *InMemoryDataHandler, name string, id int) {
	week, exists := dataHandler.weeksByName[name]

	if !exists {
		t.Errorf("Weeks by name should contain key for '%s'", name)
	}

	assertWeekFields(t, &week, name, id)
}

func assertWeekFields(t *testing.T, week *models.Week, name string, id int) {
	if week.Id != id {
		t.Errorf("Week has id %d but should have %d", week.Id, id)
	}

	if week.Name != name {
		t.Errorf("Week has name '%s' but should have '%s'", week.Name, name)
	}

	if start, _ := getTime(name); week.Start != start {
		t.Errorf("Week has start '%s' but should have '%s'", week.Start, start)
	}

	if start, _ := getTime(name); week.End != start.AddDate(0, 0, 6) {
		t.Errorf("Week has end '%s' but should have '%s'", week.End, start.AddDate(0, 0, 6))
	}
}