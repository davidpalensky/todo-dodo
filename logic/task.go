package logic

import (
	"fmt"
	"todo-dodo/db"
)

// The expected data from the client when creating tasks
type TaskCreator struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	// Unix timestamp
	Deadline uint `json:"deadline"`
	// Most correspond to a user in the users table
	User_id uint `json:"user_id"`
	// The associated tag data
	Tags []TagCreator `json:"tags"`
}

// Adds tasks to db
// TODO: Cleanup
// FIXME: Make transactions work correctly
func TaskCreateBatch(args []TaskCreator) error {
	db.DB.Exec("BEGIN TRANSACTION;")
	for _, row := range args {

		// Insert task
		res, err1 := db.DB.Exec("INSERT INTO tasks (title, content, deadline, user_id) VALUES (?, ?, ?, ?) RETURNING task_id;", row.Title, row.Content, row.Deadline, row.User_id)
		task_id, err2 := res.LastInsertId()
		if err1 != nil {
			db.DB.Exec("ROLLBACK;")
			//log.Printf("Error: Could not insert data into db: %s", err1.Error())
			return err1
		}
		if err2 != nil {
			db.DB.Exec("ROLLBACK;")
			//log.Printf("Error: Could not insert data into db: %s", err)
			return err2
		}

		// Do tag stuff
		task_id_uint64 := uint64(task_id)
		err3 := TagCreateBatch(row.Tags, &task_id_uint64)
		if err3 != nil {
			db.DB.Exec("ROLLBACK;")
			return err3
		}

	}
	db.DB.Exec("COMMIT TRANSACTION;")
	return nil
}

// Deletes a batch of tasks, including their task_tag_links entries
func TaskDeleteBatch(task_ids []uint64) error {
	for _, task_id := range task_ids {
		_, err1 := db.DB.Exec("DELETE FROM task_tag_links WHERE task_id = ?;", task_id)
		if err1 != nil {
			return &ActionError{Kind: "database", Msg: "Unable to delete task with task_id " + fmt.Sprintf("%d", task_id)}
		}
		_, err := db.DB.Exec("DELETE FROM tasks WHERE task_id = ?;", task_id)
		if err != nil {
			return &ActionError{Kind: "database", Msg: "Unable to delete task with task_id " + fmt.Sprintf("%d", task_id)}
		}
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////

// The expected information from the client when fetching tasks
type TaskFetchArgs struct {
	User_id uint64 `json:"user_id"`
}

// The data returned from the database when fetching tasks
type TaskModel struct {
	Task_id   uint64 `json:"task_id"`
	User_id   uint64 `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Creation  uint64 `json:"creation"`
	Deadline  uint64 `json:"deadline"`
	Completed bool   `json:"completed"`
}

// Fetches all tasks from db
// TODO: Make user specific and add auth
func TaskFetchAllDB(args *TaskFetchArgs) ([]TaskModel, error) {
	var result []TaskModel
	err := db.DB.Select(&result, "SELECT * FROM tasks WHERE user_id = ?;", args.User_id)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type TasksWithTagIds struct {
	Task_data TaskModel `json:"task_data"`
	Tag_ids   []uint64  `json:"tag_ids"`
}

type TasksWithTags struct {
	Tasks []TasksWithTagIds `json:"tasks"`
	Tags  []TagModel        `json:"tags"`
}

// Fetch all tasks including their associated tags
func TaskFetchAllWithTags(args *TaskFetchArgs) (*TasksWithTags, error) {
	var tasks []TaskModel
	err := db.DB.Select(&tasks, "SELECT * FROM tasks WHERE user_id = ?;", args.User_id)
	if err != nil {
		return nil, err
	}

	var tag_ids [][]uint64
	for _, task := range tasks {
		var current_ids []uint64
		err := db.DB.Select(&current_ids, "SELECT tag_id FROM task_tag_links WHERE task_id = ?;", task.Task_id)
		if err != nil {
			return nil, err
		}
		tag_ids = append(tag_ids, current_ids)
	}

	var tasks_with_tags []TasksWithTagIds
	for idx, task := range tasks {
		tasks_with_tags = append(tasks_with_tags, TasksWithTagIds{Task_data: task, Tag_ids: tag_ids[idx]})
	}

	var tags []TagModel
	err1 := db.DB.Select(&tags, "SELECT * FROM tags WHERE user_id = ?;", args.User_id)
	if err1 != nil {
		return nil, err1
	}

	return &TasksWithTags{Tasks: tasks_with_tags, Tags: tags}, nil
}