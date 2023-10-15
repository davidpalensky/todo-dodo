package main

import (
	"context"
	"database/sql"
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
		log.Printf("Failed to open db: %s", err)
		return
	}
	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, "db_result: %v\n", db_result)
	return
}

func main() {
	log.Println("Starting Server.")
	defer log.Println("Stoppped Server.")
	if DB_ERR != nil {
		log.Printf("Failed to open db %s: %s", DB_URL, DB_ERR)
		os.Exit(1)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/api/v1/task/fetch", FetchTasks)

	http.Handle("/", r)
	http.Handle("/api/v1/task/fetch", r)

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

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 100)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("Stoppped Server.")
	os.Exit(0)
}
