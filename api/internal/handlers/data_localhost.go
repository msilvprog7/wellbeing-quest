package handlers

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type LocalHostDataHandler struct {
	driver string
	connection string
	setup string
	db *sql.DB
}

func NewLocalHostDataHandler(connection string, setup string, driver string) (*LocalHostDataHandler, error) {
	dataHandler := LocalHostDataHandler{
		driver: driver,
		connection: connection,
		setup: setup,
	}
	
	if err := ping(&dataHandler); err != nil {
		return nil, err
	}

	if err := runSqlCommands(&dataHandler, dataHandler.setup); err != nil {
		return nil, err
	}

	return &dataHandler, nil
}

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