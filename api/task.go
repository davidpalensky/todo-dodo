package api

// TOPLEVEL TODO: Make it so that users can add tags when they create/delete tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"todo-dodo/db"
)

// The expected data from the client when creating tasks
type TaskCreateArgs struct {
	Title   string
	Content string
	// Unix timestamp
	Deadline uint
	// Most correspond to a user in the users table
	User_id uint
}

type TaskError struct {
	msg string
}

func (err *TaskError) Error() string {
	return err.msg
}

// API endpoint for adding a task to the database
func TaskCreate(writer http.ResponseWriter, request *http.Request) {
	// Get user data and unmarshall
	var args []TaskCreateArgs
	req, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "Could not read http request.")
		return
	}
	err1 := json.Unmarshal(req, &args)
	if err1 != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "Error could not deserialize task fetching arguements %s: %s", request.Body, err1.Error())
		return
	}

	if err := TaskCreateImpl(args); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "Error: could not create tasks: %s", err.Error())
		return
	}

	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, "Success: Task added successfully")
	return
}

// Adds task to db
func TaskCreateImpl(args []TaskCreateArgs) error {
	// Golang why dont you just allow ignoring both return values of a tuple
	a, _ := db.DB.Exec("BEGIN TRANSACTION;")
	DoNothing(a)
	for _, row := range args {
		_, err := db.DB.Exec("INSERT INTO tasks (title, content, deadline, user_id) VALUES (?, ?, ?, ?);", row.Title, row.Content, row.Deadline, row.User_id)
		if err != nil {
			log.Printf("Error: Could not fetch data from db: %s", err)
			return &TaskError{msg: "Could not fetch data from db"}
		}
	}
	b, _ := db.DB.Exec("COMMIT;")
	DoNothing(b)
	return nil
}

// The expected information from the client when fetching tasks
type TaskFetchArgs struct {
	User_id uint64
}

// The data returned from the database when fetching tasks
type TaskFetchReturnDB struct {
	Task_id   uint
	User_id   uint
	Title     string
	Content   string
	Creation  uint64
	Deadline  uint64
	Completed bool
}

// API endpoint for fetching tasks from the database
// TODO: Make user specific and add auth
func TaskFetch(writer http.ResponseWriter, request *http.Request) {
	// Get user data and unmarshall
	var args TaskFetchArgs
	req, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "Could not read http request.")
		return
	}
	err1 := json.Unmarshal(req, &args)
	if err1 != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "Error could not deserialize task fetching arguements %s: %s", request.Body, err1.Error())
		return
	}

	// Query DB
	var result []TaskFetchReturnDB
	err2 := db.DB.Select(&result, "SELECT * FROM tasks WHERE user_id = ?;", args.User_id)
	if err2 != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "Error: Could not fetch data from database")
		log.Printf("Error: db error: %s", err2)
		return
	}

	result_as_json, _ := json.Marshal(result)

	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, "db_result: %s\n", result_as_json)
	return
}
