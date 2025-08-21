package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"api.wellbeingquest.app/internal/dtos"
	"api.wellbeingquest.app/internal/models"

	_ "github.com/lib/pq"
)

type LocalHostDataHandler struct {
	driver string
	connection string
	setup []string
	db *sql.DB
}

func NewLocalHostDataHandler(driver string, connection string, setup []string) (*LocalHostDataHandler, error) {
	dataHandler := LocalHostDataHandler{
		driver: driver,
		connection: connection,
		setup: setup,
	}
	
	if err := ping(&dataHandler); err != nil {
		return nil, err
	}

	for _, setup := range dataHandler.setup {
		if err := runSqlCommands(&dataHandler, setup); err != nil {
				return nil, err
		}
	}

	return &dataHandler, nil
}

/**
 * Setup database
 */
func ping(dataHandler *LocalHostDataHandler) error {
	db, err := sql.Open(dataHandler.driver, dataHandler.connection)
	if err != nil {
		return fmt.Errorf("error opening postgres connection, error: %v", err)
	}
	
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("error pinging postgres connection, error: %v", err)
	}

	dataHandler.db = db
	return nil
}

func runSqlCommands(dataHandler *LocalHostDataHandler, sqlFile string) error {
	sqlBytes, err := os.ReadFile(sqlFile)
	if err != nil {
		return fmt.Errorf("error reading file, error: %v", err)
	}

	_, err = dataHandler.db.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("error executing sql command, error: %v", err)
	}

	return nil
}

/**
 * Database methods
 */
func (handler *LocalHostDataHandler) AddActivity(activityRequested *dtos.Activity) (models.Entry, error) {
	var entry models.Entry

	// Create or get activity, feelings, and week
	activity, err := handler.createOrGetActivity(activityRequested.Name)
	if err != nil {
		return entry, err
	}

	log.Printf("Activity, id %d, name '%s'\n", activity.Id, activity.Name)

	feelings := make([]models.Feeling, len(activityRequested.Feelings))
	for i, feelingRequested := range activityRequested.Feelings {
		feeling, err := handler.createOrGetFeeling(feelingRequested)
		if err != nil {
			return entry, err
		}

		feelings[i] = feeling
		log.Printf("Feeling, id %d, name '%s'\n", feeling.Id, feeling.Name)
	}

	weekRequested := getWeek(&activityRequested.Created)
	week, err := handler.createOrGetWeek(weekRequested)
	if err != nil {
		return entry, err
	}

	log.Printf("Week, id %d, name '%s'\n", week.Id, week.Name)

	// Link activities and feelings
	for _, feeling := range feelings {
		activityId, feelingId, err := handler.createOrGetActivityFeeling(activity.Id, feeling.Id)
		if err != nil {
			return entry, err
		}

		log.Printf("Activity feeling, %d, %d", activityId, feelingId)
	}

	// Create entry
	entry, err = handler.insertEntry(activity, week, activityRequested.Created)
	if err != nil {
		return entry, err
	}

	entry.Activity = activity.Name
	entry.Week = week.Name
	entry.Created = activityRequested.Created

	for _, feeling := range feelings {
		_, _, err := handler.insertEntryFeeling(entry.Id, feeling.Id)
		if err != nil {
			return entry, err
		}

		entry.Feelings = append(entry.Feelings, feeling.Name)
	}

	return entry, nil
}

func (handler *LocalHostDataHandler) GetWeekAndActivities(weekRequested *dtos.Week) (models.Week, []models.Entry, error) {
	var week models.Week
	var entries []models.Entry

	// Get week
	week, err := handler.getWeek(weekRequested.Name)
	if err != nil {
		return week, entries, err
	}

	start, err := getTime(week.Name)
	if err != nil {
		return week, entries, fmt.Errorf("error parsing time by week name '%s', error: %v", week.Name, err)
	}

	week.Start = start
	week.End = start.AddDate(0, 0, 6)
	log.Printf("Week '%s' retrieved\n", week.Name)

	// Get entries
	entries, err = handler.getEntries(week)
	if err != nil {
		return week, entries, err
	}

	log.Printf("Entries retrieved, count: %d\n", len(entries))
	return week, entries, nil
}

func (handler *LocalHostDataHandler) GetActivitiesAndFeelings() ([]models.Activity, []models.Feeling, error) {
	var activities []models.Activity
	var feelings []models.Feeling

	// Get activities
	activities, err := handler.getActivities()
	if err != nil {
		return activities, feelings, err
	}

	log.Printf("Activities retrieved, count: %d\n", len(activities))

	if len(activities) == 0 {
		return activities, feelings, fmt.Errorf("activities are empty")
	}

	// Get feelings
	feelings, err = handler.getFeelings()
	if err != nil {
		return activities, feelings, err
	}

	log.Printf("Feelings retrieved, count: %d\n", len(feelings))

	if len(feelings) == 0 {
		return activities, feelings, fmt.Errorf("feelings are empty")
	}

	return activities, feelings, nil
}

/**
 * Database helper methods
 */
func (handler *LocalHostDataHandler) getActivityById(id int) (models.Activity, error) {
	var activity models.Activity

	log.Printf("Querying activity by id, %d\n", id)
	err := handler.db.
		QueryRow("SELECT id, name FROM activities WHERE id = $1", id).
		Scan(&activity.Id, &activity.Name)

	if err != nil {
		log.Printf("Activity does not exist with id %d, error: %v\n", id, err)
		return activity, fmt.Errorf("error querying activity by id %d, error: %v", id, err)
	}

	return activity, nil
}

func (handler *LocalHostDataHandler) getActivities() ([]models.Activity, error) {
	var activities []models.Activity

	// Get rows
	log.Println("Querying activities")
	rows, err := handler.db.Query("SELECT id, name FROM activities")
	if err != nil {
		log.Printf("error queries activities, error: %v\n", err)
		return activities, fmt.Errorf("error querying activities, error: %v", err)
	}
	defer rows.Close()

	// Scan each row to append activities
	for rows.Next() {
		var activity models.Activity

		if err := rows.Scan(&activity.Id, &activity.Name); err != nil {
			return activities, fmt.Errorf("error scanning row for activity, error: %v", err)
		}

		activities = append(activities, activity)
	}

	if err := rows.Err(); err != nil {
		return activities, fmt.Errorf("error from rows for activities, error: %v", err)
	}

	return activities, nil
}

func (handler *LocalHostDataHandler) createOrGetActivity(name string) (models.Activity, error) {
	var activity models.Activity

	log.Printf("Querying activity by name, '%s'\n", name)
	err := handler.db.
		QueryRow("SELECT id, name FROM activities WHERE name = $1", name).
		Scan(&activity.Id, &activity.Name)

	if err == sql.ErrNoRows {
		log.Printf("Activity does not exist with name, '%s'\n", name)
		return handler.insertActivity(name)
	} else if err != nil {
		return activity, fmt.Errorf("error querying activity by name, '%s'", name)
	}

	return activity, nil
}

func (handler *LocalHostDataHandler) insertActivity(name string) (models.Activity, error) {
	var activity models.Activity

	log.Printf("Inserting activity by name, '%s'\n",name)
	err := handler.db.
		QueryRow("INSERT INTO activities (name) VALUES ($1) RETURNING id, name", name).
		Scan(&activity.Id, &activity.Name)

	if err != nil {
		return activity, fmt.Errorf("error inserting activity by name, '%s'", name)
	}

	return activity, nil
}

func (handler *LocalHostDataHandler) getFeelingById(id int) (models.Feeling, error) {
	var feeling models.Feeling

	log.Printf("Querying feeling by id, %d\n", id)
	err := handler.db.
		QueryRow("SELECT id, name FROM feelings WHERE id = $1", id).
		Scan(&feeling.Id, &feeling.Name)

	if err != nil {
		log.Printf("Feeling does not exist with id %d, error: %v\n", id, err)
		return feeling, fmt.Errorf("error querying feeling by id %d, error: %v", id, err)
	}

	return feeling, nil
}

func (handler *LocalHostDataHandler) getFeelings() ([]models.Feeling, error) {
	var feelings []models.Feeling

	// Get rows
	log.Println("Querying feelings")
	rows, err := handler.db.Query("SELECT id, name FROM feelings")
	if err != nil {
		log.Printf("error queries feelings, error: %v\n", err)
		return feelings, fmt.Errorf("error querying feelings, error: %v", err)
	}
	defer rows.Close()

	// Scan each row to append feelings
	for rows.Next() {
		var feeling models.Feeling

		if err := rows.Scan(&feeling.Id, &feeling.Name); err != nil {
			return feelings, fmt.Errorf("error scanning row for feeling, error: %v", err)
		}

		feelings = append(feelings, feeling)
	}

	if err := rows.Err(); err != nil {
		return feelings, fmt.Errorf("error from rows for feelings, error: %v", err)
	}

	return feelings, nil
}

func (handler *LocalHostDataHandler) createOrGetFeeling(name string) (models.Feeling, error) {
	var feeling models.Feeling

	log.Printf("Querying feeling by name, '%s'\n", name)
	err := handler.db.
		QueryRow("SELECT id, name FROM feelings WHERE name = $1", name).
		Scan(&feeling.Id, &feeling.Name)

	if err == sql.ErrNoRows {
		log.Printf("Feeling does not exist with name, '%s'\n", name)
		return handler.insertFeeling(name)
	} else if err != nil {
		return feeling, fmt.Errorf("error querying feeling by name, '%s'", name)
	}

	return feeling, nil
}

func (handler *LocalHostDataHandler) insertFeeling(name string) (models.Feeling, error) {
	var feeling models.Feeling

	log.Printf("Inserting feeling by name, '%s'\n",name)
	err := handler.db.
		QueryRow("INSERT INTO feelings (name) VALUES ($1) RETURNING id, name", name).
		Scan(&feeling.Id, &feeling.Name)

	if err != nil {
		return feeling, fmt.Errorf("error inserting feeling by name, '%s'", name)
	}

	return feeling, nil
}

func (handler *LocalHostDataHandler) getWeek(name string) (models.Week, error) {
	var week models.Week

	log.Printf("Querying week by name, '%s'\n", name)
	err := handler.db.
		QueryRow("SELECT id, name FROM weeks WHERE name = $1", name).
		Scan(&week.Id, &week.Name)

	if err != nil {
		log.Printf("Week does not exist with name, '%s'\n", name)
		return week, fmt.Errorf("error querying week by name, '%s'", name)
	}

	return week, nil
}

func (handler *LocalHostDataHandler) createOrGetWeek(name string) (models.Week, error) {
	var week models.Week

	log.Printf("Querying week by name, '%s'\n", name)
	err := handler.db.
		QueryRow("SELECT id, name FROM weeks WHERE name = $1", name).
		Scan(&week.Id, &week.Name)

	if err == sql.ErrNoRows {
		log.Printf("Week does not exist with name, '%s'\n", name)
		return handler.insertWeek(name)
	} else if err != nil {
		return week, fmt.Errorf("error querying week by name, '%s'", name)
	}

	return week, nil
}

func (handler *LocalHostDataHandler) insertWeek(name string) (models.Week, error) {
	var week models.Week

	log.Printf("Inserting week by name, '%s'\n",name)
	err := handler.db.
		QueryRow("INSERT INTO weeks (name) VALUES ($1) RETURNING id, name", name).
		Scan(&week.Id, &week.Name)

	if err != nil {
		return week, fmt.Errorf("error inserting week by name, '%s'", name)
	}

	return week, nil
}

func (handler *LocalHostDataHandler) getEntries(week models.Week) ([]models.Entry, error) {
	var entries []models.Entry

	// Get entry rows
	log.Printf("Querying entries by week id, %d\n", week.Id)
	rows, err := handler.db.Query("SELECT id, activityId, weekId, created FROM entries WHERE weekId = $1", week.Id)
	if err != nil {
		log.Printf("error queries entries by week id %d, error: %v\n", week.Id, err)
		return entries, fmt.Errorf("error querying entries by week id %d, error: %v", week.Id, err)
	}
	defer rows.Close()

	// Scan each row to append entries
	activitiesById := map[int]models.Activity{}
	feelingsById := map[int]models.Feeling{}

	for rows.Next() {
		var id int
		var activityId int
		var weekId int
		var created time.Time

		if err := rows.Scan(&id, &activityId, &weekId, &created); err != nil {
			return entries, fmt.Errorf("error scanning row for entry for week id %d, error: %v", weekId, err)
		}

		// Get activities
		activity, exists := activitiesById[activityId]
		if !exists {
			activity, err = handler.getActivityById(activityId)

			if err != nil {
				return entries, fmt.Errorf("error retrieving activity id %d for entry id %d, error: %v", activityId, id, err)
			}

			activitiesById[activity.Id] = activity
		}

		// Get feelings
		feelings := []string{}

		entryFeelings, err := handler.getEntryFeelingsByEntryId(id)
		if err != nil {
			return entries, fmt.Errorf("error retrieving entry feelings for entry id %d, error: %v", id, err)
		}

		for _, entryFeeling := range entryFeelings {
			feeling, exists := feelingsById[entryFeeling.FeelingId]
			if !exists {
				feeling, err = handler.getFeelingById(entryFeeling.FeelingId)

				if err != nil {
					return entries, fmt.Errorf("error retrieving feeling id %d for entry id %d, error: %v", entryFeeling.FeelingId, id, err)
				}

				feelingsById[feeling.Id] = feeling
			}

			// Add to feelings list
			feelings = append(feelings, feeling.Name)
		}

		// Add to entries
		entry := models.Entry{
			Id: id,
			Activity: activity.Name,
			Feelings: feelings,
			Week: week.Name,
			Created: created,
		}
		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return entries, fmt.Errorf("error from rows for entry for week id %d, error: %v", week.Id, err)
	}

	return entries, nil
}

func (handler *LocalHostDataHandler) insertEntry(activity models.Activity, week models.Week, created time.Time) (models.Entry, error) {
	var entry models.Entry

	log.Printf("Inserting entry, '%s', '%s'\n", activity.Name, week.Name)
	err := handler.db.
		QueryRow("INSERT INTO entries (activityId, weekId, created) VALUES ($1, $2, $3) RETURNING id", activity.Id, week.Id, created).
		Scan(&entry.Id)

	if err != nil {
		return entry, fmt.Errorf("error inserting entry, '%s', '%s'", activity.Name, week.Name)
	}

	return entry, nil
}

func (handler *LocalHostDataHandler) createOrGetActivityFeeling(activityId int, feelingId int) (int, int, error) {
	var activityIdRetrieved int
	var feelingIdRetrieved int

	log.Printf("Querying activity feeling by ids, %d, %d\n", activityId, feelingId)
	err := handler.db.
		QueryRow("SELECT activityId, feelingId FROM activityFeelings WHERE activityId = $1 AND feelingId = $2", activityId, feelingId).
		Scan(&activityIdRetrieved, &feelingIdRetrieved)

	if err == sql.ErrNoRows {
		log.Printf("Activity feeling does not exist with ids, %d, %d\n", activityId, feelingId)
		return handler.insertActivityFeeling(activityId, feelingId)
	} else if err != nil {
		return activityIdRetrieved, feelingIdRetrieved, fmt.Errorf("error querying activity feeling by ids, %d, %d", activityId, feelingId)
	}

	return activityIdRetrieved, feelingIdRetrieved, nil
}

func (handler *LocalHostDataHandler) insertActivityFeeling(activityId int, feelingId int) (int, int, error) {
	var activityIdRetrieved int
	var feelingIdRetrieved int

	log.Printf("Inserting activity feeling by ids, %d, %d\n", activityId, feelingId)
	err := handler.db.
		QueryRow("INSERT INTO activityFeelings (activityId, feelingId) VALUES ($1, $2) RETURNING activityId, feelingId", activityId, feelingId).
		Scan(&activityIdRetrieved, &feelingIdRetrieved)

	if err != nil {
		return activityIdRetrieved, feelingIdRetrieved, fmt.Errorf("error inserting activity feeling by ids, %d, %d", activityId, feelingId)
	}

	return activityIdRetrieved, feelingIdRetrieved, nil
}

func (handler *LocalHostDataHandler) getEntryFeelingsByEntryId(id int) ([]models.EntryFeeling, error) {
	var entryFeelings []models.EntryFeeling

	// Get entry feeling rows
	log.Printf("Querying entry feelings by entry id, %d\n", id)
	rows, err := handler.db.Query("SELECT entryId, feelingId FROM entryFeelings WHERE entryId = $1", id)
	if err != nil {
		log.Printf("error querying entry feelings by entry id %d, error: %v\n", id, err)
		return entryFeelings, fmt.Errorf("error querying entry feelings by entry id %d, error: %v", id, err)
	}
	defer rows.Close()

	// Scan each row to append entry feelings
	for rows.Next() {
		var entryFeeling models.EntryFeeling

		if err := rows.Scan(&entryFeeling.EntryId, &entryFeeling.FeelingId); err != nil {
			return entryFeelings, fmt.Errorf("error scanning row for entry feelings for entry id %d, error: %v", id, err)
		}

		entryFeelings = append(entryFeelings, entryFeeling)
	}

	if err := rows.Err(); err != nil {
		return entryFeelings, fmt.Errorf("error from rows for entry feelings for entry id %d, error: %v", id, err)
	}

	return entryFeelings, nil
}

func (handler *LocalHostDataHandler) insertEntryFeeling(entryId int, feelingId int) (int, int, error) {
	var entryIdRetrieved int
	var feelingIdRetrieved int

	log.Printf("Inserting entry feeling by ids, %d, %d\n", entryId, feelingId)
	err := handler.db.
		QueryRow("INSERT INTO entryFeelings (entryId, feelingId) VALUES ($1, $2) RETURNING entryId, feelingId", entryId, feelingId).
		Scan(&entryIdRetrieved, &feelingIdRetrieved)

	if err != nil {
		return entryIdRetrieved, feelingIdRetrieved, fmt.Errorf("error inserting entry feeling by ids, %d, %d", entryId, feelingId)
	}

	return entryIdRetrieved, feelingIdRetrieved, nil
}
