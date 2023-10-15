package api

import (
	"database/sql"
	"os"
)

// "libsql://[your-database].turso.io?authToken=[your-auth-token]"
var DB_URL = os.Getenv("TODO_DODO_DB_URL")
var DB_AUTH_TOKEN = os.Getenv("TODO_DODO_DB_TOKEN")

// Gotta love global variables
var DB, DB_ERR = sql.Open("libsql", DB_URL)
