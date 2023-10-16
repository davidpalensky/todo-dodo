package api

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func DBConnect() {
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
