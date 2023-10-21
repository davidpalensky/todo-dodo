package logic

import (
	"fmt"
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
	Tags []TagCreator `json:"tags"`
}

// Adds tasks to db
func TaskCreateBatch(args []TaskCreateArgs) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	// OPTIMASATION: Use tx.NamedExec() to put all values into a single query.
	for _, row := range args {

		// Insert task
		res, err1 := tx.Exec("INSERT INTO tasks (title, content, deadline, user_id) VALUES (?, ?, ?, ?) RETURNING task_id;", row.Title, row.Content, row.Deadline, row.User_id)
		task_id, err2 := res.LastInsertId()
		if err1 != nil {
			tx.Rollback()
			//log.Printf("Error: Could not insert data into db: %s", err1.Error())
			return err1
		}
		if err2 != nil {
			tx.Rollback()
			//log.Printf("Error: Could not insert data into db: %s", err)
			return err2
		}
		tx.Commit()

		// Do tag stuff
		task_id_uint64 := uint64(task_id)
		err3 := TagCreateBatch(row.Tags, task_id_uint64)
		if err3 != nil {
			tx.Rollback()
			return err3
		}

	}

	return nil
}

type TaskDeleteBatchArgs struct {
	Task_ids []uint64 `json:"task_ids"`
}

// Deletes a batch of tasks, including their task_tag_links entries
func TaskDeleteBatch(a TaskDeleteBatchArgs) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	for _, task_id := range a.Task_ids {
		_, err1 := tx.Exec("DELETE FROM task_tag_links WHERE task_id = ?;", task_id)
		if err1 != nil {
			tx.Rollback()
			return &LogicError{Kind: "database", Msg: "Unable to delete task with task_id " + fmt.Sprintf("%d", task_id)}
		}
		_, err := tx.Exec("DELETE FROM tasks WHERE task_id = ?;", task_id)
		if err != nil {
			tx.Rollback()
			return &LogicError{Kind: "database", Msg: "Unable to delete task with task_id " + fmt.Sprintf("%d", task_id)}
		}
	}
	tx.Commit()
	return nil
}

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
// TODO: Add auth
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

type TaskUpdatArgs struct {
	Task_id   uint64   `json:"task_id"`
	Completed *bool    `json:"completed"`
	Deadline  *uint64  `json:"deadline"`
	Tag_ids   []uint64 `json:"tag_ids"`
}

// Updates/replaces a task to the given information. Does not change values provided as nil / an empty list is given.
// If tags is not an empty list, all the current tags will be removed and replaced with the newly provided ones.
func TaskUpdate(a TaskUpdatArgs) error {
	if a.Completed == nil && a.Deadline == nil && len(a.Tag_ids) == 0 {
		return nil
	}
	//log.Printf("TaskUpdate: Args = %v\n", a)
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	if a.Completed != nil {
		var completed int
		if *a.Completed { // *a.Completed == true
			completed = 1
		} else {
			completed = 0
		}
		_, err := tx.Exec("UPDATE tasks SET completed = ? WHERE task_id = ?;", completed, a.Task_id)
		if err != nil {
			log.Printf("Err: %s", err.Error())
			tx.Rollback()
			return err
		}
	}
	if a.Deadline != nil {
		_, err := tx.Exec("UPDATE tasks SET deadline = ? WHERE task_id = ?;", a.Deadline, a.Task_id)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(a.Tag_ids) != 0 {
		type TaskTagLinksInserts struct {
			Task_id uint64
			Tag_id  uint64
		}

		var inserts []TaskTagLinksInserts
		for _, tag_id := range a.Tag_ids {
			inserts = append(inserts, TaskTagLinksInserts{Task_id: a.Task_id, Tag_id: tag_id})
		}

		_, err = tx.NamedExec("INSERT INTO task_tag_links (task_id, tag_id) VALUES (:task_id, :tag_id)", inserts)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}
