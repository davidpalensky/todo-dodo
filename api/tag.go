package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"todo-dodo/db"
)

type TagCreateArgs struct {
	User_id uint
	Title   string
	Color   string
}

// API endpoint for creating a tag
func TagCreate(writer http.ResponseWriter, request *http.Request) {
	// Get user data and unmarshall
	var args TagCreateArgs
	req, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "Error: Could not read http request.")
		return
	}
	err1 := json.Unmarshal(req, &args)
	if err1 != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "Error: Could not deserialize %s", request.Body)
		return
	}

	// Insert into db
	_, err2 := db.DB.Exec("INSERT INTO tag (title, color, user_id) VALUES (?, ?, ?);", args.Title, args.Color, args.User_id)
	if err2 != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "Error: Could not fetch data from database")
		log.Printf("Error: db error: %s", err2)
		return
	}

	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, "Success: Tag added successfully")
	return
}
