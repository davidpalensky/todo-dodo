package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/libsql/libsql-client-go/libsql"
)

// "libsql://[your-database].turso.io?authToken=[your-auth-token]"
var DB_URL = os.Getenv("TODO_DODO_DB_URL")
var DB_AUTH_TOKEN = os.Getenv("TODO_DODO_DB_TOKEN")

// Gotta love global variables
var DB, DB_ERR = sql.Open("libsql", DB_URL)

func main() {
	// Log messages
	log.Println("Starting Server")

	// Check DB connection
	if DB_ERR != nil {
		log.Printf("Failed to connect to db %s: %s", DB_URL, DB_ERR)
		os.Exit(1)
	}

	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/api/v1/task/fetch", FetchTasks)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
	}
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	log.Println("Server Started")

	quit := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(quit, os.Interrupt)

	// Block until we receive our signal.
	<-quit
	log.Println("Stoppping Server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server Stopped")
	os.Exit(0)
}

func HomeHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, "Message Recieved.\n")
}

func FetchTasks(writer http.ResponseWriter, request *http.Request) {
	// TODO: Make user specific and add auth
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

// TODO: Implement user-auth
type TaskCreationCommand struct {
	title   string
	content string
	due     uint
	user_id uint
}

// TODO: Implement
func CreateTask(writer http.ResponseWriter, request *http.Request) {
	var args TaskCreationCommand
	err := json.NewDecoder(request.Body).Decode(&args)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "Error could not deserialize task creation arguements: %s", err.Error())
		return
	}

}
