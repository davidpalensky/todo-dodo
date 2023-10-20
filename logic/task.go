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
	UserId uint `json:"user_id"`
	// The associated tag data
	Tags []TagCreator `json:"tags"`
}

// Adds tasks to db
func TaskCreateBatch(args []TaskCreator) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	// OPTIMASATION: Use tx.NamedExec() to put all values into a single query.
	for _, row := range args {

		// Insert task
		res, err1 := tx.Exec("INSERT INTO tasks (title, content, deadline, user_id) VALUES (?, ?, ?, ?) RETURNING task_id;", row.Title, row.Content, row.Deadline, row.UserId)
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

		// Do tag stuff
		task_id_uint64 := uint64(task_id)
		err3 := TagCreateBatch(row.Tags, &task_id_uint64)
		if err3 != nil {
			tx.Rollback()
			return err3
		}

	}
	tx.Commit()
	return nil
}

// Deletes a batch of tasks, including their task_tag_links entries
func TaskDeleteBatch(task_ids []uint64) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	for _, task_id := range task_ids {
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

///////////////////////////////////////////////////////////////////////////////

// The expected information from the client when fetching tasks
type TaskFetchArgs struct {
	User_id uint64 `json:"user_id"`
}

// The data returned from the database when fetching tasks
type TaskModel struct {
	TaskId    uint64 `json:"task_id"`
	UserId    uint64 `json:"user_id"`
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
	TaskData TaskModel `json:"task_data"`
	TagIds   []uint64  `json:"tag_ids"`
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
		err := db.DB.Select(&current_ids, "SELECT tag_id FROM task_tag_links WHERE task_id = ?;", task.TaskId)
		if err != nil {
			return nil, err
		}
		tag_ids = append(tag_ids, current_ids)
	}

	var tasks_with_tags []TasksWithTagIds
	for idx, task := range tasks {
		tasks_with_tags = append(tasks_with_tags, TasksWithTagIds{TaskData: task, TagIds: tag_ids[idx]})
	}

	var tags []TagModel
	err1 := db.DB.Select(&tags, "SELECT * FROM tags WHERE user_id = ?;", args.User_id)
	if err1 != nil {
		return nil, err1
	}

	return &TasksWithTags{Tasks: tasks_with_tags, Tags: tags}, nil
}

type TaskUpdatArgs struct {
	TaskId     uint64   `json:"task_id"`
	Completion *bool    `json:"completion"`
	Deadline   *uint64  `json:"deadline"`
	TagIds     []uint64 `json:"tag_ids"`
}

// Updates/replaces a task to the given information. Does not change values provided as nil / an empty list is given.
// If tags is not an empty list, all the current tags will be removed and replaced with the newly provided ones.
func TaskUpdate(a TaskUpdatArgs) error {
	if a.Completion == nil && a.Deadline == nil && len(a.TagIds) == 0 {
		return nil
	}
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	if a.Completion != nil {
		var completion int
		if *a.Completion { // *a.Completion == true
			completion = 1
		} else {
			completion = 0
		}
		_, err := tx.Exec("UPDATE OR IGNORE tasks SET completion = ? WHERE task_id = ?;", completion, a.TaskId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if a.Deadline != nil {
		_, err := tx.Exec("UPDATE OR IGNORE tasks SET deadline = ? WHERE task_id = ?;", a.Deadline, a.TaskId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(a.TagIds) != 0 {
		type TaskTagLinksInserts struct {
			TaskId uint64
			TagId  uint64
		}

		var inserts []TaskTagLinksInserts
		for _, tag_id := range a.TagIds {
			inserts = append(inserts, TaskTagLinksInserts{TaskId: a.TaskId, TagId: tag_id})
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
