package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// TODO: Implement user-auth
type CreateArgs struct {
	Title    string
	Content  string
	Deadline uint
	User_id  uint
}

// API endpoint for adding a task to the database
func CreateTask(writer http.ResponseWriter, request *http.Request) {
	// Get user data and unmarshall
	var args CreateArgs
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

	_, err2 := DB.Exec("INSERT INTO tasks (title, content, deadline, user_id) VALUES (?, ?, ?, ?);", args.Title, args.Content, args.Deadline, args.User_id)
	if err2 != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "Error: Could not fetch data from database")
		log.Printf("Error: db error: %s", err2)
		return
	}

	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, "Success: Task added successfully")
	return
}

type FetchArgs struct {
	User_id uint64
}

type FetchReturn struct {
	Task_id   uint
	User_id   uint
	Title     string
	Content   string
	Creation  uint64
	Deadline  uint64
	Completed bool
}

// API endpoint for fetching tasks from the database]
// TODO: Make user specific and add auth
func FetchTasks(writer http.ResponseWriter, request *http.Request) {
	// Get user data and unmarshall
	var args FetchArgs
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
	var result []FetchReturn
	err2 := DB.Select(&result, "SELECT * FROM tasks WHERE user_id = ?;", args.User_id)
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
