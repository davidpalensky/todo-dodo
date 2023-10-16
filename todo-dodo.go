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

	"github.com/gorilla/mux"
	_ "github.com/libsql/libsql-client-go/libsql"
)

func main() {
	// Log messages
	log.Println("Starting Server")

	api.DBConnect()
	defer api.DB.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/api/v1/task/fetch", api.FetchTasks)
	router.HandleFunc("/api/v1/task/create", api.CreateTask)

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
