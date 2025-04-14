package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

//go:embed data/schema.sql
var schemaSQL string

// initDB initializes the database and creates the schema if it doesn't exist.
func initDB(path string) (*sql.DB, error) {
	isNew := false
	if _, err := os.Stat(path); os.IsNotExist(err) {
		isNew = true
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("cannot open db: %v", err)
	}
	if isNew {
		if _, err := db.Exec(schemaSQL); err != nil {
			return nil, fmt.Errorf("cannot init schema: %v", err)
		}
	}

	return db, nil
}
