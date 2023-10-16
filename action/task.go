package action

import (
	"log"
	"todo-dodo/db"
)

// This error just wraps a string
type TaskError struct {
	Kind string
	// Developer friendly message
	Msg string
}

func (err *TaskError) Error() string {
	return "kind: " + err.Kind + ", msg:" + err.Msg
}

// The expected data from the client when creating tasks
type TaskCreateArgs struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	// Unix timestamp
	Deadline uint `json:"deadline"`
	// Most correspond to a user in the users table
	User_id uint `json:"user_id"`
}

// Adds tasks to db
// TODO: Cleanup
func TaskCreateImpl(args []TaskCreateArgs) error {
	// Golang why dont you just allow ignoring both return values of a tuple
	a, _ := db.DB.Exec("BEGIN TRANSACTION;")
	DoNothing(a)
	for _, row := range args {
		_, err := db.DB.Exec("INSERT INTO tasks (title, content, deadline, user_id) VALUES (?, ?, ?, ?);", row.Title, row.Content, row.Deadline, row.User_id)
		if err != nil {
			log.Printf("Error: Could not insert data into db: %s", err)
			return &TaskError{Kind: "internal", Msg: "Could not fetch data from db"}
		}
	}
	b, _ := db.DB.Exec("COMMIT;")
	DoNothing(b)
	return nil
}

// ----

// The expected information from the client when fetching tasks
type TaskFetchArgs struct {
	User_id uint64 `json:"user_id"`
}

// The data returned from the database when fetching tasks
type TaskFetchDBReturn struct {
	Task_id   uint   `json:"task_id"`
	User_id   uint   `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Creation  uint64 `json:"creation"`
	Deadline  uint64 `json:"deadline"`
	Completed bool   `json:"completed"`
}

// API endpoint for fetching tasks from the database
// TODO: Make user specific and add auth
func TaskFetch(args *TaskFetchArgs) ([]TaskFetchDBReturn, error) {
	// Query DB
	var result []TaskFetchDBReturn
	err := db.DB.Select(&result, "SELECT * FROM tasks WHERE user_id = ?;", args.User_id)

	if err != nil {
		return nil, err
	}
	return result, nil
}
