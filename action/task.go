package action

import (
	"log"
	"todo-dodo/db"
)

// The expected data from the client when creating tasks
type TaskCreateArgs struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	// Unix timestamp
	Deadline uint `json:"deadline"`
	// Most correspond to a user in the users table
	User_id uint `json:"user_id"`
	// The associated tag data
	Tags []TagModel `json:"tags"`
}

// Adds tasks to db
// TODO: Cleanup
// FIXME: Make transactions work correctly
func TaskCreate(args []TaskCreateArgs) error {
	// Golang why dont you just allow ignoring both return values of a tuple
	db.DB.Exec("BEGIN TRANSACTION;")
	for _, row := range args {

		// Insert task
		res, err := db.DB.Exec("INSERT INTO tasks (title, content, deadline, user_id) VALUES (?, ?, ?, ?) RETURNING task_id;", row.Title, row.Content, row.Deadline, row.User_id)
		task_id, err1 := res.LastInsertId()
		if err1 != nil {
			db.DB.Exec("ROLLBACK;")
			//log.Printf("Error: Could not insert data into db: %s", err)
			log.Fatalf("This query should not fail 1: %s", err1.Error())
		}
		if err != nil {
			db.DB.Exec("ROLLBACK;")
			//log.Printf("Error: Could not insert data into db: %s", err)
			return &ActionError{Kind: "database", Msg: "TaskCreate: Could not enter tasks into db: " + err.Error()}
		}

		// Do tag stuff
		task_id_uint64 := uint64(task_id)
		err3 := TagCreateBatch(row.Tags, &task_id_uint64)
		if err3 != nil {
			db.DB.Exec("ROLLBACK;")
			return &ActionError{Kind: "database", Msg: "TaskCreate: Could not enter tags into db: " + err.Error()}
		}

	}
	db.DB.Exec("COMMIT;")
	return nil
}

// Deletes a batch of tasks, including their task_tag_links entries
func TaskDeleteBatch(task_ids []uint64) error {
	for _, task_id := range task_ids {
		_, err := db.DB.Exec("DELETE FROM tasks WHERE task_id = ?;", task_id)
		if err != nil {
			db.DB.Exec("ROLLBACK;")
			return &ActionError{Kind: "database", Msg: "Unable to delete task with task_id " + string(task_id)}
		}
	}
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

// Fetches all tasks from db
// TODO: Make user specific and add auth
func TaskFetchDB(args *TaskFetchArgs) ([]TaskFetchDBReturn, error) {
	// Query DB
	var result []TaskFetchDBReturn
	err := db.DB.Select(&result, "SELECT * FROM tasks WHERE user_id = ?;", args.User_id)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type TaskFetchAllWithTagsReturn struct {
	Task_data TaskFetchDBReturn `json:"task_data"`
	Tags      []TagModel        `json:"tags"`
}

// Fetches all tasks from db aswell as the tags for each
/*
func TaskFetchAllWithTags(args *TaskFetchArgs) ([]TaskFetchAllWithTagsReturn, error) {
	// Query DB
	var result []TaskFetchAllWithTagsReturn
	err := db.DB.Select(&result.Task_data, "SELECT * FROM tasks WHERE user_id = ?;", args.User_id)

	if err != nil {
		return nil, err
	}
	return result, nil
}
*/
