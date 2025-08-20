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
	return models.Week{}, []models.Entry{}, nil
}

func (handler *LocalHostDataHandler) GetActivitiesAndFeelings() ([]models.Activity, []models.Feeling, error) {
	return []models.Activity{}, []models.Feeling{}, nil
}

/**
 * Database helper methods
 */
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
