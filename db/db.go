package db

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/jmoiron/sqlx"
)

// The database connection
var DB *sqlx.DB

// Connect to database, should only be called once.
// Make sure environment variables "TODO_DODO_DB_URL" and "TODO_DODO_DB_TOKEN" are set correctly.
func Connect() {
	// url should be in the format: libsql://{database}.turso.io?authToken={auth token}
	url := os.Getenv("TODO_DODO_DB_URL")
	token := os.Getenv("TODO_DODO_DB_TOKEN")
	url = url + "?authToken=" + token

	db, err := sqlx.Open("libsql", url)
	if err != nil {
		log.Fatalf("Could not connect to db: \nurl = %s.", url)
		os.Exit(1)
	}
	DB = db
}
