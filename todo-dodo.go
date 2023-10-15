package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	_ "github.com/libsql/libsql-client-go/libsql"
)

// "libsql://[your-database].turso.io?authToken=[your-auth-token]"
var DB_URL = os.Getenv("TODO_DODO_DB_URL")
var DB_AUTH_TOKEN = os.Getenv("TODO_DODO_DB_TOKEN")

// Gotta love global variables
var DB, DB_ERR = sql.Open("libsql", DB_URL)

func HomeHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, "Message Recieved.\n")
}

func FetchTasks(writer http.ResponseWriter, request *http.Request) {
	db_result, err := DB.Query("SELECT * FROM tasks;")
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "Error: Could not fetch data from database")
		log.Printf("Error: db error: %s", err)
		return
	}
	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, "db_result: %v\n", db_result)
	return
}

// TODO: Implement user-authS
type TaskCreationCommand struct {
	title   string
	content string
	due     uint
	user_id uint
}

func CreateTask(writer http.ResponseWriter, request *http.Request) {
	var args TaskCreationCommand
	err := json.NewDecoder(request.Body).Decode(&args)
	if err != nil {
		fmt.Fprintf(writer, "Error could not deserialize task creation arguements: %s", err.Error())
		return
	}

}

func main() {
	// Log messages
	log.Println("Starting Server.")
	defer log.Println("Stoppped Server.")

	// Check DB connection
	if DB_ERR != nil {
		log.Printf("Failed to connect to db %s: %s", DB_URL, DB_ERR)
		os.Exit(1)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/api/v1/task/fetch", FetchTasks)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
	}
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	os.Exit(0)
}
