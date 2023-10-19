package db

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// The database connection
var DB *sqlx.DB
var DB_OPENED bool = false

// Opens local database at ./db/todo-dodo.db
// Will not open again if already opened
func Open() {
	if DB_OPENED {
		return
	}
	path := "./db/todo-dodo.db"
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Could not open url with path: %s\nerr: %s\n", path, err.Error())
		os.Exit(1)
	}
	DB = db
	DB_OPENED = true
}

// Close db
func Close() {
	if !DB_OPENED {
		DB.Close()
		DB_OPENED = false
	}
}
