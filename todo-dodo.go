package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"todo-dodo/api"
	"todo-dodo/db"

	"github.com/gorilla/mux"
	_ "github.com/libsql/libsql-client-go/libsql"
)

func main() {
	// Start server
	log.Println("Starting Server")

	// Connect to db
	db.DBConnect()
	defer db.DB.Close()

	// Setup router
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/api/v1/task/fetch", api.TaskFetch)
	router.HandleFunc("/api/v1/task/create", api.TaskCreate)

	// Setup server details
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
	}

	// Run server
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	log.Println("Server Started")

	// Prepare to accept SIGINT (Ctrl+C) to kill server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Block until we receive our signal.
	<-quit

	// Shut down server
	log.Println("Stoppping Server")

	// Don't understand this
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
	writer.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(writer, "Website not yet implemented.\n")
}
