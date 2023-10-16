package api2

import (
	"log"
	"net/http"
	"todo-dodo/db"

	"github.com/gin-gonic/gin"
)

// This error just wraps a string
type TaskError struct {
	msg string
}

func (err *TaskError) Error() string {
	return err.msg
}

// The expected data from the client when creating tasks
type TaskCreateArgs struct {
	Title   string
	Content string
	// Unix timestamp
	Deadline uint
	// Most correspond to a user in the users table
	User_id uint
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

func TaskCreateEnpoint(ctx *gin.Context) {
	args := new([]TaskCreateArgs)
	if err := ctx.Bind(args); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid arguements or incorrectly encoded json provided",
		})
		return
	}
	if err := TaskCreateImpl(*args); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "could not update record",
		})
		return
	}
	return
}
