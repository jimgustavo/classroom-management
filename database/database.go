// database/database.go

package database

import (
	"database/sql"
)

var db *sql.DB

// InitializeDB initializes the database connection
func InitializeDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		return err
	}

	// Check if the database connection is successful
	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}
